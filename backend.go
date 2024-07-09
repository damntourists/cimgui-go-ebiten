package ebitenbackend

import "C"

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"runtime"
	"unsafe"
)

var Cache TextureCache

type EbitenWindowFlags int

const (
	EbitenWindowFlagsCursorMode EbitenWindowFlags = iota
	EbitenWindowFlagsCursorShape
	EbitenWindowFlagsResizingMode
	EbitenWindowFlagsFPSMode
	EbitenWindowFlagsDecorated
	EbitenWindowFlagsFloating
	EbitenWindowFlagsMaximized
	EbitenWindowFlagsMinimized
	EbitenWindowFlagsClosingHandled
	EbitenWindowFlagsMousePassthrough
)

type (
	voidCallbackFunc            func()
	DropCallback                func([]string)
	KeyCallback                 func(key, scanCode, action, mods int)
	SizeChangeCallback          func(w, h int)
	WindowCloseCallback[B ~int] func(b imgui.Backend[B])
)

var _ imgui.Backend[EbitenWindowFlags] = &EbitenBackend{}

type (
	EbitenBackend struct {
		ctx *imgui.Context

		inputChars []rune

		afterCreateContext   voidCallbackFunc
		beforeRender         voidCallbackFunc
		afterRender          voidCallbackFunc
		beforeDestroyContext voidCallbackFunc
		dropCB               DropCallback
		closeCB              imgui.WindowCloseCallback[EbitenWindowFlags]
		keyCb                KeyCallback
		loop                 voidCallbackFunc
		sizeCb               SizeChangeCallback

		// TODO: Do we need to have a handle on the window? Does ebiten let us?
		//window               uintptr

		width, height             float32
		screenWidth, screenHeight int
		bgColor                   imgui.Vec4

		clipMask bool
		lmask    *ebiten.Image

		fontAtlas *imgui.FontAtlas

		game *GameProxy
	}
)

func NewEbitenBackend() *EbitenBackend {
	Cache = NewCache()
	b := &EbitenBackend{
		//game: &GameProxy{
		//	game:       g,
		//	filter:     ebiten.FilterNearest,
		//	clipRegion: imgui.Vec2{X: 1, Y: 1},
		//},
	}
	//b.game.backend = b

	runtime.SetFinalizer(b, (*EbitenBackend).onfinalize)

	return b
}

func (b *EbitenBackend) SetGame(game ebiten.Game) {
	b.game = &GameProxy{
		backend:    b,
		game:       &game,
		filter:     ebiten.FilterNearest,
		clipRegion: imgui.Vec2{X: 1, Y: 1},
	} //game
}

func (b *EbitenBackend) Game() *ebiten.Game {
	return b.game.Game()
}

func (b *EbitenBackend) SetGameRenderDestination(dest *ebiten.Image) {
	// Cache gamescreen texture
	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, dest)
	b.game.gameScreenTextureID = tid
	b.game.gameScreen = dest
	b.SetGameScreenSize(imgui.Vec2{
		X: float32(dest.Bounds().Size().X),
		Y: float32(dest.Bounds().Size().Y),
	})
}

func (b *EbitenBackend) SetGameScreenSize(size imgui.Vec2) {
	if b.game.gameScreen == nil {
		dest := ebiten.NewImage(int(size.X), int(size.Y))
		b.SetGameRenderDestination(dest)
	}
	b.game.SetGameScreenSize(size)
}

func (b *EbitenBackend) SetAfterCreateContextHook(hook func()) {
	b.afterCreateContext = hook
}

func (b *EbitenBackend) afterCreateHook() func() {
	return b.afterCreateContext
}

func (b *EbitenBackend) SetBeforeDestroyContextHook(hook func()) {
	b.beforeDestroyContext = hook
}

func (b *EbitenBackend) beforeDestroyHook() func() {
	return b.beforeDestroyContext
}

func (b *EbitenBackend) SetBeforeRenderHook(hook func()) {
	b.beforeRender = hook
}

func (b *EbitenBackend) beforeRenderHook() func() {
	return b.beforeRender
}

func (b *EbitenBackend) SetAfterRenderHook(hook func()) {
	b.afterRender = hook
}

func (b *EbitenBackend) afterRenderHook() func() {
	return b.afterRender
}
func (b *EbitenBackend) loopFunc() func() {
	return b.loop
}
func (b *EbitenBackend) dropCallback() DropCallback {
	return b.dropCB
}
func (b *EbitenBackend) closeCallback() imgui.WindowCloseCallback[EbitenWindowFlags] {
	return b.closeCB
}
func (b *EbitenBackend) DisplaySize() (width, height int32) {
	//TODO implement me
	panic("implement me")
}
func (b *EbitenBackend) SetCursorPos(x, y float64) {
	//TODO implement me
	panic("implement me")
}
func (b *EbitenBackend) CreateWindow(title string, width, height int) {
	b.ctx = imgui.CreateContext()

	if b.afterCreateContext != nil {
		b.afterCreateContext()
	}

	io := imgui.CurrentIO()
	io.SetIniFilename("")

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(int(float64(width)), int(float64(height)))
	io.SetDisplaySize(
		imgui.Vec2{
			X: float32(float64(width)),
			Y: float32(float64(height)),
		},
	)
}
func (b *EbitenBackend) CreateTexture(pixels unsafe.Pointer, width, height int) imgui.TextureID {
	eimg := ebiten.NewImage(width, height)
	eimg.WritePixels(premultiplyPixels(pixels, width, height))

	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, eimg)
	return tid
}

