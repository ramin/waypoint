package verifier

import (
	"context"
	"fmt"
	"time"

	"github.com/celestiaorg/celestia-node/blob"
	"github.com/ramin/waypoint/generator"
	"github.com/sirupsen/logrus"
)

func (v *Verifier) StartWriter(ctx context.Context) {
	go func() {
		for {
			select {
			default:
				for height := range v.AwaitBlock(ctx) {
					go v.WriteToBlock(ctx, height)
				}
			case <-v.sig:
				fmt.Println("Received shutdown signal")
				close(v.done)
				return
			case <-ctx.Done():
				fmt.Println("Context canceled")
				close(v.done)
				return
			}
		}
	}()
}

func (v *Verifier) AwaitBlock(ctx context.Context) chan int64 {
	fmt.Println("awaiting block and writing new data")

	time.Sleep(time.Second * 2)

	heights := make(chan int64)

	go func() {
		for {
			heights <- 0

			time.Sleep(time.Second * 10)
		}
	}()

	return heights
}

func (v *Verifier) WriteToBlock(ctx context.Context, height int64) int64 {
	fmt.Println("awaiting block and writing new data")

	writeBlob, err := generator.NewBlob()
	if err != nil {
		v.errCh <- err
		v.Metrics.Errors.Add(ctx, 1)
	}

	writeHeight, err := v.rpc.Blob.Submit(ctx, []*blob.Blob{writeBlob}, nil)
	if err != nil {
		v.errCh <- err
		v.Metrics.Errors.Add(ctx, 1)
	}

	logrus.Info("wrote blob to block", writeHeight)

	// stick it in the history
	// to verify later

	v.History.Logs[fmt.Sprintf("%v", writeHeight)] = DataLog{
		BlockHeight: writeHeight,
		Namespace:   writeBlob.Namespace(),
		Data:        writeBlob.Data,
	}

	return 0
}
