package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

const (
	DefaultDeckName             = "grixis_delver"
	MaindeckSize60              = 60
	MaindeckSize80              = 80
	NumOfSideboardCards         = 15
	Col60                       = 12
	Col80                       = 16
	ImgWidth                    = 160
	ImgHeight                   = 224
	SideboardGap                = 20
	NumOfSizeboardFirstRowCards = 7
)

type Deck struct {
	Maindeck   []string
	Sideboard  []string
	CardImgMap map[string]image.Image
}

func readDecklist(deckName string) (*Deck, error) {
	file, err := os.Open(fmt.Sprintf("decklist/%s.txt", deckName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	deck := &Deck{
		CardImgMap: map[string]image.Image{},
	}
	scanner := bufio.NewScanner(file)
	decklist := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "//") {
			continue
		}
		lineParts := strings.Split(line, " ")
		numOfCopiesStr := lineParts[0]
		numOfCopies, err := strconv.Atoi(numOfCopiesStr)
		if err != nil {
			return nil, err
		}
		cardName := line[len(numOfCopiesStr)+1:]
		for i := 0; i < numOfCopies; i++ {
			decklist = append(decklist, cardName)
		}
		_, ok := deck.CardImgMap[cardName]
		if !ok {
			imgFile, extension, err := openImage(cardName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer imgFile.Close()
			img, err := decodeImage(imgFile, extension)
			if err != nil {
				fmt.Println(err)
				continue
			}
			resizedImg := resize.Resize(ImgWidth, ImgHeight, img, resize.Lanczos3)
			deck.CardImgMap[cardName] = resizedImg
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	totalCards := len(decklist)
	deck.Maindeck = decklist[:totalCards-NumOfSideboardCards]
	deck.Sideboard = decklist[totalCards-NumOfSideboardCards:]
	return deck, nil
}

func openImage(cardName string) (*os.File, string, error) {
	// Try .png
	pngName := "img/" + cardName + ".png"
	if f, err := os.Open(pngName); err == nil {
		return f, "png", nil
	}
	// Try .jpg
	jpgName := "img/" + cardName + ".jpg"
	if f, err := os.Open(jpgName); err == nil {
		return f, "jpg", nil
	}
	// Try .jpeg
	jpegName := "img/" + cardName + ".jpeg"
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

func main() {
	args := os.Args
	deckName := DefaultDeckName
	if len(args) > 1 {
		deckName = args[1]
	}
	fmt.Printf("Generating deck pic for %s\n", deckName)
	deck, err := readDecklist(deckName)
	if err != nil {
		fmt.Println(err)
		return
	}
	maindeckSize := len(deck.Maindeck)
	col := Col60
	if maindeckSize == MaindeckSize80 {
		col = Col80
	}
	sideboardHeight := ImgHeight * 2
	bottomRightPoint := image.Point{ImgWidth * col, ImgHeight*maindeckSize/col + sideboardHeight + SideboardGap}
	// whole image container
	r := image.Rectangle{image.Point{0, 0}, bottomRightPoint}
	rgba := image.NewRGBA(r)
	// render background
	bgImgFile, err := os.Open("img/background.jpg")
	if err != nil {
		fmt.Println(err)
	} else {
		bgImg, err := jpeg.Decode(bgImgFile)
		if err != nil {
			fmt.Println(err)
		} else {
			resizedBgImg := resize.Resize(uint(r.Bounds().Dx()), 0, bgImg, resize.Lanczos3)
			startingPoint := image.Point{0, r.Dy() - resizedBgImg.Bounds().Dy()}
			bound := image.Rectangle{startingPoint, startingPoint.Add(resizedBgImg.Bounds().Size())}
			draw.Draw(rgba, bound, resizedBgImg, image.Point{0, 0}, draw.Src)
		}
	}
	// render main deck
	for i, cardName := range deck.Maindeck {
		row := i % col
		col := i / col
		startingPoint := image.Point{row * ImgWidth, col * ImgHeight}
		bound := image.Rectangle{startingPoint, startingPoint.Add(image.Point{ImgWidth, ImgHeight})}
		cardImg, ok := deck.CardImgMap[cardName]
		if ok {
			draw.Draw(rgba, bound, cardImg, image.Point{0, 0}, draw.Src)
		}
	}
	// render sideboard
	maindeckRows := maindeckSize / col
	sideboardGapLeft := (r.Dx() - ImgWidth*NumOfSizeboardFirstRowCards) / 2
	for i, cardName := range deck.Sideboard[:NumOfSizeboardFirstRowCards] {
		startingPoint := image.Point{sideboardGapLeft + ImgWidth*i, maindeckRows*ImgHeight + SideboardGap}
		bound := image.Rectangle{startingPoint, startingPoint.Add(image.Point{ImgWidth, ImgHeight})}
		cardImg, ok := deck.CardImgMap[cardName]
		if ok {
			draw.Draw(rgba, bound, cardImg, image.Point{0, 0}, draw.Src)
		}
	}
	sideboardGapLeft = (r.Dx() - ImgWidth*8) / 2
	for i, cardName := range deck.Sideboard[NumOfSizeboardFirstRowCards:] {
		startingPoint := image.Point{sideboardGapLeft + ImgWidth*i, maindeckRows*ImgHeight + SideboardGap + ImgHeight}
		bound := image.Rectangle{startingPoint, startingPoint.Add(image.Point{ImgWidth, ImgHeight})}
		cardImg, ok := deck.CardImgMap[cardName]
		if ok {
			draw.Draw(rgba, bound, cardImg, image.Point{0, 0}, draw.Src)
		}
	}
	// create final image
	outputFileName := fmt.Sprintf("./decklist/%s.jpg", deckName)
	out, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println(err)
	}

	var opt jpeg.Options
	opt.Quality = 60

	jpeg.Encode(out, rgba, &opt)
	fmt.Printf("Successfully generated %s\n", outputFileName)
}
