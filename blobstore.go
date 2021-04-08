// SPDX-License-Identifier: MIT

package ssb

import (
	"context"
	"fmt"
	"io"

	"go.cryptoscope.co/luigi"
	"go.cryptoscope.co/muxrpc/v2"

	refs "go.mindeco.de/ssb-refs"
)

const (
	// BlobStoreOpPut is used in put notifications
	BlobStoreOpPut BlobStoreOp = "put"

	// BlobStoreOpRm is used in remove notifications
	BlobStoreOpRm BlobStoreOp = "rm"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mock/blobstore.go . BlobStore

// BlobStore is the storage interface to read/write blobs as well as list and remove them and get updated about changes to it.
type BlobStore interface {
	// Get returns a reader of the blob with given ref.
	Get(ref *refs.BlobRef) (io.Reader, error)

	// Put stores the data in the reader in the blob store and returns the address.
	Put(blob io.Reader) (*refs.BlobRef, error)

	// PutExpected makes sure the added blob really is the passedBlobref
	// helpful for want/get operations which don't want to wast resources
	// PutExpected(io.Reader, *refs.BlobRef) error

	// Delete deletes a blob from the blob store.
	Delete(ref *refs.BlobRef) error

	// List returns a source of the refs of all stored blobs.
	List() luigi.Source

	// Size returns the size of the blob with given ref.
	Size(ref *refs.BlobRef) (int64, error)

	// Changes returns a broadcast that emits put and remove notifications.
	Changes() luigi.Broadcast
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mock/wantmanager.go . WantManager

// WantManager handles acquiting new blobs via `createWant` muxrpc streams.
type WantManager interface {
	io.Closer
	luigi.Broadcast // todo: replace with a typed broadcast

	Want(ref *refs.BlobRef) error
	Wants(ref *refs.BlobRef) bool
	WantWithDist(ref *refs.BlobRef, dist int64) error
	//Unwant(ref *refs.BlobRef) error
	CreateWants(context.Context, *muxrpc.ByteSink, muxrpc.Endpoint) luigi.Sink

	AllWants() []BlobWant
}

type BlobWant struct {
	Ref *refs.BlobRef

	// if Dist is negative, it is the hop count to the original wanter.
	// if it is positive, it is the size of the blob.
	Dist int64
}

func (w BlobWant) String() string {
	return fmt.Sprintf("%s:%d", w.Ref.ShortRef(), w.Dist)
}

// BlobStoreNotification contains info on a single change of the blob store.
// Op is either "rm" or "put".
type BlobStoreNotification struct {
	Op  BlobStoreOp
	Ref *refs.BlobRef
}

func (bn BlobStoreNotification) String() string {
	return bn.Op.String() + ": " + bn.Ref.Ref()
}

// BlobStoreOp specifies the operation in a blob store notification.
type BlobStoreOp string

// String returns the string representation of the operation.
func (op BlobStoreOp) String() string {
	return string(op)
}
