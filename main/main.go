package main

import (
	"flag"
	"fmt"
	"github.com/aau-network-security/sandbox/config"
	"github.com/aau-network-security/sandbox/sandbox"
	"github.com/rs/zerolog/log"
)

func main() {

	var (
		tag        string
		vmName     string
		networksNo int
	)

	flag.StringVar(&tag, "tag", "test", "name of experiment")
	flag.StringVar(&vmName, "vmName", "ubuntu.ova", "name of the virtual machine")
	flag.IntVar(&networksNo, "networksNo", 3, "number of networks")

	flag.Parse()

	fmt.Println("tag:", tag)
	fmt.Println("networksNo:", networksNo)
	fmt.Println("vm:", vmName)

	configuration, err := config.NewConfig("/config/config.yml")
	if err != nil {
		panic(err)
	}

	StartSandbox(tag, vmName, networksNo, configuration)

}

func StartSandbox(tag, VmName string, NetworksNO int, config *config.Config) error {
	//wgConfig := d.config.WireguardService

	env, err := sandbox.NewSandbox(sandbox.SandConfig{
		NetworksNO: NetworksNO,
		VmName:     VmName,
		Tag:        tag,
	}, config.VmConfig.OvaDir)
	if err != nil {
		return err
	}

	if err := env.CreateSandbox(tag, VmName, NetworksNO); err != nil {
		log.Info().Err(err).Msgf("Sandbox environment is starting")
		return err
	}

	return nil

}
