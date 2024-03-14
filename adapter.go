package ebitenbackend

import (
	"fmt"
	imgui "github.com/damntourists/cimgui-go-lite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"runtime"
)

var currentAdapter *EbitenAdapter

// Adapter should proxy calls to backend.
type Adapter interface {
	SetAfterCreateContextHook(func())   //noop
	SetBeforeDestroyContextHook(func()) //noop
	SetBeforeRenderHook(func())         //noop
	SetAfterRenderHook(func())          //noop

	SetBgColor(color imgui.Vec4)
	Run(func())
	Refresh()

	SetWindowPos(x, y int)
	GetWindowPos() (x, y int32)
	SetWindowSize(width, height int)
	SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int)
	SetWindowTitle(title string)
	DisplaySize() (width, height int32)
	SetShouldClose(bool)
	ContentScale() (xScale, yScale float32)

	SetTargetFPS(fps uint)

	SetIcons(icons ...image.Image)
	CreateWindow(title string, width, height int)

	Backend() *imgui.Backend[EbitenWindowFlags]
	SetGame(ebiten.Game)
	SetUILoop(func())
	Game() ebiten.Game
	Update(float32)
	finalize()
}

type EbitenAdapter struct {
	backend imgui.Backend[EbitenWindowFlags]
	game    ebiten.Game
	loop    func()

	ClipMask bool
	lmask    *ebiten.Image
	cliptxt  string
}

type GameProxy struct {
	game    ebiten.Game
	adapter *EbitenAdapter

	width, height             float64
	screenWidth, screenHeight int

	filter ebiten.Filter
}

func (g GameProxy) Update() error {
	if g.game == nil {
		panic("No game to update!")
	}

	io := imgui.CurrentIO()
	xoff, yoff := ebiten.Wheel()

	io.SetDeltaTime(1.0 / 60.0)

	io.AddMouseWheelDelta(float32(xoff), float32(yoff))

	//imgui.CurrentStyle().ScaleAllSizes(float32(ebiten.DeviceScaleFactor()))
	imgui.NewFrame()

	err := g.game.Update()
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

	return err
}

func (g GameProxy) Draw(screen *ebiten.Image) {
	g.game.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))

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

func (g GameProxy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	width := float64(outsideWidth) * ebiten.DeviceScaleFactor()
	height := float64(outsideHeight) * ebiten.DeviceScaleFactor()

	io := imgui.CurrentIO()
	io.SetDisplaySize(imgui.Vec2{X: float32(width), Y: float32(height)})

	return int(width), int(height)
}

func (a *EbitenAdapter) SetBeforeDestroyContextHook(f func()) {
	a.backend.SetBeforeDestroyContextHook(f)
}

func (a *EbitenAdapter) SetBeforeRenderHook(f func()) {
	a.backend.SetBeforeRenderHook(f)
}

func (a *EbitenAdapter) SetAfterRenderHook(f func()) {
	a.backend.SetAfterRenderHook(f)
}

func (a *EbitenAdapter) SetBgColor(color imgui.Vec4) {
	a.backend.SetBgColor(color)
}

func (a *EbitenAdapter) Refresh() {
	a.backend.Refresh()
}

func (a *EbitenAdapter) GetWindowPos() (x, y int32) {
	return a.backend.GetWindowPos()
}

func (a *EbitenAdapter) SetWindowSize(width, height int) {
	a.backend.SetWindowSize(width, height)
}

func (a *EbitenAdapter) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	a.backend.SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (a *EbitenAdapter) SetWindowTitle(title string) {
	a.backend.SetWindowTitle(title)
}

func (a *EbitenAdapter) DisplaySize() (width, height int32) {
	return a.backend.DisplaySize()
}

func (a *EbitenAdapter) SetShouldClose(b bool) {
	a.backend.SetShouldClose(b)
}

func (a *EbitenAdapter) ContentScale() (xScale, yScale float32) {
	return a.backend.ContentScale()
}

func (a *EbitenAdapter) SetTargetFPS(fps uint) {
	a.backend.SetTargetFPS(fps)
}

func (a *EbitenAdapter) SetIcons(icons ...image.Image) {
	a.backend.SetIcons(icons...)
}

func (a *EbitenAdapter) Backend() *imgui.Backend[EbitenWindowFlags] {
	return &a.backend
}

func NewEbitenAdapter() *EbitenAdapter {
	b := &BackendBridge{
		ctx: nil,
	}

	Cache = NewCache()

	b.ctx = imgui.CreateContext()
	imgui.ImNodesCreateContext()

	//fonts := imgui.CurrentIO().Fonts()
	//_, _, _, _ = fonts.GetTextureDataAsRGBA32()

	b.SetBeforeRenderHook(func() {
		// TODO
	})
	b.SetAfterRenderHook(func() {
		// TODO
	})

	bb := (imgui.Backend[EbitenWindowFlags])(b)
	bb.SetKeyCallback(func(key, scanCode, action, mods int) {
		println("key", key, "scanCode", scanCode, "action", action, "mods", mods)
	})
	createdBackend, _ := imgui.CreateBackend(bb)

	a := EbitenAdapter{
		backend:  createdBackend,
		ClipMask: true,
	}

	runtime.SetFinalizer(&a, (*EbitenAdapter).finalize)

	currentAdapter = &a

	return &a
}

func (a *EbitenAdapter) finalize() {
	runtime.SetFinalizer(a, nil)
}

func (a *EbitenAdapter) SetGame(g ebiten.Game) {

	a.game = GameProxy{
		game:   g,
		filter: ebiten.FilterNearest,
	}

}

func (a *EbitenAdapter) Game() ebiten.Game {
	return a.game
}

func (a *EbitenAdapter) SetUILoop(f func()) {
	a.loop = f
}

func (a *EbitenAdapter) UILoop() func() {
	return a.loop
}

func (a *EbitenAdapter) Update(delta float32) {
	io := imgui.CurrentIO()
	io.SetDeltaTime(delta)

	_ = a.game.Update()
}

func (a *EbitenAdapter) SetWindowPos(x, y int) {
	a.backend.SetWindowPos(x, y)
}

func (a *EbitenAdapter) CreateWindow(title string, width, height int) {
	a.backend.CreateWindow(title, width, height)
}

func (a *EbitenAdapter) Run(f func()) {
	a.backend.Run(f)
}

func (a *EbitenAdapter) SetAfterCreateContextHook(hook func()) {
	a.backend.SetAfterCreateContextHook(hook)
}
