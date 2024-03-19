package ebitenbackend

import (
	"fmt"
	imgui "github.com/damntourists/cimgui-go-lite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameProxy struct {
	game    ebiten.Game
	adapter *EbitenAdapter

	width, height             float64
	screenWidth, screenHeight int

	gameScreenTextureID imgui.TextureID
	gameScreen          *ebiten.Image

	filter ebiten.Filter
}

// Update - Update UI and game in tandem. Handle inputs
func (g *GameProxy) Update() error {
	if g.game == nil {
		panic("No game to update!")
	}

	io := imgui.CurrentIO()

	// Sync keyboard
	currentAdapter.inputChars = sendInput(imgui.CurrentIO(), currentAdapter.inputChars)

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
	// screen will be the destination where the ui is drawn.
	// g.gameScreen is where the game is drawn.
	if g.gameScreen == nil && g.renderable() {
		g.gameScreen = ebiten.NewImage(int(g.width), int(g.height))
		Cache.SetTexture(g.gameScreenTextureID, g.gameScreen)
	}

	// Check that old frame matches new size. If not, delete old texture and create a new one.
	if g.gameScreen.Bounds().Size().X != int(g.width) ||
		g.gameScreen.Bounds().Size().Y != int(g.height) && g.renderable() {
		Cache.RemoveTexture(g.gameScreenTextureID)

		g.gameScreen = ebiten.NewImage(int(g.width), int(g.height))
		Cache.SetTexture(g.gameScreenTextureID, g.gameScreen)
	}

	if g.renderable() {
		g.game.Draw(g.gameScreen)
		ebitenutil.DebugPrint(g.gameScreen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	}

	imgui.Render()

	if currentAdapter.ClipMask {
		if currentAdapter.lmask == nil {
			w, h := screen.Size()
			currentAdapter.lmask = ebiten.NewImage(w, h)
		} else {
			w1, h1 := screen.Size()
			w2, h2 := currentAdapter.lmask.Size()
			if w1 != w2 || h1 != h2 {
				currentAdapter.lmask.Dispose()
				currentAdapter.lmask = ebiten.NewImage(w1, h1)
			}
		}
		RenderMasked(screen, currentAdapter.lmask, imgui.CurrentDrawData(), Cache, g.filter)
	} else {
		Render(screen, imgui.CurrentDrawData(), Cache, g.filter)
	}

}

func (g *GameProxy) Layout(outsideWidth, outsideHeight int) (int, int) {
	width := float64(outsideWidth) * ebiten.DeviceScaleFactor()
	height := float64(outsideHeight) * ebiten.DeviceScaleFactor()

	io := imgui.CurrentIO()
	io.SetDisplaySize(imgui.Vec2{X: float32(width), Y: float32(height)})

	return int(width), int(height)
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
