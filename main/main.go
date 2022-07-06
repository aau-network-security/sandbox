package main

import (
	"context"
	"flag"
	"github.com/aau-network-security/sandbox/store"

	//"flag"
	"fmt"
	"github.com/aau-network-security/sandbox/config"
	"github.com/aau-network-security/sandbox/sandbox"
	//"github.com/aau-network-security/sandbox2/models"
	//"github.com/aau-network-security/sandbox2/sandbox2"
	//"github.com/aau-network-security/sandbox2/store"
	//"github.com/docker/docker/integration-cli/environment"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"os"
	//"time"
)

func main() {
	//TODO: add the target vm localtion and specific VM
	//var (
	//	tag        string
	//	vmsName    string
	//	targetVM   string
	//	networksNo int
	//)

	//	flag.StringVar(&tag, "tag", "test", "name of experiment")
	//	flag.StringVar(&vmsName, "vmsName", "ubuntu.ova", "name of the target machine")
	//	flag.StringVar(&targetVM, "targetVM", "ubuntu.ova", "name for rest of virtual machines")
	//	flag.IntVar(&networksNo, "networksNo", 3, "number of networks")
	//
	//	flag.Parse()
	//
	//	fmt.Println("tag:", tag)
	//	fmt.Println("networksNo:", networksNo)
	//	fmt.Println("vm:", vmsName)
	//	fmt.Println("vm:", targetVM)

	var defaultScenarioFile = "scenarios"

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir %v", err)
	}
	fullPathToConfig := fmt.Sprintf("%s%s", dir, "/config/config.yml")
	//TODO: rezolva problema cu config file
	configuration, err := config.NewConfig(fullPathToConfig)
	if err != nil {
		panic(err)
	}

	scenFilePtr := flag.String("scenarios", defaultScenarioFile, "scenario folder")
	flag.Parse()

	scenarios, err := store.LoadScenarios(*scenFilePtr)
	if err != nil {
		log.Error().Err(err).Str("file", *scenFilePtr).Msgf("failed to read scenarios from file")
		return
	}

	sandboxConf := sandbox.SandConfig{
		Name:   "sandbox2",
		Tag:    "test",
		Config: configuration,
	}

	sand, err := sandbox.NewSandbox(&sandboxConf)
	if err != nil {
		log.Error().Msg("Problem in creating the newSandbox")
	}

	fmt.Println("Acum urmeaza problema")

	if err := sand.StartSandbox(context.TODO(), sandboxConf.Tag, sandboxConf.Name, scenarios); err != nil {
		log.Error().Msg("Problem in starting the sandbox2")
	}

}
