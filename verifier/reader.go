package verifier

import (
	"context"
	"time"

	"github.com/celestiaorg/celestia-node/share"
	"github.com/ramin/waypoint/config"
	"github.com/sirupsen/logrus"
)

func (v *Verifier) StartReader(ctx context.Context) {
	ticker := time.NewTicker(config.Read().ReadInterval)

	defer ticker.Stop()

	// start a periodic stats producer
	if config.Read().DisplayInfo {
		go v.PeriodicStats(ctx)
	}

	for {
		select {
		case <-ticker.C:
			go v.verifyRecords()
		case <-ctx.Done():
			return
		case <-v.done: // listen to the done channel for termination
			return
		}
	}
}

func (v *Verifier) PeriodicStats(ctx context.Context) {
	ticker := time.NewTicker(config.Read().ReadInterval)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logrus.Debug(v.History.Logs)
		case <-ctx.Done():
			return
		case <-v.done:
			return
		}
	}
}

func (v *Verifier) verifyRecords() {
	logrus.Debug("verifying records")

	for key, log := range v.History.Logs {
		if log.Retrieved || !log.WriteSuccess {
			delete(v.History.Logs, key)
		}

		// only proceed to check the write after the set
		// await duration has elapsed
		if time.Since(log.WrittenAt) < log.Duration {
			continue
		}

		logrus.Debug(log.BlockHeight)
		logrus.Debug(log.Namespace)
		logrus.Debug(log.Namespace.ID())

		// verify what errors come from here
		blob, err := v.rpc.Blob.GetAll(
			context.Background(),
			log.BlockHeight,
			[]share.Namespace{log.Namespace},
		)

		if err != nil {
			logrus.Error(err)
			continue
		}

		logrus.Debug(blob)

		// assume we'll need to switch on error type here
		if err != nil {
			v.Metrics.Errors.Add(context.Background(), 1)
			v.Metrics.Misses.Add(context.Background(), 1)
			v.errCh <- err
		} else {
			v.Metrics.Reads.Add(context.Background(), 1)
		}

		delete(v.History.Logs, key)
	}
}
