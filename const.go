package mediashrink

import (
	"errors"

	"github.com/haxii/filetype/matchers"
	"github.com/haxii/filetype/types"
)

// support media types
var (
	image = map[string]types.Type{
		matchers.TypeJpeg.Extension: matchers.TypeJpeg,
		matchers.TypeJpe.Extension:  matchers.TypeJpe,
		matchers.TypeJpg.Extension:  matchers.TypeJpg,
		matchers.TypePng.Extension:  matchers.TypePng,
		matchers.TypeGif.Extension:  matchers.TypeGif,
		matchers.TypeTif.Extension:  matchers.TypeTif,
		matchers.TypeTiff.Extension: matchers.TypeTiff,
		matchers.TypeBmp.Extension:  matchers.TypeBmp,
		matchers.TypeIco.Extension:  matchers.TypeIco,
		matchers.TypeJfif.Extension: matchers.TypeJfif,
	}
	audio = map[string]types.Type{
		matchers.TypeMp3.Extension:  matchers.TypeMp3,
		matchers.TypeM4a.Extension:  matchers.TypeM4a,
		matchers.TypeOgg.Extension:  matchers.TypeOgg,
		matchers.TypeFlac.Extension: matchers.TypeFlac,
		matchers.TypeWav.Extension:  matchers.TypeWav,
		matchers.TypeAac.Extension:  matchers.TypeAac,
		matchers.TypeWma.Extension:  matchers.TypeWma,
		matchers.TypeCaf.Extension:  matchers.TypeCaf,
	}
	video = map[string]types.Type{
		matchers.TypeMp4.Extension:  matchers.TypeMp4,
		matchers.TypeM4v.Extension:  matchers.TypeM4v,
		matchers.TypeMkv.Extension:  matchers.TypeMkv,
		matchers.TypeMov.Extension:  matchers.TypeMov,
		matchers.TypeAvi.Extension:  matchers.TypeAvi,
		matchers.TypeWmv.Extension:  matchers.TypeWmv,
		matchers.TypeMpeg.Extension: matchers.TypeMpeg,
		matchers.TypeMpg.Extension:  matchers.TypeMpg,
		matchers.TypeFlv.Extension:  matchers.TypeFlv,
		matchers.TypeAsf.Extension:  matchers.TypeAsf,
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
