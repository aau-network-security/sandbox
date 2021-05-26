package sandbox

import (
	"context"
	"fmt"
	"github.com/aau-network-security/sandbox/dnet/dns"
	"io"
	"net"
	"strings"

	"github.com/aau-network-security/openvswitch/ovs"

	"github.com/aau-network-security/sandbox/controller"
	"github.com/aau-network-security/sandbox/dnet/dhcp"

	"github.com/aau-network-security/sandbox/virtual/docker"
	"github.com/aau-network-security/sandbox/virtual/vbox"
	"github.com/rs/zerolog/log"
)

const (
	MAX_NET_CONN = 7
)

//var (
//	challengeURLList = map[string]string{
//		"ftp":      "registry.gitlab.com/haaukins/forensics/ftp_bf_login",
//		"hb":       "registry.gitlab.com/haaukins/web-exploitation/heartbleed",
//		"microcms": "registry.gitlab.com/haaukins/web-exploitation/micro_cms",
//		"scan":     "registry.gitlab.com/haaukins/forensics/hidden-server",
//		"rot":      "registry.gitlab.com/haaukins/crytopgraphy/rot13",
//		"csrf":     "registry.gitlab.com/haaukins/web-exploitation/csrf",
//		"uwb":      "registry.gitlab.com/haaukins/web-exploitation/webadmin-1.920-urce",
//	}
//	TemporaryScenariosPlaceHolder = map[int]Scenario{
//		1: {
//			ID: 1,
//			Networks: []network{
//				{
//					Chals: []string{"hb", "ftp", "scan"},
//					Vlan:  "vlan20",
//				},
//				{
//					Chals: []string{"scan", "csrf"},
//					Vlan:  "vlan30",
//				},
//				{
//					Chals: []string{"rot", "uwb"},
//					Vlan:  "vlan10",
//				},
//			},
//			Story:    "Scenario 1 Story placeholder",
//			Duration: "2",
//
//			Difficulty: "Easy",
//		},
//		2: {
//			ID: 2,
//			Networks: []network{
//				{
//					Chals: []string{"microcms", "joomla", "uwb"},
//					Vlan:  "vlan10",
//				},
//				{
//					Chals: []string{"jwt", "csrf"},
//					Vlan:  "vlan20",
//				},
//				{
//					Chals: []string{"rot", "uwb"},
//					Vlan:  "vlan40",
//				},
//				{
//					Chals: []string{"rot", "uwb"},
//					Vlan:  "vlan30",
//				},
//			},
//			Story:      "Scenario 2 Story placeholder",
//			Duration:   "3",
//			Difficulty: "Moderate",
//		},
//	}
//)

type Environment interface {
	GetScenarios() string
}

type network struct {
	Vlan  string
	Chals []string
}

type Scenario struct {
	ID         int
	Networks   []network
	Story      string
	Duration   string
	Difficulty string
}

type environment struct {
	// web interface microservice should stay here
	// challenge microservice should be integrated heres
	controller controller.NetController
	//wg         vpn.WireguardClient
	dockerHost docker.Host
	closer     io.Closer
	config     SandConfig
	vlib       vbox.Library
	dhcp       *dhcp.Server
	dns        *dns.Server
}

type SandConfig struct {
	NetworksNO int
	VmName     string
	Tag        string
	//WgConfig   wg.WireGuardConfig
}

func NewSandbox(conf SandConfig, vboxConf string) (*environment, error) {
	//if len(conf.Scenario.Networks) > MAX_NET_CONN {
	//	return nil, fmt.Errorf("exceeds maximum number of Networks for a environment. Max is %d", MAX_NET_CONN)
	//}
	//
	//wgClient, err := wg.NewGRPCVPNClient(conf.WgConfig)
	//if err != nil {
	//	log.Error().Msgf("Connection error on wireguard service error %v ", err)
	//	return nil, err
	//}
	netController := controller.New()
	vlib := vbox.NewLibrary(vboxConf)
	if vlib == nil {
		log.Error().Msgf("Library could not be created properly...")
		return nil, fmt.Errorf("Error on new library")
	}
	dockerHost := docker.NewHost()
	env := &environment{
		controller: *netController,
		//wg:         wgClient,
		dockerHost: dockerHost,
		config:     conf,
		vlib:       vlib,
	}
	log.Info().Msgf("New environment initialized ")

	return env, nil

}

func (g *environment) Close() error {
	//ch := make(chan interface{})
	//var wg sync.WaitGroup
	//var closers []io.Closer

	//if g.dhcp != nil {
	//	closers = append(closers, g.dhcp)
	//}

	// todo: add closers for other components as well
	//closer := func() {
	//	select {
	//		case <-ch:
	//			g.dhcp.Close()
	//
	//
	//
	//	}
	//
	//}()
	//
	//go closer()
	//if g.dockerHost != nil {
	//	close(ch)
	//}

	return nil
}

