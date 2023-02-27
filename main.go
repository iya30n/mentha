package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/disintegration/imaging"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"
)

var (
	model  *tg.Model
	labels []string

	modelPath string = "./model/2"
	imagePath string = "./images/canary.jpg"
)

func main() {
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "2")

	loadLabels(modelPath)
	model = tg.LoadModel(modelPath, []string{"serve"}, nil)
	img, err := loadImage(imagePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	classifications := mainHandler(img)
	
	for _, cl := range classifications {
		fmt.Println(cl.Label)
	}
}

func loadLabels(path string) error {
	modelLabels := path + "/labels.txt"
	f, err := os.Open(modelLabels)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type classification struct {
	Label      string
	Proability float32
}

func loadImage(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)

	return img, err
}

func mainHandler(img image.Image) []classification {
	normalizedImg, err := createTensor(img)
	if err != nil {
		log.Fatalf("unable to make a normalizedImg from image: %v", err)
	}

	results := model.Exec(
		[]tf.Output{
			model.Op("StatefulPartitionedCall", 0),
		}, map[tf.Output]*tf.Tensor{
			model.Op("serving_default_input_1", 0): normalizedImg,
		},
	)

	probabilities := results[0].Value().([][]float32)[0]
	classifications := []classification{}
	for i, p := range probabilities {
		if p < -1 {
			continue
		}
		classifications = append(classifications, classification{
			Label:      strings.ToLower(labels[i]),
			Proability: p,
		})
		labelText := strings.ToLower(labels[i])
		fmt.Printf("%s %f \n", labelText, p)
	}

	return classifications
}

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