func (b *EbitenBackend) CreateTextureRgba(img *image.RGBA, width, height int) imgui.TextureID {
	pix := img.Pix
	return b.CreateTexture(unsafe.Pointer(&pix), width, height)
}

func (b *EbitenBackend) DeleteTexture(id imgui.TextureID) {
	Cache.RemoveTexture(id)
}

func (b *EbitenBackend) SetSwapInterval(interval EbitenWindowFlags) error {
	ebiten.SetVsyncEnabled(interval > 0)
	return nil
}

func (b *EbitenBackend) SetInputMode(mode EbitenWindowFlags, value EbitenWindowFlags) {
	//TODO implement me
	panic("Not Implemented.")
}

// func (b *GLFWBackend) SetCloseCallback(cbfun WindowCloseCallback[GLFWWindowFlags]) {
func (b *EbitenBackend) SetCloseCallback(cb imgui.WindowCloseCallback[EbitenWindowFlags]) {
	b.closeCB = cb
}

func (b *EbitenBackend) SetBgColor(color imgui.Vec4) {
	b.bgColor = color
}

func (b *EbitenBackend) SetWindowFlags(flag EbitenWindowFlags, value int) {
	switch flag {
	case EbitenWindowFlagsCursorMode:
		ebiten.SetCursorMode(ebiten.CursorModeType(value))
	case EbitenWindowFlagsCursorShape:
		ebiten.SetCursorShape(ebiten.CursorShapeType(value))
	case EbitenWindowFlagsFPSMode:
		ebiten.SetVsyncEnabled(value == 0)
	case EbitenWindowFlagsResizingMode:
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeType(value))
	case EbitenWindowFlagsDecorated:
		ebiten.SetWindowDecorated(value != 0)
	case EbitenWindowFlagsFloating:
		ebiten.SetWindowFloating(value != 0)
	case EbitenWindowFlagsMaximized:
		if value != 0 {
			ebiten.MaximizeWindow()
		} else {
			ebiten.RestoreWindow()
		}
	case EbitenWindowFlagsMinimized:
		if value != 0 {
			ebiten.MinimizeWindow()
		} else {
			ebiten.RestoreWindow()
		}
	case EbitenWindowFlagsClosingHandled:
		ebiten.SetWindowClosingHandled(value != 0)
	case EbitenWindowFlagsMousePassthrough:
		ebiten.SetWindowMousePassthrough(value != 0)
	default:
		panic("Invalid flag for SetWindowFlags.")
	}
}

func (b *EbitenBackend) SetKeyCallback(callback imgui.KeyCallback) {
	b.keyCb = KeyCallback(callback)
}

func (b *EbitenBackend) SetSizeChangeCallback(callback imgui.SizeChangeCallback) {
	b.sizeCb = SizeChangeCallback(callback)
}

func (b *EbitenBackend) SetDropCallback(callback imgui.DropCallback) {
	b.dropCB = DropCallback(callback)
}

func (b *EbitenBackend) onfinalize() {
	runtime.SetFinalizer(b, nil)
	if b.beforeDestroyContext != nil {
		b.beforeDestroyContext()
	}
	b.ctx.Destroy()
}

func (b *EbitenBackend) Run(f func()) {
	f()
}

func (b *EbitenBackend) Refresh() {
	println("backend refreshing!")
}

func (b *EbitenBackend) SetWindowPos(x, y int) {
	ebiten.SetWindowPosition(x, y)
}

func (b *EbitenBackend) GetWindowPos() (x, y int32) {
	xx, yy := ebiten.WindowPosition()
	return int32(xx), int32(yy)
}

func (b *EbitenBackend) SetWindowSize(width, height int) {
	ebiten.SetWindowSize(width, height)
}

func (b *EbitenBackend) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	ebiten.SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (b *EbitenBackend) SetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func (b *EbitenBackend) SetShouldClose(shouldClose bool) {
	ebiten.SetWindowClosingHandled(shouldClose)
}

func (b *EbitenBackend) ContentScale() (xScale, yScale float32) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	return float32(scale), float32(scale)
}

func (b *EbitenBackend) SetTargetFPS(fps uint) {
	ebiten.SetTPS(int(fps))
}

func (b *EbitenBackend) SetIcons(icons ...image.Image) {
	ebiten.SetWindowIcon(icons)
}
