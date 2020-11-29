package encoding

import (
	"google.golang.org/grpc/encoding"
)

// Codec defines the interface HTTP/gRPC uses to encode and decode messages.  Note
// that implementations of this interface must be thread safe; a Codec's
// methods can be called from concurrent goroutines.
type Codec encoding.Codec

// Compressor is used for compressing and decompressing when sending or
// receiving messages.
type Compressor encoding.Compressor

// GetCodec gets a registered Codec by content-subtype, or nil if no Codec is
// registered for the content-subtype.
//
// The content-subtype is expected to be lowercase.
func GetCodec(contentSubtype string) Codec {
	return encoding.GetCodec(contentSubtype)
}

// RegisterCodec registers the provided Codec for use with all HTTP/gRPC clients and
// servers.
func RegisterCodec(codec Codec) {
	encoding.RegisterCodec(codec)
}

// RegisterCompressor registers the compressor with HTTP/gRPC by its name.
func RegisterCompressor(c Compressor) {
	encoding.RegisterCompressor(c)
}

// GetCompressor returns Compressor for the given compressor name.
func GetCompressor(name string) Compressor {
	return encoding.GetCompressor(name)
}
