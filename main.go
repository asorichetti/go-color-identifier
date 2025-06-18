package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/nfnt/resize"
)

func getDominantColor(img image.Image) color.Color {
	smallImg := resize.Resize(100, 0, img, resize.Lanczos3)

	var rTotal, gTotal, bTotal uint64
	var count uint64

	for y := 0; y < smallImg.Bounds().Dy(); y++ {
		for x := 0; x < smallImg.Bounds().Dx(); x++ {
			r, g, b, _ := smallImg.At(x, y).RGBA()
			rTotal += uint64(r >> 8)
			gTotal += uint64(g >> 8)
			bTotal += uint64(b >> 8)
			count++
		}
	}

	rAvg := rTotal / count
	gAvg := gTotal / count
	bAvg := bTotal / count

	return color.RGBA{uint8(rAvg), uint8(gAvg), uint8(bAvg), 255}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Color Identifier")

	label := widget.NewLabel("Upload an image to find the dominant color.")
	imageBox := container.NewMax()
	colorLabel := widget.NewLabel("")

	uploadButton := widget.NewButton("Upload Image", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			imgData, _, err := image.Decode(reader)
			if err != nil {
				label.SetText("Failed to decode image.")
				return
			}

			domColor := getDominantColor(imgData)

			rgba := domColor.(color.RGBA)
			hex := fmt.Sprintf("#%02X%02X%02X", rgba.R, rgba.G, rgba.B)

			// Update UI
			label.SetText(fmt.Sprintf("Dominant color: R:%d G:%d B:%d (%s)", rgba.R, rgba.G, rgba.B, hex))
			colorBox := canvas.NewRectangle(domColor)
			colorBox.SetMinSize(fyne.NewSize(100, 100))

			imageBox.Objects = []fyne.CanvasObject{colorBox}
			imageBox.Refresh()
		}, myWindow)
	})

	myWindow.SetContent(container.NewVBox(
		label,
		uploadButton,
		imageBox,
		colorLabel,
	))

	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()
}
