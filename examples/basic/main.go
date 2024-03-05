package main

import (
	imgui "github.com/AllenDang/cimgui-go"
	ebitenbackend "github.com/damntourists/cimgui-go-ebiten"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

var adapter = ebitenbackend.NewEbitenAdapter()

type MyGame struct{}

func (m MyGame) Draw(screen *ebiten.Image) {
	println("mygame draw!")

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

	imgui.Render()
}

func (m MyGame) Update() error {
	//println("mygame update called!")
	imgui.NewFrame()

	imgui.Begin("Demo Window")
	imgui.Text("Hello World!")
	imgui.End()

	defer imgui.EndFrame()
	return nil
}

func (m MyGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	println("mygame layout called!")
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
	// You MUST use the following buildtags when building:
	// * exclude_cimgui_sdli
	// * exclude_cimgui_glfw
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	adapter.CreateWindow("Hello from cimgui-go-ebiten!", 800, 600)
	adapter.SetGame(MyGame{})
	adapter.Run(func() {
		_ = ebiten.RunGame(adapter.Game())
	})

}
