package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)

type spritesheet struct {
	imagesDef
	framesDef
	animFrameDef
	texturepackerDef
}

type texturepackerDef struct {
	Texturepacker []string `json:"texturepacker"`
}

type imagesDef struct {
	Images []string `json:"images"`
}

type framesDef struct {
	Frames [][]int `json:"frames"`
}

type animFrameDef struct {
	Animations map[string]frameDef `json:"animations"`
}

type frameDef struct {
	Frames []int `json:"frames"`
}

func main() {

	ssImageFile := os.Args[1]
	ssJSONFile := os.Args[2]
	outFolder := os.Args[3]

	outDir, err := filepath.Abs(outFolder)
	if err != nil {
		panic(err)
	}
	fmt.Println(outDir)

	img, err := loadPNG(ssImageFile)
	if err != nil {
		panic(err)
	}

	sheet, err := ioutil.ReadFile(ssJSONFile)
	if err != nil {
		panic(err)
	}

	var sheetData spritesheet
	if err := json.Unmarshal([]byte(sheet), &sheetData); err != nil {
		panic(err)
	}

	imWidth := img.Bounds().Dx()
	imHeight := img.Bounds().Dy()
	sheetImage := image.NewRGBA(img.Bounds())

	for x := 0; x < imWidth; x++ {
		for y := 0; y < imHeight; y++ {
			oldColor := img.At(x, y)
			newColor := color.RGBAModel.Convert(oldColor)
			sheetImage.Set(x, y, newColor)
		}
	}

	subImageRect := image.Rect(0, 0, 0, 0)
	subImageMaxPt := image.Point{}
	subImageMinPt := image.Point{}

	for k, v := range sheetData.Animations {
		frame := sheetData.Frames[v.Frames[0]]
		x := frame[0]
		y := frame[1]
		width := frame[2]
		height := frame[3]
		subImageMinPt.X = x
		subImageMinPt.Y = y
		subImageMaxPt.X = x + width
		subImageMaxPt.Y = y + height
		subImageRect.Min = subImageMinPt
		subImageRect.Max = subImageMaxPt
		subImg := sheetImage.SubImage(subImageRect)
		createSpriteFile(outDir+"/"+k+".png", subImg)
	}
}

func createSpriteFile(name string, imgSrc image.Image) {
	outfile, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	png.Encode(outfile, imgSrc)
}

func loadPNG(path string) (image.Image, error) {
	fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}
