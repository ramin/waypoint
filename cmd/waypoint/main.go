package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ramin/waypoint/commands"
	"github.com/ramin/waypoint/config"
)

var (
	address    string
	p2pNetwork string
	jwt        string
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(formatter())

	waypoint.Flags().StringVar(&address, "address", "", "Address to use")
	waypoint.Flags().StringVar(&p2pNetwork, "p2p.network", "", "P2P network to connect to")
	waypoint.Flags().StringVar(&jwt, "jwt", "", "JWT token")

	if jwt != "" {
		config.Read().JWT = jwt
	}

	level := os.Getenv("LOG_LEVEL")

	if level == "" {
		level = "info"
	}

	l, _ := log.ParseLevel(level)
	log.SetLevel(l)
}

var waypoint = &cobra.Command{
	Use:     "waypoint",
	Version: "0.0.1",
	Short:   "waypoint",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			panic(err)
		}
	},
}

func main() {
	waypoint.AddCommand(commands.StartCmd)
	waypoint.AddCommand(commands.RunCmd)

	if err := waypoint.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func formatter() log.Formatter {
	if format := os.Getenv("LOG_FORMAT"); format == "JSON" {
		return &log.JSONFormatter{}
	}

	return &log.TextFormatter{}
}
