package mediashrink

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
)

var pngHeaderBufioPool sync.Pool

const maxPNGHeaderSize = 64

// AcquireReader acquire a buffered reader based on net connection
func acquirePNGHeaderReader(c io.Reader) *bufio.Reader {
	v := pngHeaderBufioPool.Get()
	if v == nil {
		return bufio.NewReaderSize(c, maxPNGHeaderSize)
	}
	r := v.(*bufio.Reader)
	r.Reset(c)
	return r
}

// releasePNGHeaderReader release a buffered reader
func releasePNGHeaderReader(r *bufio.Reader) {
	pngHeaderBufioPool.Put(r)
}

var (
	errNotPNG = errors.New("not a PNG file")
	pngIHDR   = []byte{0x49, 0x48, 0x44, 0x52}
)

// GetPNGInfo compatible with apple's CgBI file format
func GetPNGInfo(imagePath string) (*MediaInfo, error) {
	pngFile, err1 := os.Open(imagePath)
	if err1 != nil {
		return nil, err1
	}
	defer pngFile.Close()

	stat, err2 := pngFile.Stat()
	if err2 != nil {
		return nil, err2
	}

	reader := acquirePNGHeaderReader(pngFile)
	defer releasePNGHeaderReader(reader)

	sizeToPeek := int64(maxPNGHeaderSize)
	if sizeToPeek > stat.Size() {
		sizeToPeek = stat.Size()
	}
	header, _ := reader.Peek(int(sizeToPeek))
	isPNG := func() bool {
		return len(header) > 3 &&
			header[0] == 0x89 && header[1] == 0x50 &&
			header[2] == 0x4E && header[3] == 0x47
	}
	if !isPNG() {
		return nil, errNotPNG
	}

	// get index of "IHDR", image width and height are followed with it
	ihdrIndex := bytes.Index(header, pngIHDR)
	if ihdrIndex <= 0 || ihdrIndex+len(pngIHDR)+8 > len(header) {
		return nil, errNotPNG
	}
	widthIndex := ihdrIndex + len(pngIHDR)
	widthBytes := header[widthIndex : widthIndex+4]
	heightBytes := header[widthIndex+4 : widthIndex+8]

	// big endian
	width := int(widthBytes[3]) + int(widthBytes[2])<<8 + int(widthBytes[1])<<16 + int(widthBytes[0])<<24
	height := int(heightBytes[3]) + int(heightBytes[2])<<8 + int(heightBytes[1])<<16 + int(heightBytes[0])<<24

	return &MediaInfo{uint32(width), uint32(height), 0, "", ""}, nil
}
