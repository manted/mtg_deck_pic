package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func getImage(cardName string, withDownload bool) (image.Image, error) {
	var imgFile *os.File
	var extension string
	var err error
	imgFile, extension, err = openImage(cardName)
	if err != nil {
		// not found
		if !withDownload {
			return nil, err
		}
		err := downloadCardImage(cardName)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// open again
		imgFile, extension, err = openImage(cardName)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	defer imgFile.Close()
	img, err := decodeImage(imgFile, extension)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return img, nil
}

func openImage(cardName string) (*os.File, string, error) {
	// Try .png
	pngName := ImgDir + "/" + cardName + ".png"
	if f, err := os.Open(pngName); err == nil {
		return f, "png", nil
	}
	// Try .jpg
	jpgName := ImgDir + "/" + cardName + ".jpg"
	if f, err := os.Open(jpgName); err == nil {
		return f, "jpg", nil
	}
	// Try .jpeg
	jpegName := ImgDir + "/" + cardName + ".jpeg"
	if f, err := os.Open(jpegName); err == nil {
		return f, "jpeg", nil
	}
	// Not found
	return nil, "", fmt.Errorf("no image found for card name %q with .png, .jpg, or .jpeg extension", cardName)
}

func decodeImage(file *os.File, ext string) (image.Image, error) {
	switch ext {
	case "png":
		return png.Decode(file)
	case "jpg", "jpeg":
		return jpeg.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image extension: %s", ext)
	}
}