func (g *environment) CreateSandbox(tag, targetVM, VMname string, numNetworks int) error {
	log.Info().Str("Experiment name ", tag).
		Int("Number of Networks", numNetworks).
		Msgf("Creating sandbox")
	// bridge name will be same with event tag
	bridgeName := tag
	//selectedScenario := TemporaryScenariosPlaceHolder[scenarioNo]
	//numNetworks := len(selectedScenario.Networks)
	log.Info().Msgf("Setting openvswitch bridge %s", bridgeName)
	if err := g.initializeOVSBridge(bridgeName); err != nil {
		log.Error().Err(err).Msg("Error with initializing OVSBridge")
		return err
	}
	if err := g.createRandomNetworks(bridgeName, int(numNetworks)); err != nil {
		log.Error().Err(err).Msgf("Error on creating random networks")

	}

	var availableVlans []string
	ineti, err := net.Interfaces()
	if err != nil {
		log.Error().Err(err).Msgf("Error getting the system interfaces")
		panic(err)

	}

	for _, inter := range ineti {
		if strings.Contains(inter.Name, "vlan") {
			availableVlans = append(availableVlans, inter.Name)

		}

	}

	for _, vlans := range availableVlans {
		if strings.Contains(vlans, "vlan10") {
			if err := g.createTargetVM(vlans, targetVM); err != nil {
				//log.Error().Strs("%s", vlans).Msgf("Error in initializing Wireguard")
				log.Error().Err(err).Msgf("Error in attaching targetVM to correct network")

			}
			continue

		}

		if err := g.populateNetworks(vlans, VMname); err != nil {

			log.Error().Err(err).Msgf("Error in populating the networks with vms")

		}

	}

	return nil

}

func (g *environment) createRandomNetworks(bridge string, numberOfNetworks int) error {
	vlanTags := make(map[string]string)

	var vlansList []string

	log.Info().Msgf("Creating randomized Networks for chosen number of Networks %d", numberOfNetworks)
	//Always creating +1 network for the monitoring machine.
	for i := 1; i <= numberOfNetworks; i++ {
		vlan := fmt.Sprintf("vlan%d", i*10)
		vlansList = append(vlansList, vlan)
		vlanTags[vlan] = fmt.Sprintf("%d", i*10)
		if err := g.controller.Ovs.VSwitch.AddPortTagged(bridge, vlan, fmt.Sprintf("%d", i*10)); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
			return err
		}
		log.Info().Msgf("AddPort Set Interface Options %s", vlan)
		if err := g.controller.Ovs.VSwitch.Set.Interface(vlan, ovs.InterfaceOptions{Type: ovs.InterfaceTypeInternal}); err != nil {
			log.Error().Msgf("Error on matching interface error %v", err)
			return err
		}

		//ip tuntap add tap0 mode tap
		//ifconfig tap0 up
		//ip tuntap add tap2 mode tap
		// ifconfig tap2 up
		//ip tuntap add tap4 mode tap
		//ifconfig tap4 up
		t := fmt.Sprintf("tap%d", i)
		if err := g.controller.IPService.AddTunTap(t, "tap"); err != nil {
			log.Error().Msgf("Error happened on adding tuntap %v", err)
			return err
		}
		if err := g.controller.IFConfig.TapUp(t); err != nil {
			log.Error().Msgf("Error happened on making up tap %s %v", t, err)
			return err
		}

		tag := fmt.Sprintf("%d", i*10)
		//ovs-vsctl add-port SW tap0 tag=10
		//ovs-vsctl add-port SW tap2 tag=20
		//ovs-vsctl add-port SW tap4 tag=30
		if err := g.controller.Ovs.VSwitch.AddPortTagged(bridge, t, tag); err != nil {
			log.Error().Msgf("Error on adding port with tag err %v", err)
			return err
		}

		if err := g.controller.IFConfig.TapUp(vlan); err != nil {
			log.Error().Msgf("Error happened on making up tap %s %v", vlan, err)
			return err
		}
	}

	log.Info().Msgf("Creating the monitoring network")
	//Always creating +1 network for the monitoring machine.

	//TODO: Make assign the monitoring network smarter ! Now is hardcoded.

	//How it is happening now will be a problem for multiple games
	i := 1

	monitor := fmt.Sprintf("mon%d", i*10)

	if err := g.controller.Ovs.VSwitch.AddPort(bridge, monitor); err != nil {
		log.Error().Msgf("Error on adding port with tag err %v", err)
		return err
	}

	m := fmt.Sprintf("mon%d", i*10)
	if err := g.controller.IPService.AddTunTap(m, "tap"); err != nil {
		log.Error().Msgf("Error happened on adding monitor tuntap %v", err)
		return err
	}
	if err := g.controller.IFConfig.TapUp(m); err != nil {
		log.Error().Msgf("Error happened on making up monitor %s %v", m, err)
		return err
	}
	//adding the monitoring port in the networks
	vlanTags["monitor"] = ""

	server, err := dhcp.New(context.TODO(), vlanTags, bridge, &g.controller)
	if err != nil {
		log.Error().Msgf("Error creating DHCP server %v", err)
		return err
	}
	if err := server.Run(context.Background()); err != nil {
		log.Error().Msgf("Error in starting DHCP  %v", err)
		return err
	}
	g.dhcp = server

	return nil

}

