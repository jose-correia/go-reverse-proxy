package main

import (
	"context"
	"flag"
	"fmt"
	"go-reverse-proxy/app/common/log"
	config "go-reverse-proxy/app/handlers/configuration"
	"os"
)

var (
	configurationFileDirectory = "../../%s"
)

func main() {
	fs := flag.NewFlagSet("api", flag.ExitOnError)
	var (
		configFilename = fs.String("configuration_filename", "proxyConfig.yaml", "Name of the reverse proxy .yaml configuratio file")
	)
	_ = fs.Parse(os.Args[1:])

	ctx := context.Background()

	logger := log.NewLogger()

	// build the configuration .yaml filepath
	configFilepath := fmt.Sprintf(
		configurationFileDirectory,
		*configFilename,
	)

	configurationHandler := config.New(logger)

	// parse the .yaml configuration file into our data structure
	configuration, err := configurationHandler.FromFile(
		ctx, configFilepath,
	)
	if err != nil {
		logger.Log("module", "main", "error", err)
		os.Exit(1)
	}

	fmt.Println(configuration.Services[0].Domain)
}
