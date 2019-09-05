package main

import (
	"ImageToPdf/pkg/image_to_pdf"
	"fmt"
	"os"
)

func main() {
	imgPdf := new(image_to_pdf.ImageToPdfConverter)
	imgFile, openErr := os.Open("test_files/dark.jpg")
	if openErr != nil {
		fmt.Println(openErr)
	}
	defer imgFile.Close()

	err := imgPdf.Convert(imgFile, "jpg")
	if err != nil {
		fmt.Println(err)
	}
}
