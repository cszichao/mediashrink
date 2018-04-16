package mediashrink

import (
	"fmt"
	"os/exec"
)

// getVideoInfo get audio and video duration in secs with video dimension as well
func getVideoInfo(videoPath string) (*MediaInfo, error) {
	info := &MediaInfo{0, 0, 0, "", ""}
	// ffprobe -v quiet -print_format json -show_streams -show_format
	// ffprobe -v quiet -show_entries stream=width,height -of default=noprint_wrappers=1:nokey=1

	// get dimension
	if dimensionOutput, err := exec.Command(
		commands.FFMPEG.FFProbe,
		"-v", "quiet",
		"-show_entries", "stream=width,height",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("exec ffprobe %s with err: %s", videoPath, err)
	} else if info.Width, info.Height, err = getWidthAndHeightFromBytes(dimensionOutput); err != nil {
		return nil, fmt.Errorf("failed get width & height from %s with err %s", dimensionOutput, err)
	}

	// get duration
	if d, err := getDuration(videoPath); err == nil {
		info.Duration = d
	} else {
		return nil, err
	}
	return info, nil
}

// makeNullVideo make a null video using vInfo, returns nil if success
func (vInfo *MediaInfo) makeNullVideo(outputPath string) error {
	// ffmpeg -f lavfi -i color=#123456:s=640x480:d=10.231 \
	//        -f lavfi -i anullsrc=sample_rate=11025 -t 10.231  silence.mp4
	// ffmpeg DTS delay time -11ms
	dtsDelay := float32(0.011)
	videoDimension := fmt.Sprintf("%dx%d", vInfo.Width, vInfo.Height)
	videoDuration := fmt.Sprintf("%.2f", float32(int(vInfo.Duration/10))/100-dtsDelay)
	audioDuration := fmt.Sprintf("%.3f", float32(vInfo.Duration)/1000-dtsDelay)

	if info, err := exec.Command(
		commands.FFMPEG.FFMpeg,
		"-loglevel", "panic",
		"-f", "lavfi", "-i", "color=#"+vInfo.Signature+":s="+videoDimension+":d="+videoDuration,
		"-f", "lavfi", "-i", "anullsrc=sample_rate=128000", "-t", audioDuration,
		outputPath,
	).CombinedOutput(); err != nil {
		return fmt.Errorf("exec ffmpeg %s with err: %s, info: %s", outputPath, err, info)
	}
	return nil
}
