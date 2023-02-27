package imageClassifier

import (
	"fmt"
	"image"
	"runtime/debug"

	"github.com/disintegration/imaging"
	tf "github.com/galeone/tensorflow/tensorflow/go"
)

func createTensor(img image.Image) (*tf.Tensor, error) {
	nrgbaImg := imaging.Fill(img, 512, 512, imaging.Center, imaging.Lanczos)

	return imageToTensor(nrgbaImg, 512, 512)
}

func imageToTensor(img image.Image, imageHeight, imageWidth int) (tfTensor *tf.Tensor, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("classify: %s (panic)\nstack: %s", r, debug.Stack())
		}
	}()

	if imageHeight <= 0 || imageWidth <= 0 {
		return tfTensor, fmt.Errorf("classify: image width and height must be > 0")
	}

	var tfImage [1][][][3]float32

	for j := 0; j < imageHeight; j++ {
		tfImage[0] = append(tfImage[0], make([][3]float32, imageWidth))
	}

	for i := 0; i < imageWidth; i++ {
		for j := 0; j < imageHeight; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			tfImage[0][j][i][0] = convertValue(r)
			tfImage[0][j][i][1] = convertValue(g)
			tfImage[0][j][i][2] = convertValue(b)
		}
	}
	return tf.NewTensor(tfImage)
}

func convertValue(value uint32) float32 {
	return (float32(value >> 8)) / float32(255)
}