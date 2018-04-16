package mediashrink

import (
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
