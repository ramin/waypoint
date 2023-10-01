package verifier

import (
	"context"
	"fmt"
	"time"

	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/header"
	"github.com/ramin/waypoint/config"
	"github.com/ramin/waypoint/generator"
	"github.com/sirupsen/logrus"
)

func (v *Verifier) StartWriter(ctx context.Context) {
	go func() {
		for {
			select {
			default:
				for header := range v.AwaitBlock(ctx) {
					// grab width of data here too
					// length of row routes (of DA header) / 2
					go v.WriteToBlock(ctx, header.Height())
				}
			case <-v.sig:
				logrus.Info("Received shutdown signal")
				close(v.done)
				return
			case <-ctx.Done():
				logrus.Info("Context canceled")
				close(v.done)
				return
			}
		}
	}()
}

func (v *Verifier) AwaitBlock(ctx context.Context) <-chan *header.ExtendedHeader {
	logrus.Info("awaiting block and writing new data")

	heights, err := v.rpc.Header.Subscribe(ctx)
	if err != nil {
		v.errCh <- err
		v.Metrics.Errors.Add(ctx, 1)
		panic(err)
	}

	return heights
}

func (v *Verifier) WriteToBlock(ctx context.Context, height uint64) uint64 {
	logrus.Info("awaiting block and writing new data")
	logrus.Info(height)

	writeBlob, err := generator.NewBlob()
	if err != nil {
		v.errCh <- err
		v.Metrics.Errors.Add(ctx, 1)
	}

	writeHeight, err := v.rpc.Blob.Submit(ctx, []*blob.Blob{writeBlob}, nil)
	if err != nil {
		v.errCh <- err
		logrus.Info(err)
		logrus.Error("failed to wrote blob to block ", writeHeight)
		v.Metrics.Errors.Add(ctx, 1)
		return writeHeight
	}

	v.Metrics.Writes.Add(context.Background(), 1)
	logrus.Info("wrote blob to block ", writeHeight)

	v.History.Logs[fmt.Sprintf("%v", writeHeight)] = DataLog{
		BlockHeight: writeHeight,
		Namespace:   writeBlob.Namespace(),
		Data:        writeBlob.Data,
		WrittenAt:   time.Now(),
		Duration:    config.Read().ReadInterval,
	}

	return writeHeight
}
