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
	EbitenWindowFlagsNone = EbitenWindowFlags(iota)
	EbitenWindowFlagsResizable
	EbitenWindowFlagsMaximized
	EbitenWindowFlagsMinimized
	EbitenWindowFlagsDecorated
	EbitenWindowFlagsFloating
	EbitenWindowFlagsMousePassthrough

	/*
		Refer to the following for more modes/settings:
			https://github.com/hajimehoshi/ebiten/blob/main/internal/ui/ui.go#L69

			Ex:
				type WindowResizingMode int
				const (
					WindowResizingModeDisabled WindowResizingMode = iota
					WindowResizingModeOnlyFullscreenEnabled
					WindowResizingModeEnabled
				)
	*/
)

var _ imgui.Backend[EbitenWindowFlags] = &BackendBridge{}

type (
	BackendBridge struct {
		hookAfterCreateContext   func()
		hookBeforeDestroyContext func()
		hookLoop                 func()
		hookBeforeRender         func()
		hookAfterRender          func()

		dropCBFn        imgui.DropCallback
		closeCBFn       imgui.WindowCloseCallback[EbitenWindowFlags]
		keyCBFn         imgui.KeyCallback
		sizeChangedCbFn imgui.SizeChangeCallback

		fontAtlas *imgui.FontAtlas
		io        *imgui.IO
		ctx       *imgui.Context

		//filter ebiten.Filter

		//uiTx   *ebiten.Image
		//gameTx *ebiten.Image

		lmask *ebiten.Image

		//cache                     TextureCache
		width, height             float32
		screenWidth, screenHeight int
		bgColor                   imgui.Vec4

		ClipMask bool
		Game     ebiten.Game
	}
)

//func NewBackend() imgui.Backend[EbitenWindowFlags] {
//}

//
//func (b *BackendBridge) SetGame(game ebiten.Game) {
//	b.game = game
//}
//
//func (b *BackendBridge) Game() *ebiten.Game {
//	return &b.game
//}

func (b *BackendBridge) SetCloseCallback(cb imgui.WindowCloseCallback[EbitenWindowFlags]) {
	b.closeCBFn = cb
}

func (b *BackendBridge) SetBgColor(color imgui.Vec4) {
	b.bgColor = color
}

func (b *BackendBridge) SetWindowFlags(flag EbitenWindowFlags, value int) {
	//TODO implement me
	switch flag {
	case EbitenWindowFlagsResizable:
		fallthrough
	case EbitenWindowFlagsMaximized:
		fallthrough
	case EbitenWindowFlagsMinimized:
		fallthrough
	case EbitenWindowFlagsDecorated:
		fallthrough
	case EbitenWindowFlagsFloating:
		fallthrough
	case EbitenWindowFlagsMousePassthrough:
		fallthrough
	default: // EbitenWindowFlagsNone
		//
	}

}

func (b *BackendBridge) SetAfterCreateContextHook(f func()) {
	b.hookAfterCreateContext = f
}

func (b *BackendBridge) SetBeforeDestroyContextHook(f func()) {
	b.hookBeforeDestroyContext = f
}

func (b *BackendBridge) SetBeforeRenderHook(f func()) {
	b.hookBeforeRender = f
}

func (b *BackendBridge) SetAfterRenderHook(f func()) {
	b.hookAfterRender = f
}

func (b *BackendBridge) SetKeyCallback(callback imgui.KeyCallback) {
	b.keyCBFn = callback
}

func (b *BackendBridge) SetSizeChangeCallback(callback imgui.SizeChangeCallback) {
	b.sizeChangedCbFn = callback
}

func (b *BackendBridge) SetDropCallback(callback imgui.DropCallback) {
	b.dropCBFn = callback
}

func (b *BackendBridge) onfinalize() {
	runtime.SetFinalizer(b, nil)
	b.ctx.Destroy()
}

func (b *BackendBridge) Update() error {

	return nil
}

// Layout accepts a native outside size in device-independent pixels and returns the pxlgame's logical screen
// size.
//
// On desktops, the outside is a window or a monitor (fullscreen mode). On browsers, the outside is a body
// element. On mobiles, the outside is the view's size.
//
// Even though the outside size and the screen size differ, the rendering scale is automatically adjusted to
// fit with the outside.
//
// Layout is called almost every frame.
//
// It is ensured that Layout is invoked before Update is called in the first frame.
//
// If Layout returns non-positive numbers, the caller can panic.
//
// You can return a fixed screen size if you don't care, or you can also return a calculated screen size
// adjusted with the given outside size.
//
// If the pxlgame implements the interface LayoutFer, Layout is never called and LayoutF is called instead.

func (b *BackendBridge) CreateWindow(title string, width, height int) {
	// actually just sets up window. Run creates the window. This is
	// to satisfy the interface.
	b.ctx = imgui.CreateContext()
	b.io = imgui.CurrentIO()

	imgui.PlotCreateContext()
	imgui.ImNodesCreateContext()

	// build fonts

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(
		width*int(ebiten.DeviceScaleFactor()),
		height*int(ebiten.DeviceScaleFactor()),
	)
	b.io.SetDisplaySize(
		imgui.Vec2{
			X: float32(width * int(ebiten.DeviceScaleFactor())),
			Y: float32(width * int(ebiten.DeviceScaleFactor())),
		},
	)

}

func (b *BackendBridge) SetLoop(update func()) {
	b.hookLoop = update
}

func (b *BackendBridge) Run(f func()) {
	f()
}

func (b *BackendBridge) Refresh() {

}

func (b *BackendBridge) SetWindowPos(x, y int) {
	ebiten.SetWindowPosition(x, y)
}

func (b *BackendBridge) GetWindowPos() (x, y int32) {
	xx, yy := ebiten.WindowPosition()
	return int32(xx), int32(yy)
}

func (b *BackendBridge) SetWindowSize(width, height int) {
	ebiten.SetWindowSize(width, height)
}

func (b *BackendBridge) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	ebiten.SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (b *BackendBridge) SetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func (b *BackendBridge) DisplaySize() (width, height int32) {
	return int32(b.width), int32(b.height)
}

func (b *BackendBridge) SetShouldClose(shouldClose bool) {
	ebiten.SetWindowClosingHandled(shouldClose)
}

func (b *BackendBridge) ContentScale() (xScale, yScale float32) {
	scale := ebiten.DeviceScaleFactor()
	return float32(scale), float32(scale)
}

func (b *BackendBridge) SetTargetFPS(fps uint) {
	ebiten.SetTPS(int(fps))
}

func (b *BackendBridge) SetIcons(icons ...image.Image) {
	ebiten.SetWindowIcon(icons)
}

func (b *BackendBridge) CreateTexture(pixels unsafe.Pointer, width, height int) imgui.TextureID {
	eimg := ebiten.NewImage(width, height)
	eimg.WritePixels(premultiplyPixels(pixels, width, height))

	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, eimg)

	return tid
}

func (b *BackendBridge) CreateTextureRgba(img *image.RGBA, width, height int) imgui.TextureID {
	pix := img.Pix
	return b.CreateTexture(unsafe.Pointer(&pix), width, height)
}

func (b *BackendBridge) DeleteTexture(id imgui.TextureID) {
	Cache.RemoveTexture(id)
}
