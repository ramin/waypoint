package verifier

import (
	"time"

	"github.com/celestiaorg/celestia-node/share"
)

// DataLog represents the logging structure for each block.
type DataLog struct {
	BlockHeight    uint64          // Block height
	Namespace      share.Namespace // Namespace where data is written
	Data           []byte          // Data blob written
	WriteSuccess   bool            // Whether the write was successful or not
	WrittenAt      time.Time       // Time of write
	Retrieved      bool            // Whether the data was retrieved
	RetrievedAt    time.Time       // Time of retrieval
	RetrievedBlock int64           // Block at which the data was retrieved
	Duration       time.Duration   // Duration for retrieval
}

type LogStore map[string]DataLog
