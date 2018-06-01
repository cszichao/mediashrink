package mediashrink

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/haxii/filetype/matchers"
)

func isImage(ext string) bool {
	_, exists := image[ext]
	return exists
}

func isPNG(ext string) bool {
	_, exists := imagePNG[ext]
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
	Duration  uint32 // in ms
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
// set guessMissingExt to true to guess the media type when no ext presented in path.
func GetMediaInfo(sig string, guessMissingExt bool, path string) (*MediaInfo, error) {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		ext = strings.ToLower(ext[1:])
	} else if guessMissingExt {
		ext = guessExt(path)
	}
	if len(ext) == 0 {
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
	if isPNG(ext) { //test image in a more fast and compatible way
		if mediaInfo, err = getImagePNGInfo(path); err != nil {
			return nil, err
		} else if mediaInfo.Width <= 0 || mediaInfo.Height <= 0 {
			return nil, ErrUnknownMediaType
		}
	} else if isImage(ext) {
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
	} else {
		return nil, ErrUnknownMediaType
	}
	mediaInfo.Ext = ext
	mediaInfo.Signature = signature
	return mediaInfo, nil
}

// Shrink makes a shrink media using info
func (info *MediaInfo) Shrink(outputPath string) error {
	var err error
	safeOutputPath := outputPath + "." + info.Ext
	if isImage(info.Ext) {
		err = info.makeNullImage(safeOutputPath)
	} else if isVideo(info.Ext) {
		err = info.makeNullVideo(safeOutputPath)
	} else if isAudio(info.Ext) {
		err = info.makeNullAudio(safeOutputPath)
	}
	if err != nil {
		return err
	}
	if _, err := os.Stat(safeOutputPath); err == nil {
		return os.Rename(safeOutputPath, outputPath)
	}
	return fmt.Errorf("unsupported media format %s", info.ToString())
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

// guessExt guess if file is a supported media, if so return the ext, otherwise nil string
func guessExt(filePath string) string {
	guessedExt := ""
	if err := readFileHeader(filePath, func(header []byte, err error) error {
		if err != nil {
			return err
		}
		for ext, matcher := range image {
			if matcher(header) {
				guessedExt = ext
				return nil
			}
		}
		for ext, matcher := range video {
			if matcher(header) {
				guessedExt = ext
				return nil
			}
		}
		for ext, matcher := range audio {
			if matcher(header) {
				guessedExt = ext
				return nil
			}
		}
		return nil
	}); err != nil {
		return ""
	}
	// uniform the jpg extensions
	if guessedExt == "jpe" || guessedExt == "jpeg" {
		guessedExt = "jpg"
	}
	return guessedExt
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

// CheckCompatibility check compatibility of ImageMagicK and FFMpeg
func CheckCompatibility(output io.Writer, exportDir string) {
	mediaList := map[string]*MediaInfo{}
	duration := 5000
	for imgExt := range image {
		infoStr := "32x32x0x123456." + imgExt
		if i, err := MediaInfoFromString(infoStr); err == nil {
			mediaList[imgExt] = i
		} else {
			fmt.Fprintf(output, "failed to parse media info %s with error %s", infoStr, err)
			return
		}
	}
	for audioExt := range audio {
		infoStr := "0x0x" + strconv.Itoa(duration) + "x123456." + audioExt
		if a, err := MediaInfoFromString(infoStr); err == nil {
			mediaList[audioExt] = a
		} else {
			fmt.Fprintf(output, "failed to parse media info %s with error %s", infoStr, err)
			return
		}
	}
	for videoExt := range video {
		infoStr := "128x128x" + strconv.Itoa(duration) + "x123456." + videoExt
		if v, err := MediaInfoFromString(infoStr); err == nil {
			mediaList[videoExt] = v
		} else {
			fmt.Fprintf(output, "failed to parse media info %s with error %s", infoStr, err)
			return
		}
	}

	for ext, media := range mediaList {
		sample := filepath.Join(exportDir, "shrink."+ext)
		fmt.Fprintf(output, "generating %s to %s : ", ext, sample)
		if err := media.Shrink(sample); err != nil {
			fmt.Fprintf(output, "Failed with error %s\n", err)
		} else {
			fmt.Fprintf(output, "OKay\n")
			fmt.Fprintln(output, "guessing ext of ", sample, ":", guessExt(sample))
			fmt.Fprintf(output, "reading %s from %s :", ext, sample)
			if info, err := GetMediaInfo("", false, sample); err != nil {
				fmt.Fprintf(output, "Failed with error %s\n", err)
			} else {
				durationMargin := 0
				if info.Duration > 0 {
					durationMargin = int(info.Duration) - duration
				}
				fmt.Fprintln(output, "Okay:", info.ToString(), "duration margin: ", durationMargin, "ms")
			}
		}
		os.Remove(sample)
	}
}

// ImageMatchers image matchers
func ImageMatchers() map[string]matchers.Matcher {
	return image
}

// AudioMatchers audio matchers
func AudioMatchers() map[string]matchers.Matcher {
	return audio
}

// VideoMatchers video matchers
func VideoMatchers() map[string]matchers.Matcher {
	return video
}
