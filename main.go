package main

import (
	"encoding/json"
	"fmt"
	"image"
	"mentha/imageClassifier"
	"net/http"
)

func handleFile(w http.ResponseWriter, req *http.Request) {
	file, _, err := req.FormFile("image")
	if err != nil {
		fmt.Fprintf(w, "error on upload file: %v", err.Error())
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(w, "error on upload file: %v", err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	results := imageClassifier.Classify(img)
	json.NewEncoder(w).Encode(map[string][]string{"result": results})
}

func main() {
	http.HandleFunc("/upload", handleFile)
	http.ListenAndServe(":9090", nil)
}
