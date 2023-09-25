package commands

import (
	"context"

	"github.com/celestiaorg/celestia-node/state"
	"github.com/ramin/waypoint/config"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command{
	Use: "info",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		defer ctx.Done()

		rpc, err := newClient(ctx)
		if err != nil {
			panic(err)
		}

		info, err := rpc.Node.Info(context.Background())
		if err != nil {
			panic(err)
		}

		logrus.Info(info.Type)
		logrus.Info(info.APIVersion)

		var address state.Address
		err = address.UnmarshalJSON([]byte(config.Read().Address))
		if err != nil {
			logrus.Error(err)
			return
		}

		balance, err := rpc.State.BalanceForAddress(context.Background(), address)
		if err != nil {
			panic(err)
		}
		logrus.Info("Got balance: %v", balance)

	},
}
