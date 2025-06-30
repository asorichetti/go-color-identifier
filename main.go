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
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ImageClickOverlay struct {
	widget.BaseWidget
	img        image.Image
	imgSize    fyne.Size
	imgDisplay *canvas.Image
	label      *widget.Label
	swatch     *canvas.Rectangle
}

func NewImageClickOverlay(img image.Image, imgDisplay *canvas.Image, label *widget.Label, swatch *canvas.Rectangle) *ImageClickOverlay {
	overlay := &ImageClickOverlay{
		img:        img,
		imgDisplay: imgDisplay,
		imgSize:    imgDisplay.Size(),
		label:      label,
		swatch:     swatch,
	}
	overlay.ExtendBaseWidget(overlay)
	return overlay
}

func (c *ImageClickOverlay) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(c.imgSize)
	return widget.NewSimpleRenderer(rect)
}

func (c *ImageClickOverlay) MouseDown(ev *desktop.MouseEvent) {
	imgBounds := c.img.Bounds()        // real image size (in pixels)
	displaySize := c.imgDisplay.Size() // actual widget size (on screen)

	// Click position inside widget (not global screen)
	clickX := ev.Position.X
	clickY := ev.Position.Y

	// Map from displayed coordinates to original image pixel coordinates
	ratioX := float64(imgBounds.Dx()) / float64(displaySize.Width)
	ratioY := float64(imgBounds.Dy()) / float64(displaySize.Height)

	px := int(float64(clickX) * ratioX)
	py := int(float64(clickY) * ratioY)

	// Flip Y if image is flipped or offset — may need tweaking
	if px >= 0 && py >= 0 && px < imgBounds.Dx() && py < imgBounds.Dy() {
		col := c.img.At(px, py)
		r, g, b, _ := col.RGBA()
		r8 := uint8(r >> 8)
		g8 := uint8(g >> 8)
		b8 := uint8(b >> 8)
		hex := fmt.Sprintf("#%02X%02X%02X", r8, g8, b8)

		c.label.SetText(fmt.Sprintf("Clicked (%d, %d): R:%d G:%d B:%d (%s)", px, py, r8, g8, b8, hex))
		c.swatch.FillColor = color.RGBA{r8, g8, b8, 255}
		c.swatch.Refresh()
	} else {
		c.label.SetText("Click was outside the image area.")
	}
}

func (c *ImageClickOverlay) MouseUp(ev *desktop.MouseEvent) {}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Click to Identify Color")

	label := widget.NewLabel("Upload an image and click to get the color at that spot.")
	imageContainer := container.NewMax()
	swatch := canvas.NewRectangle(color.White)
	swatch.SetMinSize(fyne.NewSize(100, 100))

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

			imgCanvas := canvas.NewImageFromImage(imgData)
			imgCanvas.FillMode = canvas.ImageFillContain
			imgCanvas.SetMinSize(fyne.NewSize(400, 400))

			// Wait until canvas is sized before overlaying
			myApp.SendNotification(&fyne.Notification{
				Title:   "Ready",
				Content: "Image loaded. You can now click it.",
			})

			overlay := NewImageClickOverlay(imgData, imgCanvas, label, swatch)

			imageContainer.Objects = []fyne.CanvasObject{
				container.NewStack(imgCanvas, overlay),
			}
			imageContainer.Refresh()
		}, myWindow)
	})

	myWindow.SetContent(container.NewVBox(
		label,
		uploadButton,
		imageContainer,
		swatch,
	))

	myWindow.Resize(fyne.NewSize(600, 500))
	myWindow.ShowAndRun()
}
