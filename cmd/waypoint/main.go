package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ramin/waypoint/commands"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(formatter())

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
