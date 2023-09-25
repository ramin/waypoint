package commands

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
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

		logrus.Info(header)

		<-startServer(context.Background())
	},
}
