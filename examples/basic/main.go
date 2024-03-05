package main

import (
	imgui "github.com/AllenDang/cimgui-go"
	ebitenbackend "github.com/damntourists/cimgui-go-ebiten"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MyGame struct {}

func (m MyGame) Update() error {
	//TODO implement me

	return nil
}

func (m MyGame) Draw(screen *ebiten.Image) {
	var tileSize = 32

	newStyle := imgui.NewStyle()
	imgui.StyleColorsDarkV(newStyle)




	var gridColor = imgui.Vec4{
		X: 0,
		Y: 0,
		Z: 0,
		W: 0,
	}
	var gridColorAlt =

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

func (m MyGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	// You MUST use the following buildtags when building:
	// * exclude_cimgui_sdli
	// * exclude_cimgui_glfw
	// TODO: Work in progress... currently not possible to cast to Backend

	adapter := ebitenbackend.NewEbitenAdapter()
	b, _ := imgui.CreateBackend(adapter.Backend())

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	b.CreateWindow("Hello from cimgui-go-ebiten!", 1200, 900)
	adapter.SetGame(MyGame{})
	b.Run(func() {
		_ = ebiten.RunGame(adapter.Game())
	})

}