func (g *environment) initializeOVSBridge(bridgeName string) error {
	log.Info().Msgf("Game brigde name is set to game tag %s", bridgeName)
	if err := g.controller.Ovs.VSwitch.AddBridge(bridgeName); err != nil {
		log.Error().Msgf("Error on creating OVS bridge %v", err)
		return err
	}
	return nil
}

//func (g *environment) attachChallenge(bridge string, challengeList []string, cli *controller.NetController, vlan string) error {
//	ctx := context.Background()
//	log.Info().Msgf("Starting challenges for the game %s", bridge)
//	for _, ch := range challengeList {
//		container := docker.NewContainer(docker.ContainerConfig{
//			Image: challengeURLList[ch],
//			Labels: map[string]string{
//				"nap": "challenges",
//			}})
//
//		container.ID()
//		if err := container.Create(ctx); err != nil {
//			log.Error().Msgf("Error in creating container  %v", err)
//			return err
//		}
//		if err := container.Start(ctx); err != nil {
//			log.Error().Msgf("Error in creating container  %v", err)
//			return err
//		}
//
//		cid := container.ID()
//		if cid == "" {
//			return fmt.Errorf("Container ID could be fetched correctly")
//		}
//
//		if err := cli.Ovs.Docker.AddPort(bridge, "eth0", cid, ovs.DockerOptions{DHCP: true, VlanTag: vlan}); err != nil {
//			log.Error().Msgf("Error on adding port on docker %v", err)
//			return err
//		}
//
//		//if err := cli.Ovs.Docker.SetVlan(bridge, "eth0", cid, vlan); err != nil {
//		//	log.Error().Msgf("Error on ovs-docker SetVlan %v", err)
//		//	return err
//		//}
//
//	}
//
//	return nil
//
//}
//
//func (g *environment) initializeScenarios(bridge string, cli *controller.NetController, scenarioNumber int) error {
//	log.Debug().Msgf("Inializing scenarios for game [ %s ]", bridge)
//	networks := TemporaryScenariosPlaceHolder[scenarioNumber].Networks
//	var vlans []string
//	if scenarioNumber > 3 || scenarioNumber < 0 {
//		return fmt.Errorf("Invalid senario selection, make a selection between 1 to 3 ")
//	}
//	for _, net := range networks {
//		vlans = append(vlans, net.Vlan)
//
//	}
//	log.Debug().Strs("Network Vlans", vlans).Msgf("Vlans")
//
//
//	// initializing SOC
//
//
//
//	// initializing scenarios by attaching correct challenge to correct network
//	for _, net := range networks {
//		if err := g.attachChallenge(bridge, net.Chals, cli, net.Vlan[len(net.Vlan)-2:]); err != nil {
//
//			fmt.Printf("Error in attach challenge %v", err)
//			return err
//		}
//
//	}
//
//	return nil
//}

