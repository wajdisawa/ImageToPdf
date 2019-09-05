package image_to_pdf

import (
	"ImageToPdf/pkg/config"
	"github.com/disintegration/imaging"
	"github.com/jung-kurt/gofpdf"
	"image"
	"image/jpeg"
	"io"
)

type ImageToPdfConverter struct {
}

func (imgToPdfConverter ImageToPdfConverter) Convert(file io.ReadSeeker, imageType string) error {

	pdf := gofpdf.New(config.PdfOrientation, config.PdfMeasurementUnit, config.PdfPaperType, config.PdfFont)
	pdf.AddPage()
	if seekErr := seekFile(file); seekErr != nil {
		return seekErr
	}
	width, height, dimErr := getDimensions(file)
	if dimErr != nil {
		return dimErr
	}
	options := gofpdf.ImageOptions{}
	if width > height {
		fr, fw := io.Pipe()
		go rotateImage(fw, file)
		width, height = height, width
		options.ImageType = "jpg"

		pdf.RegisterImageOptionsReader(
			config.ImageName,
			options,
			fr)
	} else {
		options.ImageType = imageType
		options.AllowNegativePosition = false
		options.ReadDpi = true
		if seekErr := seekFile(file); seekErr != nil {
			return seekErr
		}
		pdf.RegisterImageOptionsReader(
			config.ImageName,
			options,
			file)
	}
	width, height = adjustToPdfDimensions(width, height)

	pdf.ImageOptions(
		config.ImageName,
		config.Margin,
		config.Margin,
		width,
		height,
		false,
		options,
		0,
		"",
	)
	if pdfErr := pdf.Error(); pdfErr != nil {
		return pdfErr
	}
	pdf.OutputFileAndClose(config.TestFilePath + config.ConvertedFileName)
	return nil
}

func adjustToPdfDimensions(width float64, height float64) (float64, float64) {
	//TODO: add later adjustment based on the paper size
	//what we do here is to make sure that the pdf have at least 127 Ppi
	width, height = width/config.PdfPpiFactor, height/config.PdfPpiFactor
	if width > config.PdfMaxWidth {
		width = config.PdfMaxWidth
		height = 0
	}
	if height > config.PdfMaxHeight {
		height = config.PdfMaxHeight
		width = 0
	}
	return width, height
}

func getDimensions(file io.Reader) (float64, float64, error) {
	imageFile, _, decodeErr := image.DecodeConfig(file)
	if decodeErr != nil {
		return 0, 0, decodeErr
	}
	return float64(imageFile.Width), float64(imageFile.Height), nil
}

func seekFile(file io.ReadSeeker) error {
	_, SeekErr := file.Seek(0, 0)
	if SeekErr != nil {
		return SeekErr
	}
	return nil
}
func rotateImage(fw *io.PipeWriter, file io.ReadSeeker) {
	if seekErr := seekFile(file); seekErr != nil {
		fw.CloseWithError(seekErr)
		return
	}
	srcImage, _, decodeErr := image.Decode(file)
	if decodeErr != nil {
		fw.CloseWithError(decodeErr)
		return
	}
	jpgOptions := jpeg.Options{Quality: config.ImageQuality}
	if encodeErr := jpeg.Encode(fw, imaging.Rotate90(srcImage), &jpgOptions); encodeErr != nil {
		fw.CloseWithError(encodeErr)
	}
	if writerErr := fw.Close(); writerErr != nil {
		return
	}
}
