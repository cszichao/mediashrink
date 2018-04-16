package mediashrink

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ErrHashSum weird error, should it actually occurred?
var ErrHashSum = errors.New("error occurred when making hash sum")

// fileMD5 cal MD5 of a giving file
func fileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)
	if len(hashInBytes) >= 16 {
		return hex.EncodeToString(hashInBytes[:16]), nil
	}
	return "", ErrHashSum
}

// getWidthAndHeightFromBytes get w & h from "1024\n768\n..." bytes
func getWidthAndHeightFromBytes(info []byte) (uint32, uint32, error) {
	// cut width & height from original bytes
	width := uint32(0)
	widthIndex := 0
	height := uint32(0)
	heightIndex := 0
	for index, b := range info {
		if widthIndex == 0 && b == '\n' {
			widthIndex = index
			continue
		}
		if heightIndex == 0 && b == '\n' {
			heightIndex = index
			break
		}
	}
	if widthIndex <= 0 || heightIndex <= widthIndex+1 || heightIndex > len(info) {
		return 0, 0, fmt.Errorf("error occurred when convert %s to int", info)
	}
	// width
	widthStr := string(info[0:widthIndex])
	if i, err := strconv.Atoi(widthStr); err == nil {
		width = uint32(i)
	} else {
		return 0, 0, fmt.Errorf("error occurred when convert %s to int:%s", widthStr, err)
	}
	// height
	heightStr := string(info[widthIndex+1 : heightIndex])
	if i, err := strconv.Atoi(heightStr); err == nil {
		height = uint32(i)
	} else {
		width = 0
		return 0, 0, fmt.Errorf("error occurred when convert %s to int:%s", heightStr, err)
	}
	return width, height, nil
}
