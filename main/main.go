package main

import (
	"flag"
	"fmt"
	"github.com/aau-network-security/sandbox/config"
	"github.com/aau-network-security/sandbox/sandbox"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {

	var (
		tag        string
		vmsName    string
		targetVM   string
		networksNo int
	)

	flag.StringVar(&tag, "tag", "test", "name of experiment")
	flag.StringVar(&vmsName, "vmsName", "ubuntu.ova", "name of the target machine")
	flag.StringVar(&targetVM, "targetVM", "ubuntu.ova", "name for rest of virtual machines")
	flag.IntVar(&networksNo, "networksNo", 3, "number of networks")

	flag.Parse()

	fmt.Println("tag:", tag)
	fmt.Println("networksNo:", networksNo)
	fmt.Println("vm:", vmsName)
	fmt.Println("vm:", targetVM)

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir %v", err)
	}
	fullPathToConfig := fmt.Sprintf("%s%s", dir, "/config/config.yml")

	configuration, err := config.NewConfig(fullPathToConfig)
	if err != nil {
		panic(err)
	}

	StartSandbox(tag, vmsName, targetVM, networksNo, configuration)

}

func StartSandbox(tag, vmsName, targetVM string, NetworksNO int, config *config.Config) error {
	//wgConfig := d.config.WireguardService

	env, err := sandbox.NewSandbox(sandbox.SandConfig{
		NetworksNO: NetworksNO,
		VmName:     vmsName,
		Tag:        tag,
	}, config.VmConfig.OvaDir)
	if err != nil {
		return err
	}

	if err := env.CreateSandbox(tag, targetVM, vmsName, NetworksNO); err != nil {
		log.Info().Err(err).Msgf("Sandbox environment is starting")
		return err
	}

	return nil

}
