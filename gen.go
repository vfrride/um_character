/*
 * Overlays character information on the pdf
 *
 * Run as: go run gen.go <input.json> 
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"github.com/unidoc/unidoc/pdf/model/fonts"
)

type Character struct {
	Name         string     `json:"name"`
	Player       string     `json:"player"`
	Calling      string     `json:"calling"`
	Age          int        `json:"age"`
	WealthRating int        `json:"wealth_rating"`
	Vitality     int        `json:"vitality"`
	Coordination int        `json:"coordination"`
	Wit          int        `json:"wit"`
	Intellect    int        `json:"intellect"`
	Charm        int        `json:"charm"`
	Will         int        `json:"will"`
	Prowess      int        `json:"prowess"`
	Actions      int        `json:"actions"`
	Experience   int        `json:"experience"`
	Corruption   Corruption `json:"corruption"`
	Skills       []Skill    `json:"skills"`
	Qualities    []string   `json:"qualities"`
	Impediments  []string   `json:"impediments"`
	Features     []string   `json:"features"`
	Wounds       []string   `json:"wounds"`
	Armour       []string   `json:"armour"`
	Weapons      []Weapon   `json:"weapons"`
	Possessions  []string   `json:"possessions"`
}

type Corruption struct {
	Physical struct {
		Affliction string `json:"affliction"`
		Value      int    `json:"value"`
	}
	Desire struct {
		Affliction string `json:"affliction"`
		Value      int    `json:"value"`
	}
	Drive struct {
		Affliction string `json:"affliction"`
		Value      int    `json:"value"`
	}
}

type Skill struct {
	Name   string   `json:"name"`
	Rating int      `json:"rating"`
	Values []string `json:"values"`
}

type Weapon struct {
	Name     string `json:"name"`
	Rating   int    `json:"rating"`
	Damage   int    `json:"damage"`
	Range    int    `json:"range"`
	Ammo     int    `json:"ammo"`
	Cost     string `json:"cost"`
	AmmoCost string `json:"ammo_cost"`
	Features string `json:"features"`
}

type StringListConfig struct {
	XPos        float64
	YPos        float64
	YOffset     float64
	WrapAt      int
	WrapXOffset float64
	FontSize    float64
}

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Requires at least 1 argument: input path\n")
		fmt.Printf("Usage: go run gen.go <input>.json\n")
		os.Exit(0)
	}

	outputPath := ""
	inputPath := ""

	// Sanity check the input arguments.
	for i, arg := range os.Args {
		if i == 0 {
			continue
		} else if i == 1 {
			inputPath = arg
			outputPath = generateOutputPath(inputPath)
			continue
		}
	}

	character := readJson(inputPath)

	pdfReader, err := readPdf("CharacterSheets.pdf")

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		os.Exit(1)
	}

	c := creator.New()

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			os.Exit(1)
		}
		err = c.AddPage(page)
		if err != nil {
			os.Exit(1)
		}

		if i == 0 {
			_ = c.Draw(getTextNode(character.Name, 293, 35))
			_ = c.Draw(getTextNode(character.Player, 445, 35))
			_ = c.Draw(getTextNode(character.Calling, 302, 60))
			_ = c.Draw(getTextNodeFromInt(character.Age, 287, 85))
			_ = c.Draw(getTextNodeFromInt(character.WealthRating, 488, 85))
			_ = c.Draw(getTextNodeFromInt(character.Vitality, 145, 155))
			_ = c.Draw(getTextNodeFromInt(character.Coordination, 145, 183))
			_ = c.Draw(getTextNodeFromInt(character.Wit, 145, 210))
			_ = c.Draw(getTextNodeFromInt(character.Intellect, 145, 237))
			_ = c.Draw(getTextNodeFromInt(character.Charm, 145, 264))
			_ = c.Draw(getTextNodeFromInt(character.Will, 145, 292))
			_ = c.Draw(getTextNodeFromInt(character.Prowess, 103, 330))
			_ = c.Draw(getTextNodeFromInt(character.Actions, 190, 330))
			_ = c.Draw(getTextNodeFromInt(character.Corruption.Physical.Value, 190, 393))
			_ = c.Draw(getTextNode(character.Corruption.Physical.Affliction, 106, 410))
			_ = c.Draw(getTextNodeFromInt(character.Corruption.Desire.Value, 190, 413))
			_ = c.Draw(getTextNode(character.Corruption.Desire.Affliction, 106, 430))
			_ = c.Draw(getTextNodeFromInt(character.Experience, 103, 571))

			for index, skill := range character.Skills {
				var xPos float64 = 230
				var yPos float64 = float64(147) + float64(index)*64.5
				if index > 6 {
					xPos = float64(415)
					yPos = float64(147) + float64(index-7)*64.5
				}
				c.Draw(getTextNode(skill.Name, xPos, yPos))
				c.Draw(getTextNodeFromInt(skill.Rating, (xPos + 136), (yPos - 2.5)))
				for idx, value := range skill.Values {
					var vYPos float64 = yPos + float64(12) + float64(idx)*9
					c.Draw(getTextNode(value, xPos+25, vYPos, float64(8)))
				}
			}

			config := StringListConfig{float64(17), float64(628), float64(10), 12, float64(88), float64(8)}
			outputStringList(c, character.Qualities, config)

			config.XPos = float64(198)
			outputStringList(c, character.Impediments, config)

			config.XPos = float64(382)
			outputStringList(c, character.Features, config)
		} else if i == 1 {
			config := StringListConfig{float64(72), float64(308), float64(10), 12, float64(88), float64(8)}
			outputStringList(c, character.Wounds, config)

			config.XPos = float64(357)
			outputStringList(c, character.Armour, config)

			outputWeaponsList(c, character.Weapons)

			config = StringListConfig{float64(10), float64(604), float64(11), 12, float64(88), float64(9)}
			outputStringList(c, character.Possessions, config)
		}
	}

	err = c.WriteToFile(outputPath)
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func outputStringList(c *creator.Creator, strs []string, config StringListConfig) {
	for index, str := range strs {
		var xPos float64 = config.XPos
		var yPos float64 = config.YPos + float64(index)*config.YOffset
		if index > config.WrapAt {
			xPos = config.XPos + config.WrapXOffset
			yPos = float64(config.YPos) + float64(index-config.WrapAt-1)*config.YOffset
		}
		c.Draw(getTextNode(str, xPos+25, yPos, config.FontSize))
	}
}

func outputWeaponsList(c *creator.Creator, weapons []Weapon) {
	for index, weapon := range weapons {
		var xPos float64 = float64(38)
		var yPos float64 = float64(411) + float64(index)*11
		c.Draw(getTextNode(weapon.Name, xPos, yPos, float64(9)))
		c.Draw(getTextNodeFromInt(weapon.Rating, float64(xPos+128), yPos, float64(9)))
		c.Draw(getTextNode(fmt.Sprintf("+ %d", weapon.Damage), float64(xPos+198), yPos, float64(9)))
		c.Draw(getTextNode(fmt.Sprintf("%d feet", weapon.Range), float64(xPos+241), yPos, float64(9)))
		c.Draw(getTextNodeFromInt(weapon.Ammo, float64(xPos+295), yPos, float64(9)))
		c.Draw(getTextNode(weapon.Cost, float64(xPos+323), yPos, float64(9)))
		c.Draw(getTextNode(weapon.AmmoCost, float64(xPos+363), yPos, float64(9)))
		c.Draw(getTextNode(weapon.Features, float64(xPos+450), yPos, float64(9)))
	}
}

func getTextNode(text string, posX float64, posY float64, font_size_optional ...float64) creator.Drawable {
	p := creator.NewParagraph(text)
	// Change to times bold font (default is helvetica).
	p.SetFont(fonts.NewFontTimesBold())
	p.SetPos(posX, posY)
	if len(font_size_optional) == 1 {
		p.SetFontSize(font_size_optional[0])
	}
	return p
}

func getTextNodeFromInt(val int, posX float64, posY float64, font_size_optional ...float64) creator.Drawable {
	return getTextNode(strconv.Itoa(val), posX, posY, font_size_optional...)
}

func generateOutputPath(inputPath string) string {
	input := strings.Split(inputPath, ".")
	if input[len(input)-1] == "json" {
		input = input[0:(len(input) - 1)]
	}
	return strings.Join(input, ".") + ".pdf"
}

func readJson(inputPath string) Character {
	raw, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Character

	json.Unmarshal(raw, &c)
	return c
}

func readPdf(pdfPath string) (pdfReader *pdf.PdfReader, err error) {
	f, err := os.Open(pdfPath)
	if err != nil {
		return
	}

	pdfReader, err = pdf.NewPdfReader(f)

	if err != nil {
		return
	}

	return pdfReader, nil
}

func processPage(page *pdf.PdfPage) error {
	mBox, err := page.GetMediaBox()
	if err != nil {
		return err
	}
	pageWidth := mBox.Urx - mBox.Llx
	pageHeight := mBox.Ury - mBox.Lly

	fmt.Printf(" Page: %+v\n", page)
	fmt.Printf(" Page mediabox: %+v\n", page.MediaBox)
	fmt.Printf(" Page height: %f\n", pageHeight)
	fmt.Printf(" Page width: %f\n", pageWidth)

	return nil
}
