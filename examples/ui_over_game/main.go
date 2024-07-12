package main

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/damntourists/cimgui-go-ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type MyGame struct{}

func (m *MyGame) Draw(screen *ebiten.Image) {
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

func (m *MyGame) Update() error {
	imgui.ShowDemoWindow()
	return nil
}

func (m *MyGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	backend := ebitenbackend.NewEbitenBackend()
	backend.SetGame(&MyGame{})

	backend.CreateWindow("Hello from cimgui-go-ebiten!", 800, 600)
	backend.Run(func() {
		_ = ebiten.RunGame(backend.Game())
	})

}
