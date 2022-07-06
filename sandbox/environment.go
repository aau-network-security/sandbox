package sandbox

import (
	"context"
	"errors"
	"fmt"
	"github.com/aau-network-security/sandbox/config"
	"github.com/aau-network-security/sandbox/controller"
	"github.com/aau-network-security/sandbox/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"sync"

	//"github.com/aau-network-security/sandbox2/models"
	"github.com/aau-network-security/sandbox/store"
	"github.com/aau-network-security/sandbox/virtual"
	"github.com/aau-network-security/sandbox/virtual/docker"
	"github.com/aau-network-security/sandbox/virtual/vbox"
	"github.com/rs/zerolog/log"
)

var (
	redListenPort  uint = 5181
	blueListenPort uint = 5182
	min                 = 7900
	max                 = 7950
	gmin                = 5350
	gmax                = 5375
	smin                = 3000
	smax                = 3500
	rmin                = 5000
	rmax                = 5300

	ErrVMNotCreated       = errors.New("no VM created")
	ErrGettingContainerID = errors.New("could not get container ID")
)

type environment struct {
	// challenge microservice should be integrated heres
	controller controller.NetController
	//wg         vpn.WireguardClient
	//dhcp       dhproto.DHCPClient
	dockerHost docker.Host
	instances  []virtual.Instance
	ports      []string
	vlib       vbox.Library
	//dnsServer  *dns.Server
}

