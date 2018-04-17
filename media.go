package mediashrink

import (
	"fmt"
	"path/filepath"
	"strconv"
)

func isImage(ext string) bool {
	_, exists := image[ext]
	return exists
}

func isVideo(ext string) bool {
	_, exists := video[ext]
	return exists
}

func isAudio(ext string) bool {
	_, exists := audio[ext]
	return exists
}

// MediaInfo shows the media's dimension & duration
type MediaInfo struct {
	Width     uint32
	Height    uint32
	Duration  uint32
	Signature string
	Ext       string
}

// ToString convert MediaInfo To String width[x]height[x]duration[x]signature[.]ext
func (info *MediaInfo) ToString() string {
	return fmt.Sprintf("%dx%dx%dx%s.%s", info.Width,
		info.Height, info.Duration, info.Signature, info.Ext)
}

// GetMediaInfo return the MediaInfo if path is a valid media, otherwise return null.
// sig: hex string in min length of 6, should be a MD5 string normally,
func GetMediaInfo(sig string, path string) (*MediaInfo, error) {
	ext := filepath.Ext(path)
	if len(ext) <= 1 {
		return nil, ErrUnknownMediaType
	}

	if len(sig) == 0 {
		if md5, err := fileMD5(path); err == nil {
			sig = md5
		} else {
			return nil, err
		}
	}
	signature := validateSignature(sig)
	if len(signature) == 0 {
		return nil, fmt.Errorf("wrong signature %s for file %s", sig, path)
	}

	var mediaInfo *MediaInfo
	var err error
	if isImage(ext) {
		if mediaInfo, err = getImageInfo(path); err != nil {
			return nil, err
		} else if mediaInfo.Width <= 0 || mediaInfo.Height <= 0 {
			return nil, ErrUnknownMediaType

		}
	} else if isVideo(ext) {
		if mediaInfo, err = getVideoInfo(path); err != nil {
			return nil, err
		} else if mediaInfo.Width <= 0 || mediaInfo.Height <= 0 || mediaInfo.Duration <= 0 {
			return nil, ErrUnknownMediaType
		}
	} else if isAudio(ext) {
		if mediaInfo, err = getAudioInfo(path); err != nil {
			return nil, err
		} else if mediaInfo.Duration <= 0 {
			return nil, ErrUnknownMediaType
		}
	}
	mediaInfo.Ext = ext
	mediaInfo.Signature = signature
	return mediaInfo, nil

}

// Shrink makes a shrink media using info
func (info *MediaInfo) Shrink(outputPath string) error {
	if isImage(info.Ext) {
		return info.makeNullImage(outputPath)
	} else if isVideo(info.Ext) {
		return info.makeNullVideo(outputPath)
	} else if isAudio(info.Ext) {
		return info.makeNullAudio(outputPath)
	}

	return fmt.Errorf("making a unsupported format %s", info.ToString())
}

func validateSignature(s string) string {
	if len(s) < 6 {
		return ""
	}
	for _, b := range s[:6] {
		if !(('a' <= b && b <= 'f') || ('0' <= b && b <= '9')) {
			return ""
		}
	}
	return s[:6]
}

// MediaInfoFromString convert String width[x]height[x]duration[[x]signature[.]ext To MediaInfo
func MediaInfoFromString(str string) (*MediaInfo, error) {
	heightIndex := 0
	durationIndex := 0
	signatureIndex := 0
	extIndex := 0
	for index, b := range str {
		if heightIndex == 0 && b == 'x' {
			heightIndex = index
			continue
		}
		if durationIndex == 0 && b == 'x' {
			durationIndex = index
			continue
		}
		if signatureIndex == 0 && b == 'x' {
			signatureIndex = index
			continue
		}
		if extIndex == 0 && b == '.' {
			extIndex = index
			break
		}
	}
	if heightIndex <= 0 ||
		durationIndex <= heightIndex+1 ||
		signatureIndex <= durationIndex+1 ||
		extIndex <= signatureIndex+1 ||
		len(str) <= extIndex+1 {
		return nil, fmt.Errorf("error occurred when convert %s to media info", str)
	}

	info := &MediaInfo{0, 0, 0, "", ""}
	// width
	widthStr := str[0:heightIndex]
	if i, err := strconv.Atoi(widthStr); err == nil {
		info.Width = uint32(i)
	} else {
		return nil, fmt.Errorf("error occurred when convert %s to int:%s", widthStr, err)
	}
	// height
	heightStr := str[heightIndex+1 : durationIndex]
	if i, err := strconv.Atoi(heightStr); err == nil {
		info.Height = uint32(i)
	} else {
		info.Width = 0
		return nil, fmt.Errorf("error occurred when convert %s to int:%s", heightStr, err)
	}
	// duration
	durationStr := str[durationIndex+1 : signatureIndex]
	if i, err := strconv.Atoi(durationStr); err == nil {
		info.Duration = uint32(i)
	} else {
		info.Width = 0
		info.Height = 0
		return nil, fmt.Errorf("error occurred when convert %s to int:%s", durationStr, err)
	}
	// signature
	sigStr := validateSignature(str[signatureIndex+1 : extIndex])
	if len(sigStr) > 0 {
		info.Signature = sigStr
	} else {
		info.Width = 0
		info.Height = 0
		info.Duration = 0
		return nil, fmt.Errorf("error occurred when convert %s to signature", sigStr)
	}
	// ext
	info.Ext = str[extIndex+1:]
	return info, nil
}
