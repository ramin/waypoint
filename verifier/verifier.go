package verifier

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
)

type Verifier struct {
	History *History
	Metrics *Metrics

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

	m, err := NewMetrics()
	if err != nil {
		return nil, err
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

func (v *Verifier) WithClient(c *client.Client) *Verifier {
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
			fmt.Println("error", err)
		}
	}()
	return v.done
}
