package id3v2

import (
	"bytes"
	"errors"
	"github.com/bogem/id3v2/util"
	"io"
)

const (
	ID3Identifier = "ID3"
	TagHeaderSize = 10
)

type TagHeader struct {
	FramesSize uint32
	Version    byte
}

func ParseHeader(rd io.Reader) (*TagHeader, error) {
	data := make([]byte, TagHeaderSize)
	n, err := rd.Read(data)
	if n < TagHeaderSize {
		err = errors.New("Size of tag header is less than expected")
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	if !isID3Tag(data[0:3]) {
		return nil, nil
	}

	size, err := util.ParseSize(data[6:])
	if err != nil {
		return nil, err
	}

	header := &TagHeader{
		Version:    data[3],
		FramesSize: size,
	}

	return header, nil
}

func isID3Tag(data []byte) bool {
	if len(data) != len(ID3Identifier) {
		return false
	}
	return string(data[0:3]) == ID3Identifier
}

func FormTagHeader(framesSize []byte) []byte {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()

	// Identifier
	b.WriteString(ID3Identifier)

	// Version
	b.WriteByte(4)

	// Revision
	b.WriteByte(0)

	// Flags
	b.WriteByte(0)

	// Size of frames
	b.Write(framesSize)

	bytesBufPool.Put(b)
	return b.Bytes()
}
