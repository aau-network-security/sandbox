package sandbox

import (
	"context"
	"errors"
	"fmt"
	"github.com/aau-network-security/openvswitch/ovs"
	"github.com/aau-network-security/sandbox/config"
	"github.com/aau-network-security/sandbox/controller"
	"github.com/aau-network-security/sandbox/dnet/dns"
	"github.com/aau-network-security/sandbox/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"sync"
	"time"

	//"github.com/aau-network-security/sandbox2/models"
	"github.com/aau-network-security/sandbox/store"
	"github.com/aau-network-security/sandbox/virtual"
	"github.com/aau-network-security/sandbox/virtual/docker"
	"github.com/aau-network-security/sandbox/virtual/vbox"
	"github.com/rs/zerolog/log"
)

var (
	mgmin = 10443
	mgmax = 10553

	rmin = 5000
	rmax = 5300

	ErrVMNotCreated       = errors.New("no VM created")
	ErrGettingContainerID = errors.New("could not get container ID")
)

type environment struct {
	// challenge microservice should be integrated heres
	controller controller.NetController
	dockerHost docker.Host
	instances  []virtual.Instance
	ports      []string
	vlib       vbox.Library
	dnsServer  *dns.Server
}

type SandConfig struct {
	Name   string
	Tag    string
	env    *environment
	Config *config.Config
}

func NewSandbox(sandconf *SandConfig) (*SandConfig, error) {

	netController := controller.New()
	netController.IPPool = controller.NewIPPoolFromHost()

	dockerHost := docker.NewHost()
	vlib := vbox.NewLibrary(sandconf.Config.VmConfig.OvaDir)
	env := &environment{
		controller: *netController,
		dockerHost: dockerHost,
		vlib:       vlib,
	}

	sandconf.env = env

	log.Info().Msgf("New environment initialized ")
	return sandconf, nil
}

func (gc *SandConfig) StartSandbox(ctx context.Context, tag, name string, scenarios map[int]store.Scenario) error {

	scenario, ok := scenarios[0]
	if !ok {
		return status.Errorf(codes.InvalidArgument, "No scenario exists with that ID - See valid ID using list command")
	}

	log.Info().Str("Sandbox Tag", tag).
		Str("Sandbox Name", name).
		Str("Scenario", scenario.Name).
		Msg("starting sandbox")

	log.Debug().Str("Sandbox", name).Str("bridgeName", tag).Msg("creating openvswitch bridge")
	if err := gc.env.initializeOVSBridge(tag); err != nil {
		return err
	}

	log.Debug().Str("Game", name).Int("Networks", len(scenario.Networks)).Msg("Creating networks")
	if err := gc.env.createNetworks(tag, scenario.Networks); err != nil {
		return err
	}

	var vlanPorts []string
	for _, network := range scenario.Networks {
		vlanPorts = append(vlanPorts, fmt.Sprintf("%s_%s", tag, network.Name))
	}
	vlanPorts = append(vlanPorts, fmt.Sprintf("%s_monitoring", tag))

	log.Debug().Str("Sandbox", tag).Msgf("Initilizing OpnSense VM")

	//assign mgmt port to opnsense vm in 443 and 22
	mngtPort := getRandomPort(mgmin, mgmax)

	routerPort := getRandomPort(rmin, rmax)

	log.Debug().Str("Game", name).Msg("configuring monitoring")
	if err := gc.env.configureMonitor(ctx, tag, scenario.Networks); err != nil {
		log.Error().Err(err).Msgf("configuring monitoring")
		return err
	}

	if err := gc.env.initOpnSenseVM(ctx, tag, vlanPorts, mngtPort, routerPort); err != nil {
		log.Error().Err(err).Msg("Problem booting OpnSense VM")
		return err
	}

	log.Debug().Str("Game  ", name).Msg("starting DNS server")

	if err := gc.env.initDNSServer(ctx, tag); err != nil {
		log.Error().Err(err).Msg("attaching the DNS")
		return err
	}

	//initFTPMalws
	if err := gc.env.initFTPMalws(ctx, tag); err != nil {
		log.Error().Err(err).Msg("Problem booting targetWin VM")
		return err
	}

	log.Debug().Str("Game", name).Msg("initializing scenario")
	if err := gc.env.initializeScenario(ctx, tag, scenario); err != nil {
		return err
	}

	//TODO: Create FTP here!!!

	if err := gc.env.addTargetVM(ctx, tag); err != nil {
		log.Error().Err(err).Msg("Problem booting targetWin VM")
		return err
	}

	log.Info().Str("Game Tag", tag).
		Str("Game Name", name).
		Msg("started sandbox")

	return nil
}
func (gc *SandConfig) CloseSandbox(ctx context.Context) error {
	var waitg sync.WaitGroup
	var failed bool

	log.Info().Str("Game Name", gc.Name).Str("Game Tag", gc.Tag).Msg("Stopping sandbox")
	for _, instance := range gc.env.instances {
		waitg.Add(1)
		go func(vi virtual.Instance) {
			defer waitg.Done()
			if err := vi.Stop(); err != nil {
				log.Error().Str("Instance Type", vi.Info().Type).Str("Instance Name", vi.Info().Id).Msg("failed to stop virtual instance")
				failed = true
			}
			log.Debug().Str("Instance Type", vi.Info().Type).Str("Instance Name", vi.Info().Id).Msg("stopped instance")
			if err := vi.Close(); err != nil {
				log.Error().Str("Instance Type", vi.Info().Type).Str("Instance Name", vi.Info().Id).Msg("failed to close virtual instance")
				failed = true
			}

			if vi.Info().Type == "docker" {
				if err := gc.env.controller.Ovs.Docker.DeletePorts(gc.Tag, vi.Info().Id); err != nil {
					log.Error().Str("Instance Name", vi.Info().Id).Msg("Deleted all ports on docker image")
					failed = true
				}
			}
			log.Debug().Str("Instance Type", vi.Info().Type).Str("Instance Name", vi.Info().Id).Msg("closed instance")
		}(instance)

	}
	waitg.Wait()
	if failed {
		return errors.New("failed to stop an virtual instance")
	}

	if err := gc.env.removeNetworks(gc.Tag); err != nil {
		return errors.New("failed to remove networks")
	}

	return nil
}

