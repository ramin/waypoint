package verifier

import (
	"context"
	"fmt"

	"github.com/celestiaorg/celestia-node/share"
)

// History is a data structure that maintains logs and provides summary statistics.
type History struct {
	Logs LogStore

	Misc MiscStats
}

// NewBlock creates a new block log and appends it to the History logs.
func (bt *History) NewBlock(bh uint64, ns share.Namespace, data []byte, success bool) {
	bl := DataLog{
		BlockHeight:  bh,
		Namespace:    ns,
		Data:         data,
		WriteSuccess: success,
	}
	bt.Logs[fmt.Sprintf("%v", bh)] = bl
}

// LogRetrievalAttempt logs the retrieval data for a given block height.
func (bt *History) LogRetrievalAttempt(data []byte, currentContext context.Context) {
	// todo,

	// using the current data (retrieval)
	// verify against stored data
	// does it match?
	// if not,
}

func (bt *History) RecordAtHeight(height int64) (*DataLog, error) {
	if val, ok := bt.Logs[fmt.Sprintf("%v", height)]; ok {
		return &val, nil
	}
	bt.Misc.Misses++

	return nil, fmt.Errorf(
		fmt.Sprintf("no data found at height: %v", height),
		ErrEntryNotFound,
	)
}

// Stats generates and returns summary statistics.
func (bt *History) Stats() (int, int, int, []string) {
	totalWrites := len(bt.Logs)
	successfulWrites := 0
	failedWrites := 0
	failureReasons := make([]string, 0)

	for _, log := range bt.Logs {
		if log.WriteSuccess {
			successfulWrites++
		} else {
			failedWrites++
			failureReasons = append(failureReasons, "Common failure reason")
		}
	}

	return totalWrites, successfulWrites, failedWrites, failureReasons
}

func NewHistory() (*History, error) {
	return &History{
		Logs: make(LogStore),
	}, nil
}
