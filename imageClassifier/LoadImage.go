package imageClassifier

import (
	"image"
	"os"
)

func loadImage(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)

	return img, err
}