package core

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// A Stream type wraps a Git object and facilitates formatting the object into
// the full on-disk format: a header followed by the actual stream of bytes
// that is the object's content.
type Stream struct {
	Object
	checksum Sha1
}

var _ EncodeDecoder = &Stream{}

// Reader returns an io.Reader that prepends the underlying encoded object
// stream with a valid header, which consists of the human-readable identifier
// of the blob type, followed by a space, followed by the ASCII representation
// of the number of bytes in the object content, and terminated with a NULL
// byte.
func (stream *Stream) Reader() io.Reader {
	object := stream.Object
	return io.MultiReader(
		strings.NewReader(object.Type()),
		bytes.NewReader([]byte{' '}),
		strings.NewReader(strconv.FormatInt(int64(object.Size()), 10)),
		bytes.NewReader([]byte{0}),
		object.Reader(),
	)
}

// Bytes returns a whole byte array, obtained by allocating a byte array buffer
// and draining this stream's reader into it while growing the buffer as needed.
func (stream *Stream) Bytes() []byte {
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(stream.Reader()); err != nil {
		Die(err)
	}
	return buffer.Bytes()
}

// Hash returns the SHA-1 checksum of this stream's byte representation. This
// checksum is only calculated once and then cached.
func (stream *Stream) Hash() Sha1 {
	if stream.checksum != "" {
		return stream.checksum
	}
	return stream.Rehash()
}

// Rehash clears the cached checksum, recalculates the checksum of this stream's
// byte representation, and caches this new checksum. You must call this method
// if you change the underlying object wrapped by a Stream.
func (stream *Stream) Rehash() (checksum Sha1) {
	hash := sha1.New()
	if _, err := io.Copy(hash, stream.Reader()); err != nil {
		Die(err)
	}
	checksum = Sha1(fmt.Sprintf("%x", hash.Sum(nil)))
	stream.checksum = checksum
	return
}

// Decode parses an object represented in its entirety by a byte sequence
// yielded from an io.Reader and reconstructs the object with the correct
// type as inferred from the sequence's header.
func (stream *Stream) Decode(reader io.Reader) error {
	r := bufio.NewReader(reader)

	if typeString, err := r.ReadString(byte(' ')); err != nil {
		return err
	} else {
		typeString = typeString[:len(typeString)-1]

		switch typeString {
		case "blob":
			stream.Object = &Blob{}
		case "tree":
			stream.Object = &Tree{}
		default:
			return errors.New(fmt.Sprintf("%v is not a known object type", typeString))
		}
	}

	if lenString, err := r.ReadString(byte(0)); err != nil {
		return err
	} else {
		lenString = lenString[:len(lenString)-1]
		length, err := strconv.ParseInt(lenString, 10, 64)
		if err != nil {
			return err
		}

		rest := &io.LimitedReader{R: r, N: length}
		if err := stream.Object.Decode(rest); err != nil {
			return err
		}
	}

	return nil
}
