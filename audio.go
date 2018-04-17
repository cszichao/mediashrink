package mediashrink

import (
	"fmt"
	"os/exec"
	"strconv"
)

// getAudioInfo get audio and video duration in secs with video dimension as well
func getAudioInfo(audioPath string) (*MediaInfo, error) {
	d, err := getDuration(audioPath)
	if err != nil {
		return nil, err
	}
	return &MediaInfo{0, 0, d, "", ""}, nil
}

// makeNullAudio make a null audio using aInfo, returns nil if success
func (aInfo *MediaInfo) makeNullAudio(outputPath string) error {
	// ffmpeg -f lavfi -i anullsrc=sample_rate=11025 -t 10.231  -metadata title="signature" silence.mp4
	// ffmpeg DTS delay time -11ms
	dtsDelay := float32(0.011)
	audioDuration := fmt.Sprintf("%.3f", float32(aInfo.Duration)/1000-dtsDelay)

	if info, err := exec.Command(
		commands.FFMPEG.FFMpeg,
		"-loglevel", "fatal",
		"-y", "-f", "lavfi", "-i", "anullsrc=sample_rate=128000",
		"-t", audioDuration,
		"-metadata", "title=\""+aInfo.Signature+"\"",
		outputPath,
	).CombinedOutput(); err != nil {
		return fmt.Errorf("exec ffmpeg %s with err: %s, info: %s", outputPath, err, info)
	}
	return nil
}

func getDuration(filePath string) (uint32, error) {
	// ffprobe -v quiet -show_entries format=duration -of default=noprint_wrappers=1:nokey=1

	// get raw duration output
	durationOutput, err1 := exec.Command(
		commands.FFMPEG.FFProbe,
		"-v", "quiet",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath).CombinedOutput()
	if err1 != nil {
		return 0, fmt.Errorf("exec ffprobe %s with err: %s", filePath, err1)
	}

	// convert raw duration into numeric
	duration, err2 := getDurationFromBytes(durationOutput)
	if err2 != nil {
		return 0, fmt.Errorf("failed get media duration from %s with err %s", durationOutput, err2)
	}

	return duration, nil
}

func getDurationFromBytes(durationOutput []byte) (uint32, error) {
	// parse duration
	durationIndex := 0
	for index, b := range durationOutput {
		if durationIndex == 0 && b == '\n' {
			durationIndex = index
			break
		}
	}
	if durationIndex <= 0 || durationIndex > len(durationOutput) {
		return 0, fmt.Errorf("error occurred when convert %s to duration", durationOutput)
	}
	durationStr := string(durationOutput[0:durationIndex])
	i, err := strconv.ParseFloat(durationStr, 32)
	if err == nil {
		return uint32(i * 1000), nil
	}
	return 0, fmt.Errorf("error occurred when convert %s to int: %s", durationStr, err)
}
