package verifier

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/ramin/waypoint/config"
	"github.com/ramin/waypoint/metrics"
	"github.com/sirupsen/logrus"
)

type Verifier struct {
	History *History
	Metrics *metrics.Metrics

	checkpoint Checkpoint

	rpc *client.Client

	done  chan bool
	errCh chan error
	sig   chan os.Signal
}

type Checkpoint struct {
	Time   time.Time
	Height int64
}

func NewVerifier(ctx context.Context) (*Verifier, error) {
	h, err := NewHistory()
	if err != nil {
		return nil, err
	}

	var m *metrics.Metrics

	if config.Read().MetricsEnabled {
		m, err = metrics.NewMetrics()
		if err != nil {
			return nil, err
		}
	}

	return &Verifier{
		History:    h,
		Metrics:    m,
		checkpoint: Checkpoint{},
		sig:        make(chan os.Signal, 1),
		done:       make(chan bool),
		errCh:      make(chan error),
	}, nil
}

func NewVerifierWithClient(ctx context.Context) (*Verifier, error) {
	v, err := NewVerifier(ctx)
	if err != nil {
		return nil, err
	}
	v = v.WithClient(ctx)
	return v, nil
}

func (v *Verifier) WithClient(ctx context.Context) *Verifier {
	cfg := config.Read()

	// by default, celestia-nodes run RPC on port 26658
	rpc, err := client.NewClient(
		ctx,
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		cfg.JWT,
	)
	if err != nil {
		panic(err)
	}

	return v.AddClient(rpc)
}

func (v *Verifier) AddClient(c *client.Client) *Verifier {
	v.rpc = c
	return v
}

func (v *Verifier) Start(ctx context.Context) chan bool {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go v.StartWriter(ctx)

	go v.StartReader(ctx)

	// print out errors in background
	go func() {
		for err := range v.errCh {
			logrus.Error("error", err)
		}
	}()
	return v.done
}
