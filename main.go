package main

import (
	"bufio"
	"image"
	"log"
	"os"

	tg "github.com/galeone/tfgo"
)

var (
	model  *tg.Model
	labels []string

	modelPath string = "/your/path"
	imagePath string = "/your/path"
)

func main() {
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "2")

	loadLabels(modelPath)
	model = tg.LoadModel(modelPath, []string{"serve"}, nil)
	img, err := loadImage(imagePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	mainHandler(img)
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

func loadImage(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)

	return img, err
}

func mainHandler(img image.Image) {}
