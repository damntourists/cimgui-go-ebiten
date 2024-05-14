package ebitenbackend

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameProxy struct {
	game    ebiten.Game
	adapter *EbitenAdapter

	width, height             float64
	screenWidth, screenHeight int

	gameScreenTextureID imgui.TextureID
	gameScreen          *ebiten.Image

	filter ebiten.Filter

	clipRegion imgui.Vec2

	Resizeable bool
}

func (g *GameProxy) Game() ebiten.Game {
	return g.game
}

// Update - Update UI and game in tandem. Handle inputs
func (g *GameProxy) Update() error {
	if g.game == nil {
		panic("No game to update!")
	}

	io := imgui.CurrentIO()

	// Sync keyboard
	CurrentAdapter.inputChars = sendInput(imgui.CurrentIO(), CurrentAdapter.inputChars)

	// Sync mouse wheel
	xoff, yoff := ebiten.Wheel()
	io.AddMouseWheelDelta(float32(xoff), float32(yoff))

	// Sync mouse position
	cx, cy := ebiten.CursorPosition()
	io.SetMousePos(imgui.Vec2{X: float32(cx), Y: float32(cy)})

	// Sync mouse button
	io.SetMouseButtonDown(0, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
	io.SetMouseButtonDown(1, ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
	io.SetMouseButtonDown(2, ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle))

	// Sync mouse cursor
	switch imgui.CurrentMouseCursor() {
	case imgui.MouseCursorNone:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	case imgui.MouseCursorArrow:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	case imgui.MouseCursorTextInput:
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	case imgui.MouseCursorResizeAll:
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	case imgui.MouseCursorResizeEW:
		ebiten.SetCursorShape(ebiten.CursorShapeEWResize)
	case imgui.MouseCursorResizeNS:
		ebiten.SetCursorShape(ebiten.CursorShapeNSResize)
	case imgui.MouseCursorHand:
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	default:
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	io.SetDeltaTime(1.0 / 60.0)

	if io.Fonts().FontCount() == 0 {
		// The font atlas is empty. Add the default font and set up
		// scaling to match ebiten's scale factor. It's recommended to set up fonts on
		// your own (per imgui docs) because the default font does not scale well.
		io.SetFontGlobalScale(float32(ebiten.Monitor().DeviceScaleFactor()))

		io.Fonts().AddFontDefault()
	}

	if !io.Fonts().IsBuilt() {
		_, _, _, _ = io.Fonts().GetTextureDataAsRGBA32()
	}

	imgui.NewFrame()
	err := g.game.Update()
	imgui.EndFrame()
	return err
}

func (g *GameProxy) renderable() bool {
	if g.width > 0 && g.height > 0 && g.screenWidth > 0 && g.screenHeight > 0 {
		return true
	}
	return false
}

func (g *GameProxy) Draw(screen *ebiten.Image) {
	destination := screen
	if g.gameScreen != nil {
		destination = g.gameScreen
	}
	destination.Clear()
	g.game.Draw(destination)

	imgui.Render()

	if CurrentAdapter.ClipMask {
		if CurrentAdapter.lmask == nil {
			w, h := screen.Size()
			CurrentAdapter.lmask = ebiten.NewImage(w, h)
		} else {
			w1, h1 := screen.Size()
			w2, h2 := CurrentAdapter.lmask.Size()
			if w1 != w2 || h1 != h2 {
				CurrentAdapter.lmask.Dispose()
				CurrentAdapter.lmask = ebiten.NewImage(w1, h1)
			}
		}
		RenderMasked(screen, CurrentAdapter.lmask, imgui.CurrentDrawData(), Cache, g.filter)
	} else {
		Render(screen, imgui.CurrentDrawData(), Cache, g.filter)
	}

}

func (g *GameProxy) Layout(outsideWidth, outsideHeight int) (int, int) {
	// This is the full window dimensions
	width := float64(outsideWidth)
	height := float64(outsideHeight)

	io := imgui.CurrentIO()
	io.SetDisplaySize(imgui.Vec2{X: float32(width), Y: float32(height)})

	// Set game screen height/width to match wrapped game
	//screenWidth, screenHeight := g.game.Layout(outsideWidth, outsideHeight)
	//g.screenWidth = screenWidth   //g.Screen().Bounds().Dx()
	//g.screenHeight = screenHeight //g.Screen().Bounds().Dy()

	return int(width), int(height)
	//return g.screenWidth, g.screenHeight
}

func (g *GameProxy) Screen() *ebiten.Image {
	return g.gameScreen
}

func (g *GameProxy) ScreenTextureID() imgui.TextureID {
	return g.gameScreenTextureID
}

func (g *GameProxy) SetGameScreenSize(v imgui.Vec2) {
	g.height = float64(max(1, v.Y))
	g.width = float64(max(1, v.X))
}
