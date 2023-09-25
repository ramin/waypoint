package commands

import (
	"context"
	"fmt"

	"github.com/celestiaorg/celestia-node/state"
	"github.com/ramin/waypoint/config"

	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command{
	Use: "info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("waypoint")
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

		fmt.Println(info.Type)
		fmt.Println(info.APIVersion)

		var address state.Address
		err = address.UnmarshalJSON([]byte(config.Read().Address))
		if err != nil {
			fmt.Println(err)

			return
		}

		balance, err := rpc.State.BalanceForAddress(context.Background(), address)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got balance: %v", balance)

	},
}