func (env *environment) initDNSServer(ctx context.Context, bridge string) error {
	//New(bridge, IPanswer string)
	//defer wg.Done()
	DNS, err := dns.New(ctx, bridge)
	if err != nil {
		log.Error().Msgf("Error creating DNS server %v", err)
		return err
	}

	if DNS == nil {
		return ErrVirtualInstanceNil
	}

	env.instances = append(env.instances, DNS)
	//env.instances = append(env.instances, server )

	if err := DNS.Run(ctx); err != nil {
		log.Error().Msgf("Error in starting DNS  %v", err)
		return err
	}

	contID := DNS.ID()
	//HardCoded Mac Address Container
	fmt.Printf("AICI e ID = %s\n", contID)

	i := 1

	macAddress := "8a:3d:ec:9c:b6:a5"
	vlantag := "0"
	//sudo ovs-docker add-port test eth0 09 --vlan=10 --macaddress="8a:3d:ec:9c:b6:a5" --dhcp=true
	//TODO: Check if you need a vlan for DNS server
	if err := env.controller.Ovs.Docker.AddPort(bridge, fmt.Sprintf("eth%d", i), contID, ovs.DockerOptions{MACAddress: macAddress, VlanTag: vlantag, DHCP: true}); err != nil {

		log.Error().Err(err).Str("container", contID).Msg("adding port to DNS container")
		return err
	}

	return nil
}

func (env *environment) initFTPMalws(ctx context.Context, bridge string) error {
	//New(bridge, IPanswer string)
	//defer wg.Done()
	//Todo: asta nu merge !! FMMMSIII
	//var string malwarePath
	malwarePath := "/home/rvm/sandbox/bad/upload"
	ftp := docker.NewContainer(docker.ContainerConfig{
		Image: "atmoz/sftp",
		Mounts: []string{
			fmt.Sprintf("%s:/home/foo/upload", malwarePath),
			//fmt.Sprintf("",dir ),
		},
		Labels: map[string]string{
			"nap-sandbox": bridge,
			//"sandbox-networks": strings.Join(nets, ","),

		},
		Cmd: []string{"foo:pass:1001"},
	})

	if err := ftp.Create(ctx); err != nil {
		log.Error().Err(err).Msg("creating container")
		return err
	}

	if err := ftp.Start(ctx); err != nil {
		log.Error().Err(err).Msg("starting container")
		return err
	}

	cid := ftp.ID()
	if cid == "" {
		log.Error().Msg("getting ID for container")
		return ErrGettingContainerID
	}

	vlantag := "40"
	if err := env.controller.Ovs.Docker.AddPort(bridge, "eth0", cid, ovs.DockerOptions{MACAddress: "8a:3d:af:44:1b:f7", VlanTag: vlantag, DHCP: true}); err != nil {
		log.Error().Err(err).Str("container", cid).Msg("adding port to container")
		return err
	}

	if ftp == nil {
		return ErrVirtualInstanceNil
	}

	env.instances = append(env.instances, ftp)

	return nil
}

//configureMonitor will configure the monitoring VM by attaching the correct interfaces
func (env *environment) configureMonitor(ctx context.Context, bridge string, nets []models.Network) error {

	log.Info().Str("sandbox tag", bridge).Msg("creating monitoring network")
	if err := env.createPort(bridge, "monitoring", 0); err != nil {
		return err
	}

	//TODO: Create mirror for each VLAN
	//		Change the mirror AllBlue to VLAN specific
	//mirror := fmt.Sprintf("%s_mirror", bridge)
	//
	//log.Info().Str("sandbox tag", bridge).Msg("Creating the network mirror")
	//if err := env.controller.Ovs.VSwitch.CreateMirrorforBridge(mirror, bridge); err != nil {
	//	log.Error().Err(err).Msg("creating mirror")
	//	return err
	//}
	//
	//if err := env.createPort(bridge, "AllBlue", 0); err != nil {
	//	return err
	//}
	//
	//portUUID, err := env.controller.Ovs.VSwitch.GetPortUUID(fmt.Sprintf("%s_AllBlue", bridge))
	//if err != nil {
	//	log.Error().Err(err).Str("port", fmt.Sprintf("%s_AllBlue", bridge)).Msg("getting port uuid")
	//	return err
	//}
	//
	//var vlans []string
	//for _, network := range nets {
	//	vlans = append(vlans, fmt.Sprint(network.Tag))
	//}
	//
	//if err := env.controller.Ovs.VSwitch.MirrorAllVlans(mirror, portUUID, vlans); err != nil {
	//	log.Error().Err(err).Msgf("mirroring traffic")
	//	return err
	//}
	//
	return nil
}

