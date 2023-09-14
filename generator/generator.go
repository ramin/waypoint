package generator

import (
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/blob/blobtest"
)

func NewBlob() (*blob.Blob, error) {
	appBlobs, err := blobtest.GenerateV0Blobs([]int{8}, false)
	if err != nil {
		panic(err)
	}

	return blob.NewBlob(
		appBlobs[0].ShareVersion,
		append([]byte{appBlobs[0].NamespaceVersion}, appBlobs[0].NamespaceID...),
		appBlobs[0].Data,
	)
}