//configureMonitor will configure the monitoring VM by attaching the correct interfaces
//func (g *environment) configureMonitor(bridge string, numberNetworks int) error {
//
//	var ifaces []string
//	var vlanTags []string
//	var getBlue string  // mirrorName
//	var bluePort string // port in OVS for mirror traffic
//
//	getBlue = "blueMirror"
//	if err := g.controller.Ovs.VSwitch.CreateMirrorforBridge(getBlue, bridge); err != nil {
//		log.Error().Err(err).Msgf("Error on creating mirror")
//		return err
//
//	}
//
//	for i := 1; i <= numberNetworks; i++ {
//		tag := fmt.Sprintf("%d", i*10)
//		vlanTags = append(vlanTags, tag)
//
//	}
//
//	bluePort = "ALLblue"
//
//	if err := g.controller.IPService.AddTunTap(bluePort, "tap"); err != nil {
//		log.Error().Msgf("Error happened on adding monitor tuntap %v", err)
//		return err
//	}
//	if err := g.controller.IFConfig.TapUp(bluePort); err != nil {
//		log.Error().Msgf("Error happened on making up monitor %s %v", bluePort, err)
//		return err
//	}
//
//	if err := g.controller.Ovs.VSwitch.AddPort(bridge, bluePort); err != nil {
//		log.Error().Err(err).Msgf("Error on adding port to mirror traffic, err %v", err)
//		return err
//	}
//	//
//	//log.Info().Msgf("AddPort for mirroring Set Interface Options %s", bluePort)
//	//if err := g.controller.Ovs.VSwitch.Set.Interface(bluePort, ovs.InterfaceOptions{Type: ovs.InterfaceTypeInternal}); err != nil {
//	//	log.Error().Msgf("Error on matching interface error %v", err)
//	//	return err
//	//}
//
//	portUUID, err := g.controller.Ovs.VSwitch.GetPortUUID(bluePort)
//	if err != nil {
//		log.Error().Err(err).Msgf("Error on getting port uuid")
//		return err
//	}
//
//	if err := g.controller.Ovs.VSwitch.MirrorAllVlans(getBlue, portUUID, vlanTags); err != nil {
//		log.Error().Err(err).Msgf("Error on adding port to mirror traffic")
//		return err
//
//	}
//
//	//if err := g.controller.Ovs.VSwitch.MirrorAllVlans(getBlue, bluePort, vlanTags); err != nil {
//	//	log.Error().Err(err).Msgf("Error on adding port to mirror traffic")
//	//	return err
//	//
//	//}
//
//	//err, monitoringNetwr
//	ifaces = append(ifaces, bluePort)
//
//	ineti, err := net.Interfaces()
//	if err != nil {
//		log.Error().Err(err).Msgf("Error getting the system interfaces")
//		panic(err)
//
//	}
//
//	for _, inter := range ineti {
//		if strings.Contains(inter.Name, "mon") {
//			ifaces = append(ifaces, inter.Name)
//			if len(ifaces) != 2 {
//				log.Error().Err(err).Msgf("error on creating the list of interfaces")
//
//			}
//
//		}
//		continue
//
//	}
//
//	macAddress := g.dhcp.GetMAC()
//	macAddressClean := strings.ReplaceAll(macAddress, ":", "")
//	nicNumber := len(ifaces) + 1
//
//	fmt.Println(macAddressClean)
//	fmt.Println(nicNumber)
//	//
//	fmt.Println(ifaces)
//	if err := g.initializeSOC(ifaces, macAddressClean, nicNumber); err != nil {
//		log.Error().Err(err).Msgf("error starting VM with given interfaces")
//		return err
//	}
//
//	return nil
//
//}

//Attach one windows HostOnly Machine

//todo: 1 target machine with HOSTOnly

// the rest of the Devies in Bridged, Preferably some windows.

func (env *environment) createTargetVM(network string, targetVM string) error {
	log.Debug().Str("On vlan ", network).Msgf("is the target machine")
	vm, err := env.vlib.GetCopy(context.Background(),
		vbox.InstanceConfig{Image: targetVM,
			CPU:      1,
			MemoryMB: 2048},
		// SetBridge parameter cleanFirst should be enabled if you want to delete all interfaces first
		// is attaching to openvswitch network
		vbox.SetHostOnly(true),
		vbox.SetBridge(network, false),
		//vbox.SetNameofVM(),

	)

	if err != nil {
		log.Error().Msgf("Error while getting copy of VM err : %v", err)
		return err
	}
	if vm != nil {

		//log.Debug().Msgf("VM needs to be started now.. ", vm.Info().Id)

		log.Debug().Msgf("VM [ %s ] is starting .... ", vm.Info().Id)

		if err := vm.Start(context.Background()); err != nil {
			log.Error().Msgf("Failed to start target machine")
			return err
		}
	}
	return nil
}

func (env *environment) populateNetworks(network string, vmName string) error {

	//log.Debug().Str("Elastic Port", "9200").
	//	Str("Kibana Port", "5601").
	//	Msgf("Initalizing SoC for the game")

	// todo: Solve problem with the soc ovaFile
	vm, err := env.vlib.GetCopy(context.Background(),
		vbox.InstanceConfig{Image: vmName,
			CPU:      1,
			MemoryMB: 2048},
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network

		vbox.SetBridge(network, true),
		//vbox.SetMAC(mac, nic),
	)

	if err != nil {
		log.Error().Msgf("Error while getting copy of VM err : %v", err)
		return err
	}
	if vm != nil {

		//log.Debug().Msgf("VM needs to be started now.. ", vm.Info().Id)

		log.Debug().Msgf("VM [ %s ] is starting .... ", vm.Info().Id)

		if err := vm.Start(context.Background()); err != nil {
			log.Error().Msgf("Failed to start virtual machine on Vlan ")
			return err
		}
	}
	return nil

}
