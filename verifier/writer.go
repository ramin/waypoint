package verifier

import (
	"context"
	"fmt"

	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/header"
	"github.com/ramin/waypoint/generator"
	"github.com/sirupsen/logrus"
)

func (v *Verifier) StartWriter(ctx context.Context) {
	go func() {
		for {
			select {
			default:
				for header := range v.AwaitBlock(ctx) {
					go v.WriteToBlock(ctx, header.Height())
					// fmt.Println(height)
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

func (v *Verifier) AwaitBlock(ctx context.Context) <-chan *header.ExtendedHeader {
	fmt.Println("awaiting block and writing new data")

	heights, err := v.rpc.Header.Subscribe(ctx)
	if err != nil {
		v.errCh <- err
		v.Metrics.Errors.Add(ctx, 1)
		panic(err)
	}

	return heights
}

func (v *Verifier) WriteToBlock(ctx context.Context, height uint64) uint64 {
	fmt.Println("awaiting block and writing new data")
	fmt.Println(height)

	writeBlob, err := generator.NewBlob()
	if err != nil {
		v.errCh <- err
		v.Metrics.Errors.Add(ctx, 1)
	}

	writeHeight, err := v.rpc.Blob.Submit(ctx, []*blob.Blob{writeBlob}, nil)
	if err != nil {
		v.errCh <- err
		fmt.Println(err)
		logrus.Error("failed to wrote blob to block ", writeHeight)
		v.Metrics.Errors.Add(ctx, 1)
		return writeHeight
	}

	logrus.Info("wrote blob to block ", writeHeight)

	// stick it in the history
	// to verify later

	v.History.Logs[fmt.Sprintf("%v", writeHeight)] = DataLog{
		BlockHeight: writeHeight,
		Namespace:   writeBlob.Namespace(),
		Data:        writeBlob.Data,
	}

	return writeHeight
}
