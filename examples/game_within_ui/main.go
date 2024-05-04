package main

import (
	"bytes"
	imgui "github.com/AllenDang/cimgui-go"
	backend "github.com/damntourists/cimgui-go-ebiten"
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
	adapter      = backend.NewEbitenAdapter()
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

	Image(adapter.ScreenTextureID(), imgui.Vec2{X: screenWidth, Y: screenHeight})
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
	//
	// The build tags listed below are required to compile with AllenDang/cimgui-go. You
	// may, however, use the damntourists/cimgui-go-lite to bypass this requirement.
	// Please refer to the go.mod file for more info.
	//
	// * exclude_cimgui_sdl
	// * exclude_cimgui_glfw
	//
	w, h := gophersImage.Bounds().Dx(), gophersImage.Bounds().Dy()
	g := &MyGame{
		gophersRenderTarget: ebiten.NewImage(w/mosaicRatio, h/mosaicRatio),
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	renderDestination := ebiten.NewImage(320, 240)
	adapter.CreateWindow("Hello from cimgui-go-ebiten!", 800, 600)
	adapter.SetGame(g)
	adapter.SetGameScreenSize(imgui.Vec2{X: screenWidth, Y: screenHeight})
	adapter.SetGameRenderDestination(renderDestination)
	adapter.Run(func() {
		_ = ebiten.RunGame(adapter.Game())
	})
}
