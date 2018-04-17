package mediashrink

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// getImageInfo get image size with width x height
func getImageInfo(imagePath string) (*MediaInfo, error) {
	info := &MediaInfo{0, 0, 0, "", ""}
	if output, err := exec.Command(
		commands.ImageMagicK.Identify,
		"-format", "%[fx:w]\n%[fx:h]\n", imagePath,
	).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("exec identify %s with err: %s, info: %s", imagePath, err, output)
	} else if info.Width, info.Height, err = getWidthAndHeightFromBytes(output); err != nil {
		return nil, fmt.Errorf("failed get width & height from %s with err %s", output, err)
	}
	return info, nil
}

// makeNullImage make a null image using imgInfo
func (imgInfo *MediaInfo) makeNullImage(outputPath string) error {
	// convert -size 1024x768 xc:white canvas.jpg
	imageSize := fmt.Sprintf("%dx%d", imgInfo.Width, imgInfo.Height)
	if info, err := exec.Command(
		commands.ImageMagicK.Convert,
		"-size", imageSize,
		"xc:#"+imgInfo.Signature, outputPath,
	).CombinedOutput(); err != nil {
		return fmt.Errorf("exec convert %s with err: %s, info: %s", outputPath, err, info)
	}
	return nil
}

var (
	errNotPNG = errors.New("not a PNG file")
	pngIHDR   = []byte{0x49, 0x48, 0x44, 0x52}
)

// getImagePNGInfo optimized info getter for png, compatible with apple's CgBI file format
func getImagePNGInfo(imagePath string) (*MediaInfo, error) {
	isPNG := func(header []byte) bool {
		return len(header) > 3 &&
			header[0] == 0x89 && header[1] == 0x50 &&
			header[2] == 0x4E && header[3] == 0x47
	}
	var width, height int
	if err := readFileHeader(imagePath, func(header []byte, err error) error {
		if err != nil {
			return err
		}
		if !isPNG(header) {
			return errNotPNG
		}

		// get index of "IHDR", image width and height are followed with it
		ihdrIndex := bytes.Index(header, pngIHDR)
		if ihdrIndex <= 0 || ihdrIndex+len(pngIHDR)+8 > len(header) {
			return errNotPNG
		}
		widthIndex := ihdrIndex + len(pngIHDR)
		widthBytes := header[widthIndex : widthIndex+4]
		heightBytes := header[widthIndex+4 : widthIndex+8]

		// big endian
		width = int(widthBytes[3]) + int(widthBytes[2])<<8 + int(widthBytes[1])<<16 + int(widthBytes[0])<<24
		height = int(heightBytes[3]) + int(heightBytes[2])<<8 + int(heightBytes[1])<<16 + int(heightBytes[0])<<24
		return nil
	}); err != nil {
		return nil, err
	}
	return &MediaInfo{uint32(width), uint32(height), 0, "", ""}, nil
}
