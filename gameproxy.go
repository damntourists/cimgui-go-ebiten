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

func (g *GameProxy) Update() error {
	if g.game == nil {
		panic("No game to update!")
	}

	io := imgui.CurrentIO()
	io.SetDeltaTime(1.0 / 60.0)

	xoff, yoff := ebiten.Wheel()

	io.AddMouseWheelDelta(float32(xoff), float32(yoff))

	currentAdapter.inputChars = sendInput(imgui.CurrentIO(), currentAdapter.inputChars)

	imgui.NewFrame()

	cx, cy := ebiten.CursorPosition()
	io.SetMousePos(imgui.Vec2{X: float32(cx), Y: float32(cy)})
	io.SetMouseButtonDown(0, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
	io.SetMouseButtonDown(1, ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
	io.SetMouseButtonDown(2, ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle))

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

	imgui.EndFrame()

	err := g.game.Update()
	return err
}

func (g *GameProxy) Draw(screen *ebiten.Image) {

	if g.gameScreen == nil {
		g.gameScreen = ebiten.NewImage(g.screenWidth, g.screenHeight)
		Cache.SetTexture(g.gameScreenTextureID, g.gameScreen)
	}

	//currentTexture := Cache.GetTexture(g.gameScreenTextureID)

	// Check that old frame matches new size. If not, delete old texture and create a new one.
	if g.gameScreen.Bounds().Size().X != g.screenWidth ||
		g.gameScreen.Bounds().Size().Y != g.screenHeight {
		Cache.RemoveTexture(g.gameScreenTextureID)

		g.gameScreen = ebiten.NewImage(g.screenWidth, g.screenHeight)
		Cache.SetTexture(g.gameScreenTextureID, g.gameScreen)
	}

	g.game.Draw(g.gameScreen)

	ebitenutil.DebugPrint(g.gameScreen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))

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

func (g *GameProxy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
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
	g.screenHeight = int(v.Y)
	g.screenWidth = int(v.X)
}
