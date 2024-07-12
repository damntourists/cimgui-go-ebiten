package main

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/damntourists/cimgui-go-ebiten"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"unsafe"
)

const (
	fontFilename = "Sharpie_Complete/Fonts/WEB/fonts/Sharpie-Regular.ttf"
)

var (
	fontSize = float32(math.Floor(24 * ebiten.Monitor().DeviceScaleFactor()))
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
	imgui.Begin("Custom Font")
	imgui.Text("Hello world!")
	imgui.End()
	return nil
}

func (m *MyGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func rebuildFonts() {
	// Clear out any existing fonts. First font to be added is considered "default"
	fontAtlas := imgui.CurrentIO().Fonts()
	if fontAtlas.FontCount() > 0 {
		fontAtlas.Clear()
	}

	// Read embedded font
	fontData, err := webFonts.ReadFile(fontFilename)
	if err != nil {
		panic(err)
	}

	dataPtr := uintptr(unsafe.Pointer(imgui.SliceToPtr(fontData)))
	dataLen := int32(len(fontData))

	cfg := imgui.NewFontConfig()

	ranges := fontAtlas.GlyphRangesDefault()

	if fontAtlas.FontCount() > 0 {
		cfg.SetFontDataOwnedByAtlas(true)
	}

	// Force font atlas to rebuild tex cache
	_ = fontAtlas.AddFontFromMemoryTTFV(dataPtr, dataLen, fontSize, cfg, ranges)
	_, _, _, _ = fontAtlas.GetTextureDataAsRGBA32()
}

func main() {
	backend := ebitenbackend.NewEbitenBackend()
	backend.SetGame(&MyGame{})
	backend.SetWindowFlags(
		ebitenbackend.EbitenWindowFlagsResizingMode,
		int(ebiten.WindowResizingModeEnabled),
	)

	// register backend with imgui
	_, _ = imgui.CreateBackend[ebitenbackend.EbitenWindowFlags](backend)

	backend.CreateWindow("Hello from cimgui-go-ebiten!", 800, 600)

	backend.Run(func() {
		rebuildFonts()
		_ = ebiten.RunGame(backend.Game())
	})
}
