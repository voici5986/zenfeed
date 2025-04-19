package index

import (
	"context"
	"io"
)

// Codec defines interface for encoding and decoding index.
type Codec interface {
	// EncodeTo encodes the index to the given writer.
	EncodeTo(ctx context.Context, w io.Writer) (err error)
	// DecodeFrom decodes the index from the given reader.
	DecodeFrom(ctx context.Context, r io.Reader) (err error)
}
