package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/exp/rand"
)

const (
	SignedImgWidth  = 120
	SignedImgHeight = 168

	SignedCol = 20
	SignedRow = 5
)

func main1() {
	cards := []string{
		"Abrade",
		"Abrupt Decay",
		"Ancestral Vision",
		"Ancient Grudge",
		"Baleful Strix",
		"Blood Crypt",
		"Bloodbraid Elf",
		"Bloodstained Mire",
		"Blooming Marsh",
		"Blue Elemental Blast",
		"Bojuka Bog",
		"Brainstorm",
		"Brazen Borrower",
		"Breeding Pool",
		"Cabal Therapy",
		"Chain Lightning",
		"Choke",
		"Collector Ouphe",
		"Concealed Courtyard",
		"Containment Priest",
		"Courser of Kruphix",
		"Dack Fayden",
		"Dark Confidant",
		"Daze",
		"Death's Shadow",
		"Deathrite Shaman",
		"Delver of Secrets",
		"Dismember",
		"Dragon's Rage Channeler",
		"Dread of Night",
		"Drown in the Loch",
		"Engineered Explosives",
		"Ethersworn Canonist",
		"Fatal Push",
		"Fire::Ice",
		"Flooded Strand",
		"Flusterstorm",
		"Force of Will",
		"Forest",
		"Godless Shrine",
		"Grim Flayer",
		"Huntmaster of the Fells",
		"Hydroblast",
		"Ice-Fang Coatl",
		"Island",
		"Izzet Staticaster",
		"Jace, the Mind Sculptor",
		"Klothys, God of Destiny",
		"Leovold, Emissary of Trest",
		"Life from the Loam",
		"Lightning Bolt",
		"Liliana of the Veil",
		"Liliana, the Last Hope",
		"Marsh Flats",
		"Meddling Mage",
		"Misty Rainforest",
		"Mountain",
		"Narset, Parter of Veils",
		"Nimble Mongoose",
		"Null Rod",
		"Overgrown Tomb",
		"Pernicious Deed",
		"Pithing Needle",
		"Plains",
		"Polluted Delta",
		"Ponder",
		"Preordain",
		"Punishing Fire",
		"Pyroblast",
		"Red Elemental Blast",
		"Rest in Peace",
		"Scalding Tarn",
		"Scavenging Ooze",
		"Shambling Vent",
		"Shardless Agent",
		"Shatterstorm",
		"Siege Rhino",
		"Snapcaster Mage",
		"Snow-Covered Forest",
		"Snow-Covered Island",
		"Snow-Covered Plains",
		"Snow-Covered Swamp",
		"Spirebluff Canal",
		"Steam Vents",
		"Stomping Ground",
		"Surgical Extraction",
		"Swamp",
		"Swords to Plowshares",
		"Temple Garden",
		"Thoughtseize",
		"Tireless Tracker",
		"Toxic Deluge",
		"Twilight Mire",
		"Veil of Summer",
		"Wasteland",
		"Watery Grave",
		"Windswept Heath",
		"Winter Orb",
		"Wooded Foothills",
		"Yorion, Sky Nomad",
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	for i := range cards {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}

	bottomRightPoint := image.Point{SignedImgWidth * SignedCol, SignedImgHeight * SignedRow}
	// whole image container
	r := image.Rectangle{image.Point{0, 0}, bottomRightPoint}
	rgba := image.NewRGBA(r)
	// render main deck
	for i, cardName := range cards {
		row := i % SignedCol
		col := i / SignedCol
		startingPoint := image.Point{row * SignedImgWidth, col * SignedImgHeight}
		bound := image.Rectangle{startingPoint, startingPoint.Add(image.Point{SignedImgWidth, SignedImgHeight})}
		imgFile, err := os.Open(fmt.Sprintf("img/%s.png", cardName))
		if err != nil {
			fmt.Println(err)
			continue
		}
		img, err := png.Decode(imgFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cardImg := resize.Resize(SignedImgWidth, SignedImgHeight, img, resize.Lanczos3)
		draw.Draw(rgba, bound, cardImg, image.Point{0, 0}, draw.Src)
	}
	// create final image
	outputFileName := "./decklist/signed_cards.jpg"
	out, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println(err)
	}

	var opt jpeg.Options
	opt.Quality = 60

	jpeg.Encode(out, rgba, &opt)
	fmt.Printf("Successfully generated %s\n", outputFileName)
}
