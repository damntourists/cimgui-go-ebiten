package main

import (
	"bytes"
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/damntourists/cimgui-go-ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"image"
	"log"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

const mosaicRatio = 16

var (
	gophersImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Gophers_jpg))
	if err != nil {
		log.Fatal(err)
	}
	gophersImage = ebiten.NewImageFromImage(img)
}

type MyGame struct {
	backend             *ebitenbackend.EbitenBackend
	gophersRenderTarget *ebiten.Image
}

func (m *MyGame) Draw(screen *ebiten.Image) {
	// Shrink the image once.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1.0/mosaicRatio, 1.0/mosaicRatio)
	m.gophersRenderTarget.DrawImage(gophersImage, op)

	// Enlarge the shrunk image.
	// The filter is the nearest filter, so the result will be mosaic.
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(mosaicRatio, mosaicRatio)
	screen.DrawImage(m.gophersRenderTarget, op)
}

func (m *MyGame) Update() error {
	imgui.ShowDemoWindow()
	imgui.Begin("Game In UI")

	Image(m.backend.Game().ScreenTextureID(), imgui.Vec2{X: screenWidth, Y: screenHeight})
	imgui.End()
	return nil
}

func (m *MyGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func Image(tid imgui.TextureID, size imgui.Vec2) {
	uv0 := imgui.NewVec2(0, 0)
	uv1 := imgui.NewVec2(1, 1)
	border_col := imgui.NewVec4(0, 0, 0, 0)
	tint_col := imgui.NewVec4(1, 1, 1, 1)
	imgui.ImageV(tid, size, uv0, uv1, tint_col, border_col)
}

func main() {
	w, h := gophersImage.Bounds().Dx(), gophersImage.Bounds().Dy()

	backend := ebitenbackend.NewEbitenBackend()
	backend.SetGame(
		&MyGame{
			backend:             backend,
			gophersRenderTarget: ebiten.NewImage(w/mosaicRatio, h/mosaicRatio),
		},
	)
	backend.SetWindowFlags(
		ebitenbackend.EbitenWindowFlagsResizingMode,
		int(ebiten.WindowResizingModeEnabled),
	)
	backend.SetGameScreenSize(imgui.Vec2{X: screenWidth, Y: screenHeight})

	renderDestination := ebiten.NewImage(320, 240)
	backend.SetGameRenderDestination(renderDestination)

	backend.CreateWindow("Hello from cimgui-go-ebiten!", 800, 600)
	backend.Run(func() {
		_ = ebiten.RunGame(backend.Game())
	})
}
