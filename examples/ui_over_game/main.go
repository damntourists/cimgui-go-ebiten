package main

import (
	imgui "github.com/AllenDang/cimgui-go"
	backend "github.com/damntourists/cimgui-go-ebiten"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

var adapter = backend.NewEbitenAdapter()

type MyGame struct{}

func (m MyGame) Draw(screen *ebiten.Image) {
	var tileSize = 32

	var gridColor = color.RGBA{R: 100, G: 100, B: 100, A: 1}
	var gridColorAlt = color.RGBA{R: 150, G: 150, B: 150, A: 1}

	w := screen.Bounds().Dx()
	h := screen.Bounds().Dy()

	for y := 0; y <= (h / tileSize); y++ {
		for x := 0; x <= (w / tileSize); x++ {
			c := gridColor
			if (x%2 == 0 && y%2 == 1) || (x%2 == 1 && y%2 == 0) {
				c = gridColorAlt
			}

			vector.DrawFilledRect(screen,
				float32(x*tileSize),
				float32(y*tileSize),
				float32(tileSize),
				float32(tileSize), c, false)
		}
	}
}

func (m MyGame) Update() error {
	imgui.ShowDemoWindow()
	return nil
}

func (m MyGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
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
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	adapter.CreateWindow("Hello from cimgui-go-ebiten!", 800, 600)
	adapter.SetGame(MyGame{})
	adapter.Run(func() {
		_ = ebiten.RunGame(adapter.Game())
	})

}
