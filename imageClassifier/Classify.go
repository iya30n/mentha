package imageClassifier

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"
)

var (
	model  *tg.Model
	labels []string
)

type classification struct {
	Label      string
	Proability float32
}

func Classify(img image.Image)[]string {
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "2")

	modelPath := "./model/2"

	loadLabels(modelPath)
	model = tg.LoadModel(modelPath, []string{"serve"}, nil)

	classifications := mainHandler(img)

	var results []string
	for _, cl := range classifications {
		results = append(results, cl.Label)
	}

	return results
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