func (env *environment) addTargetVM(ctx context.Context, bridge string) error {

	var special []string

	log.Info().Str("sandbox tag", bridge).Msg("creating special interface")
	if err := env.createPort(bridge, "special", 10); err != nil {
		log.Error().Err(err).Msg("creating mirror")

		return err
	}

	//TODO: Aici trebuie chestia aia cu TCPDUMP traffic
	//		pentru portul unde e masina compromisa

	dt := time.Now()
	dt.Format("01022006_150405_Mon")
	targetIntf := fmt.Sprintf("%s_special", bridge)
	log.Debug().Msg("ACUM urmeaza functia problema ")
	go func() {
		if err := env.controller.TCPdump.DumpTraffic(targetIntf, fmt.Sprintf("special_%s", dt.Format("01022006_150405_Mon"))); err != nil {
			log.Error().Err(err).Str("interface: ", targetIntf).Msg("problem starting tcpdump")

		}
	}()

	fmt.Println("este blocat aici sau doar dureaza mai mult?")
	log.Debug().Msg("Oare chiar ramane blocat aici ")
	////TODO: Vezi de ce moare
	///AICI MOARE, DE CE MORI ?? DE CE NU MERGI MAI DEPARTE?? DE CE??

	special = append(special, targetIntf)

	vm, err := env.vlib.GetCopy(ctx, bridge,
		vbox.InstanceConfig{Image: "pain3.ova",
			CPU:      2,
			MemoryMB: 4500},
		vbox.SetBridge(special, true),
		vbox.SetMAC("6CF0491A6E12", 1),
	)

	if err != nil {
		log.Error().Err(err).Msg("creating copy of SoC VM")
		return err
	}
	if vm == nil {
		return ErrVMNotCreated
	}
	log.Debug().Str("VM", vm.Info().Id).Msg("starting VM")

	if err := vm.Start(ctx); err != nil {
		log.Error().Err(err).Msgf("starting virtual machine")
		return err
	}
	env.instances = append(env.instances, vm)

	return nil
}

func (env *environment) initializeSOC(ctx context.Context, networks []string, mac string, tag string, nic int) error {

	vm, err := env.vlib.GetCopy(ctx, tag,
		vbox.InstanceConfig{Image: "socWorking.ova",
			CPU:      2,
			MemoryMB: 8096},
		//vbox.MapVMPort([]virtual.NatPortSettings{
		//	{
		//		HostPort:    strconv.FormatUint(uint64(socPort), 10),
		//		GuestPort:   "22",
		//		ServiceName: "sshd",
		//		Protocol:    "tcp",
		//	},
		//}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(networks, false),
		vbox.SetMAC("04D3B09BEAD6", nic),
	)

	if err != nil {
		log.Error().Err(err).Msg("creating copy of SoC VM")
		return err
	}
	if vm == nil {
		return ErrVMNotCreated
	}
	log.Debug().Str("VM", vm.Info().Id).Msg("starting VM")

	if err := vm.Start(ctx); err != nil {
		log.Error().Err(err).Msgf("starting virtual machine")
		return err
	}
	env.instances = append(env.instances, vm)

	return nil
}

func (env *environment) initOpnSenseVM(ctx context.Context, tag string, vlanPorts []string, mngtPort, routerPort uint) error {

	vm, err := env.vlib.GetCopy(ctx,
		tag,
		vbox.InstanceConfig{Image: "opnfuck.ova",
			CPU:      1,
			MemoryMB: 1600},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				// this is for management opnSense port
				HostPort:    strconv.FormatUint(uint64(mngtPort), 10),
				GuestPort:   "443",
				ServiceName: "mngtP",
				Protocol:    "tcp",
			},
			{
				HostPort:    strconv.FormatUint(uint64(routerPort), 10),
				GuestPort:   "22",
				ServiceName: "sshd",
				Protocol:    "tcp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(vlanPorts, false),
	)

	if err != nil {
		log.Error().Err(err).Msg("creating OpnSense VM")
		return err
	}
	if vm == nil {
		fmt.Println("Banuiesc ca aici e o problema :)")
		log.Debug().Str("VM", vm.Info().Id).Msg("starting VM")
		return ErrVMNotCreated
	}
	//log.Debug().Str("VM", vm.Info().Id).Msg("starting VM")

	if err := vm.Start(ctx); err != nil {
		log.Error().Err(err).Msgf("starting virtual machine")
		return err
	}
	env.instances = append(env.instances, vm)

	return nil
}
