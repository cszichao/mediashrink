package mediashrink

import (
	"errors"

	"github.com/haxii/filetype/matchers"
)

// support media types
var (
	image = map[string]matchers.Matcher{
		matchers.TypeJpeg.Extension: matchers.Jpeg,
		matchers.TypeJpe.Extension:  matchers.Jpe,
		matchers.TypeJpg.Extension:  matchers.Jpg,
		matchers.TypePng.Extension:  matchers.Png,
		matchers.TypeGif.Extension:  matchers.Gif,
		matchers.TypeTif.Extension:  matchers.Tif,
		matchers.TypeTiff.Extension: matchers.Tiff,
		matchers.TypeBmp.Extension:  matchers.Bmp,
		matchers.TypeIco.Extension:  matchers.Ico,
	}
	imagePNG = map[string]matchers.Matcher{matchers.TypePng.Extension: matchers.Png}

	audio = map[string]matchers.Matcher{
		matchers.TypeMp3.Extension:  matchers.Mp3,
		matchers.TypeM4a.Extension:  matchers.M4a,
		matchers.TypeOgg.Extension:  matchers.Ogg,
		matchers.TypeFlac.Extension: matchers.Flac,
		matchers.TypeWav.Extension:  matchers.Wav,
		matchers.TypeAac.Extension:  matchers.Aac,
		matchers.TypeWma.Extension:  matchers.Wma,
		matchers.TypeCaf.Extension:  matchers.Caf,
	}

	video = map[string]matchers.Matcher{
		matchers.TypeMp4.Extension:  matchers.Mp4,
		matchers.TypeM4v.Extension:  matchers.M4v,
		matchers.TypeMkv.Extension:  matchers.Mkv,
		matchers.TypeMov.Extension:  matchers.Mov,
		matchers.TypeAvi.Extension:  matchers.Avi,
		matchers.TypeWmv.Extension:  matchers.Wmv,
		matchers.TypeMpeg.Extension: matchers.Mpeg,
		matchers.TypeMpg.Extension:  matchers.Mpeg,
		matchers.TypeFlv.Extension:  matchers.Flv,
		matchers.TypeAsf.Extension:  matchers.Asf,
	}
)

var (
	// ErrUnknownMediaType error message when the given file is not recognized as a shrinkable media type
	ErrUnknownMediaType = errors.New("unknown media type")
)

var (
	// commands used for exec by golang
	commands = &CommandNames{
		ImageMagicK: &ImageMagicKExec{
			Identify: "identify",
			Convert:  "convert",
		},
		FFMPEG: &FFMPEGExec{
			FFMpeg:  "ffmpeg",
			FFProbe: "ffprobe",
		},
		P7Zip: &P7ZipExec{
			P7z: "7z",
		},
	}
)

// CommandNames for exec
type CommandNames struct {
	FFMPEG      *FFMPEGExec
	ImageMagicK *ImageMagicKExec
	P7Zip       *P7ZipExec
}

// FFMPEGExec ...
type FFMPEGExec struct {
	FFProbe string
	FFMpeg  string
}

// ImageMagicKExec ...
type ImageMagicKExec struct {
	Identify string
	Convert  string
}

// P7ZipExec ...
type P7ZipExec struct {
	P7z string
}
