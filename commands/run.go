package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("waypoint: run")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		rpc, err := newClient(ctx)
		if err != nil {
			panic(err)
		}

		header, err := rpc.Header.GetByHeight(context.Background(), 20)
		if err != nil {
			panic(err)
		}

		fmt.Println(header)

		<-startServer(context.Background())
	},
}