type SandConfig struct {
	//ID       string
	Name string
	Tag  string
	//WgConfig      wg.WireGuardConfig
	//Host       string
	env    *environment
	Config *config.Config

	//NetworksIP map[string]string
	//NetworksNO int

	//redVPNIp      string
	//blueVPNIp     string
	//redPort       uint
	//bluePort      uint
	//CreatedAt     time.Time
	//RedPanicLeft  uint
	//BluePanicLeft uint
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

	log.Info().Str("Game Tag", tag).
		Str("Game Name", name).
		Str("Scenario", scenario.Name).
		Msg("starting sandbox")

	log.Debug().Str("Game", name).Str("bridgeName", tag).Msg("creating openvswitch bridge")
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

	log.Debug().Str("Game", name).Msg("configuring monitoring")
	if err := gc.env.configureMonitor(ctx, tag, scenario.Networks); err != nil {
		log.Error().Err(err).Msgf("configuring monitoring")
		return err
	}

	//log.Debug().Str("Game", tag).Msgf("Initilizing VPN VM")

	//assign connection port to RED users
	//redTeamVPNPort := getRandomPort(min, max)

	// assign grpc port to wg vms
	//wgPort := getRandomPort(gmin, gmax)

	//routerPort := getRandomPort(rmin, rmax)

	//assign connection port to Blue users
	//blueTeamVPNPort := getRandomPort(min, max)

	//if err := gc.env.initWireguardVM(ctx, tag, vlanPorts, redTeamVPNPort, blueTeamVPNPort, wgPort, routerPort); err != nil {
	//
	//	return err
	//}

	//log.Debug().Str("Game", name).Msg("waiting for wireguard vm to boot")
	//
	//dhcpClient, err := dhcp.NewDHCPClient(ctx, gc.WgConfig, wgPort)
	//if err != nil {
	//	log.Error().Err(err).Msg("connecting to DHCP service")
	//	return err
	//}
	//
	//gc.env.dhcp = dhcpClient

	//log.Debug().Str("Game  ", name).Msg("starting DHCP server")

	//gc.NetworksIP, err, ipMail, ipDC = gc.env.initDHCPServer(ctx, len(scenario.Networks), scenario)
	//if err != nil {
	//	return err
	//}

	//log.Debug().Str("Game  ", name).Msg("starting DNS server")
	//
	//if err := gc.env.initDNSServer(ctx, tag, gc.NetworksIP, scenario, ipMail, ipDC); err != nil {
	//	log.Error().Err(err).Msg("connecting to DHCP service")
	//	return err
	//}

	//wgClient, err := wg.NewGRPCVPNClient(ctx, gc.WgConfig, wgPort)
	//if err != nil {
	//	log.Error().Err(err).Msg("connecting to wireguard service")
	//	return err
	//}
	//gc.env.wg = wgClient

	log.Debug().Str("Game", name).Msg("initializing scenario")
	if err := gc.env.initializeScenario(ctx, tag, scenario); err != nil {
		return err
	}

	//ethInterfaceName := "eth0" // can be customized later

	//redTeamVPNIp, err := gc.env.getRandomIp()
	//if err != nil {
	//	log.Error().Err(err).Msg("Problem in generating red team VPNip")
	//	return err
	//}
	////
	//gc.redVPNIp = fmt.Sprintf("%s.0/24", redTeamVPNIp)
	////Assigning a connection port for Red team
	//
	//gc.redPort = redTeamVPNPort

	// create wireguard interface for red team
	//wgNICred := fmt.Sprintf("%s_red", tag)

	// initializing VPN endpoint for red team
	//if err := gc.env.initVPNInterface(gc.redVPNIp, redListenPort, wgNICred, ethInterfaceName); err != nil {
	//	return err
	//}
	//
	//blueTeamVPNIp, err := gc.env.getRandomIp()
	//if err != nil {
	//	log.Error().Err(err).Msg("")
	//	return err
	//}
	//
	//gc.blueVPNIp = fmt.Sprintf("%s.0/24", blueTeamVPNIp)
	//
	////Assigning a connection port for blue team
	//
	//gc.bluePort = blueTeamVPNPort
	// initializing VPN endpoint for blue team

	//create wireguard interface for blue team
	//wgNICblue := fmt.Sprintf("%s_blue", tag)

	//if err := gc.env.initVPNInterface(gc.blueVPNIp, blueListenPort, wgNICblue, ethInterfaceName); err != nil {
	//	return err
	//}

	macAddress := "04:d3:b0:9b:ea:d6"
	macAddressClean := strings.ReplaceAll(macAddress, ":", "")

	log.Debug().Str("sandbox", tag).Msg("Initalizing SoC")
	socPort := getRandomPort(smin, smax)
	ifaces := []string{fmt.Sprintf("%s_monitoring", tag), fmt.Sprintf("%s_AllBlue", tag)}

	//Todo: add also internet interface

	//ifaces :=
	if err := gc.env.initializeSOC(ctx, ifaces, macAddressClean, tag, 2, socPort); err != nil {
		log.Error().Err(err).Str("sandbox", tag).Msg("starting SoC vm")
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

//func (env *environment) initVPNInterface(ipAddress string, port uint, vpnInterfaceName, ethInterface string) error {
//
//	// ipAddress should be in this format : "45.11.23.1/24"
//	// port should be unique per interface
//
//	_, err := env.wg.InitializeI(context.Background(), &vpn.IReq{
//		Address:            ipAddress,
//		ListenPort:         uint32(port),
//		SaveConfig:         true,
//		Eth:                ethInterface,
//		IName:              vpnInterfaceName,
//		DownInterfacesFile: "/etc/network/downinterfaces",
//	})
//	if err != nil {
//		log.Error().Msgf("Error in initializing interface %v", err)
//		return err
//	}
//	return nil
//}

//func (env *environment) initDHCPServer(ctx context.Context, numberNetworks int, scenario store.Scenario) (map[string]string, error, string, string) {
//	var networks []*dhproto.Network
//	var staticHosts []*dhproto.StaticHost
//	var ipMail, ipDC string
//
//	ipList := make(map[string]string)
//
//	for i := 1; i <= numberNetworks; i++ {
//		var network dhproto.Network
//		randIP, _ := env.controller.IPPool.Get()
//		network.Network = randIP + ".0"
//		network.Min = randIP + ".6"
//		network.Max = randIP + ".250"
//		network.Router = randIP + ".1"
//
//		ipList[fmt.Sprintf("%d", 10*i)] = randIP + ".0/24"
//		network.DnsServer = randIP + ".2"
//		networks = append(networks, &network)
//
//	}
//
//	// Setup monitoring network
//
//	monitoringNet := dhproto.Network{
//		Network:   "10.10.10.0",
//		Min:       "10.10.10.6",
//		Max:       "10.10.10.199",
//		Router:    "10.10.10.1",
//		DnsServer: "10.10.10.2",
//	}
//	ipList[""] = "10.10.10.0/24"
//
//	networks = append(networks, &monitoringNet)
//	//Todo: This is scenario based method to make it work
//	// in future this needs to be scenario indenpent
//
//	for _, item := range scenario.Hosts {
//
//		//cast la string acum e lista de stringuri
//		if item.Name == "mailserver" {
//
//			ipMail = ConstructStaticIP(ipList, item.Networks, item.IPAddr)
//			host := dhproto.StaticHost{
//				Name: item.Name,
//
//				MacAddress: "04:d3:04:54:fe:15",
//				Address:    ipMail,
//				Router:     ConstructStaticIP(ipList, item.Networks, ".1"),
//				DomainName: fmt.Sprintf("\"%s\"", item.DNS),
//				DnsServer:  ConstructStaticIP(ipList, item.Networks, ".2"),
//			}
//
//			staticHosts = append(staticHosts, &host)
//			continue
//
//		} else if item.Name == "DCcon" {
//			fmt.Printf("Este in bucla cu DCcon \n")
//
//			ipDC = ConstructStaticIP(ipList, item.Networks, item.IPAddr)
//			fmt.Printf("DCcon IP: %s\n", ipDC)
//
//			host := dhproto.StaticHost{
//				Name:       item.Name,
//				MacAddress: "04:d3:b0:c7:57:c7",
//				Address:    ipDC,
//				Router:     ConstructStaticIP(ipList, item.Networks, ".1"),
//				DomainName: fmt.Sprintf("\"%s\"", item.DNS),
//				DnsServer:  ConstructStaticIP(ipList, item.Networks, ".2"),
//			}
//			staticHosts = append(staticHosts, &host)
//		} else {
//			fmt.Printf("Este in bucla cu Else. \n")
//			continue
//		}
//
//	}
//
//	host := dhproto.StaticHost{
//		Name:       "SOC",
//		MacAddress: "04:d3:b0:9b:ea:d6",
//		Address:    "10.10.10.200",
//		Router:     "10.10.10.1",
//		DomainName: "\"blue.monitor.soc\"",
//		DnsServer:  "10.10.10.2",
//	}
//
//	staticHosts = append(staticHosts, &host)
//
//	_, err := env.dhcp.StartDHCP(ctx, &dhproto.StartReq{Networks: networks, StaticHosts: staticHosts})
//	if err != nil {
//		return ipList, err, ipMail, ipDC
//	}
//
//	return ipList, nil, ipMail, ipDC
//}

//func (env *environment) initDNSServer(ctx context.Context, bridge string, ipList map[string]string, scenario store.Scenario, IPMail, IPdc string) error {
//
//	server, err := dns.New(bridge, ipList, scenario, IPMail, IPdc)
//	if err != nil {
//		log.Error().Msgf("Error creating DNS server %v", err)
//		return err
//	}
//	env.dnsServer = server
//	//env.instances = append(env.instances, server )
//
//	if err := server.Run(ctx); err != nil {
//		log.Error().Msgf("Error in starting DNS  %v", err)
//		return err
//	}
//
//	contID := server.Container().ID()
//	fmt.Printf("AICI e ID = %s\n", contID)
//
//	i := 1
//	for _, network := range ipList {
//
//		if network == "10.10.10.0/24" {
//
//			ipAddrs := strings.TrimSuffix(network, ".0/24")
//			ipAddrs = ipAddrs + ".2/24"
//
//			fmt.Println(ipAddrs)
//
//			if err := env.controller.Ovs.Docker.AddPort(bridge, fmt.Sprintf("eth%d", i), contID, ovs.DockerOptions{IPAddress: ipAddrs}); err != nil {
//
//				log.Error().Err(err).Str("container", contID).Msg("adding port to DNS container")
//				return err
//			}
//			i++
//			fmt.Println(i)
//
//		} else {
//			ipAddrs := strings.TrimSuffix(network, ".0/24")
//			ipAddrs = ipAddrs + ".2/24"
//
//			fmt.Println(ipAddrs)
//			//fmt.Sprintf("eth%d", vlan)
//			tag := i * 10
//
//			sTag := strconv.Itoa(tag)
//
//			fmt.Println(sTag)
//			if err := env.controller.Ovs.Docker.AddPort(bridge, fmt.Sprintf("eth%d", i), contID, ovs.DockerOptions{VlanTag: sTag, IPAddress: ipAddrs}); err != nil {
//
//				log.Error().Err(err).Str("container", contID).Msg("adding port to DNS container")
//				return err
//			}
//			i++
//			fmt.Println(i)
//
//		}
//
//	}
//
//	return nil
//}

//configureMonitor will configure the monitoring VM by attaching the correct interfaces
func (env *environment) configureMonitor(ctx context.Context, bridge string, nets []models.Network) error {

	log.Info().Str("sandbox tag", bridge).Msg("creating monitoring network")
	if err := env.createPort(bridge, "monitoring", 0); err != nil {
		return err
	}

	mirror := fmt.Sprintf("%s_mirror", bridge)

	log.Info().Str("sandbox tag", bridge).Msg("Creating the network mirror")
	if err := env.controller.Ovs.VSwitch.CreateMirrorforBridge(mirror, bridge); err != nil {
		log.Error().Err(err).Msg("creating mirror")
		return err
	}

	if err := env.createPort(bridge, "AllBlue", 0); err != nil {
		return err
	}

	portUUID, err := env.controller.Ovs.VSwitch.GetPortUUID(fmt.Sprintf("%s_AllBlue", bridge))
	if err != nil {
		log.Error().Err(err).Str("port", fmt.Sprintf("%s_AllBlue", bridge)).Msg("getting port uuid")
		return err
	}

	var vlans []string
	for _, network := range nets {
		vlans = append(vlans, fmt.Sprint(network.Tag))
	}

	if err := env.controller.Ovs.VSwitch.MirrorAllVlans(mirror, portUUID, vlans); err != nil {
		log.Error().Err(err).Msgf("mirroring traffic")
		return err
	}

	return nil
}

func (env *environment) initializeSOC(ctx context.Context, networks []string, mac string, tag string, nic int, socPort uint) error {

	vm, err := env.vlib.GetCopy(ctx, tag,
		vbox.InstanceConfig{Image: "soc2022.ova",
			CPU:      4,
			MemoryMB: 32384},
		vbox.MapVMPort([]virtual.NatPortSettings{
			{
				HostPort:    strconv.FormatUint(uint64(socPort), 10),
				GuestPort:   "22",
				ServiceName: "sshd",
				Protocol:    "tcp",
			},
		}),
		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
		// is attaching to openvswitch network
		vbox.SetBridge(networks, false),
		vbox.SetMAC(mac, nic),
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

//func (env *environment) initWireguardVM(ctx context.Context, tag string, vlanPorts []string, redTeamVPNport, blueTeamVPNport, wgPort uint, routerPort uint) error {
//
//	vm, err := env.vlib.GetCopy(ctx,
//		tag,
//		vbox.InstanceConfig{Image: "Routerfix.ova",
//			CPU:      2,
//			MemoryMB: 2048},
//		vbox.MapVMPort([]virtual.NatPortSettings{
//			{
//				// this is for gRPC service
//				HostPort:    strconv.FormatUint(uint64(wgPort), 10),
//				GuestPort:   "5353",
//				ServiceName: "wgservice",
//				Protocol:    "tcp",
//			},
//			{
//				HostPort:    strconv.FormatUint(uint64(redTeamVPNport), 10),
//				GuestPort:   strconv.FormatUint(uint64(redListenPort), 10),
//				ServiceName: "wgRedConnection",
//				Protocol:    "udp",
//			},
//			{
//				HostPort:    strconv.FormatUint(uint64(blueTeamVPNport), 10),
//				GuestPort:   strconv.FormatUint(uint64(blueListenPort), 10),
//				ServiceName: "wgBlueConnection",
//				Protocol:    "udp",
//			},
//			{
//				HostPort:    strconv.FormatUint(uint64(routerPort), 10),
//				GuestPort:   "22",
//				ServiceName: "sshd",
//				Protocol:    "tcp",
//			},
//		}),
//		// SetBridge parameter cleanFirst should be enabled when wireguard/router instance
//		// is attaching to openvswitch network
//		vbox.SetBridge(vlanPorts, false),
//	)
//
//	if err != nil {
//		log.Error().Err(err).Msg("creating VPN VM")
//		return err
//	}
//	if vm == nil {
//		return ErrVMNotCreated
//	}
//	log.Debug().Str("VM", vm.Info().Id).Msg("starting VM")
//
//	if err := vm.Start(ctx); err != nil {
//		log.Error().Err(err).Msgf("starting virtual machine")
//		return err
//	}
//	env.instances = append(env.instances, vm)
//
//	return nil
//}

//func (gc *SandConfig) CreateVPNConfig(ctx context.Context, isRed bool, idUser string) (VPNConfig, error) {
//
//	var nicName string
//
//	var allowedIps []string
//	var peerIP string
//	var endpoint string
//	//var dns string
//	if isRed {
//		//dns = ""
//		nicName = fmt.Sprintf("%s_red", gc.Tag)
//
//		for key := range gc.NetworksIP {
//			if gc.NetworksIP[key] == "10.10.10.0/24" {
//				continue
//			}
//			allowedIps = append(allowedIps, gc.NetworksIP[key])
//			break
//		}
//
//		peerIP = gc.redVPNIp
//		allowedIps = append(allowedIps, peerIP)
//
//		endpoint = fmt.Sprintf("%s.%s:%d", gc.Tag, gc.Host, gc.redPort)
//	} else {
//
//		nicName = fmt.Sprintf("%s_blue", gc.Tag)
//		for key := range gc.NetworksIP {
//			allowedIps = append(allowedIps, gc.NetworksIP[key])
//		}
//
//		peerIP = gc.blueVPNIp
//		allowedIps = append(allowedIps, peerIP)
//		endpoint = fmt.Sprintf("%s.%s:%d", gc.Tag, gc.Host, gc.bluePort)
//	}
//
//	serverPubKey, err := gc.env.wg.GetPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: nicName, PrivKeyName: nicName})
//	if err != nil {
//		log.Error().Err(err).Str("User", idUser).Msg("Err get public nicName wireguard")
//		return VPNConfig{}, err
//	}
//
//	_, err = gc.env.wg.GenPrivateKey(ctx, &vpn.PrivKeyReq{PrivateKeyName: gc.Tag + "_" + idUser + "_"})
//	if err != nil {
//		//fmt.Printf("Err gen private nicName wireguard  %v", err)
//		log.Error().Err(err).Str("User", idUser).Msg("Err gen private nicName wireguard")
//		return VPNConfig{}, err
//	}
//
//	//generate client public nicName
//	//log.Info().Msgf("Generating public nicName for team %s", evTag+"_"+team+"_"+strconv.Itoa(ipAddr))
//	_, err = gc.env.wg.GenPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: gc.Tag + "_" + idUser + "_", PrivKeyName: gc.Tag + "_" + idUser + "_"})
//	if err != nil {
//		log.Error().Err(err).Str("User", idUser).Msg("Err gen public nicName client")
//		return VPNConfig{}, err
//	}
//
//	clientPubKey, err := gc.env.wg.GetPublicKey(ctx, &vpn.PubKeyReq{PubKeyName: gc.Tag + "_" + idUser + "_"})
//	if err != nil {
//		fmt.Printf("Error on GetPublicKey %v", err)
//		return VPNConfig{}, err
//	}
//
//	pIP := fmt.Sprintf("%d/32", IPcounter())
//
//	peerIP = strings.Replace(peerIP, "0/24", pIP, 1)
//
//	_, err = gc.env.wg.AddPeer(ctx, &vpn.AddPReq{
//		Nic:        nicName,
//		AllowedIPs: peerIP,
//		PublicKey:  clientPubKey.Message,
//	})
//
//	if err != nil {
//		log.Error().Err(err).Msg("Error on adding peer to interface")
//		return VPNConfig{}, err
//
//	}
//
//	clientPrivKey, err := gc.env.wg.GetPrivateKey(ctx, &vpn.PrivKeyReq{PrivateKeyName: gc.Tag + "_" + idUser + "_"})
//	if err != nil {
//		log.Error().Err(err).Msg("getting priv NIC")
//		return VPNConfig{}, err
//	}
//
//	return VPNConfig{
//		ServerPublicKey:  serverPubKey.Message,
//		PrivateKeyClient: clientPrivKey.Message,
//		Endpoint:         endpoint,
//		AllowedIPs:       strings.Join(allowedIps, ", "),
//		PeerIP:           peerIP,
//	}, nil
//
//}