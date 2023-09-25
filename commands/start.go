package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ramin/waypoint/config"
	"github.com/ramin/waypoint/server"
	"github.com/ramin/waypoint/verifier"

	"github.com/celestiaorg/celestia-node/api/rpc/client"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("waypoint")
		cfg := config.Read()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, err := verifier.NewVerifier(ctx)
		if err != nil {
			panic(err)
		}

		rpc, err := newClient(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(cfg.JWT)

		v.AddClient(rpc)

		go v.Start(ctx)

		<-startServer(ctx)
	},
}

func newClient(ctx context.Context) (*client.Client, error) {
	cfg := config.Read()
	fmt.Println(cfg.JWT)

	return client.NewClient(
		ctx,
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		cfg.JWT,
	)
}

func startServer(ctx context.Context) chan bool {
	srv := server.NewServer()

	signal.Notify(srv.Close, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s := <-srv.Close
		logrus.Infof("received signal %v", s)
		srv.GracefulShutdown()
	}()

	return srv.Start()
}
