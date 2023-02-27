package imageClassifier

import (
	"bufio"
	"os"
)

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