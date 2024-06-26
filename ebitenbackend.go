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
		dropCBFn        imgui.DropCallback
		closeCBFn       imgui.WindowCloseCallback[EbitenWindowFlags]
		keyCBFn         imgui.KeyCallback
		sizeChangedCbFn imgui.SizeChangeCallback

		fontAtlas *imgui.FontAtlas
		io        *imgui.IO
		ctx       *imgui.Context

		lmask *ebiten.Image

		width, height             float32
		screenWidth, screenHeight int
		bgColor                   imgui.Vec4

		ClipMask bool
		Game     ebiten.Game
	}
)

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

func (b *BackendBridge) SetAfterCreateContextHook(_ func()) {
	// noop
}

func (b *BackendBridge) SetBeforeDestroyContextHook(_ func()) {
	// noop
}

func (b *BackendBridge) SetBeforeRenderHook(_ func()) {
	// noop
}

func (b *BackendBridge) SetAfterRenderHook(_ func()) {
	// noop
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

func (b *BackendBridge) CreateWindow(title string, width, height int) {
	// actually just sets up window. Run creates the window. This is
	// to satisfy the interface.
	b.ctx = imgui.CreateContext()
	b.io = imgui.CurrentIO()
	b.io.SetIniFilename("")

	imgui.PlotCreateContext()
	imgui.ImNodesCreateContext()

	sf := ebiten.Monitor().DeviceScaleFactor()
	imgui.CurrentStyle().ScaleAllSizes(float32(sf))

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(int(float64(width)), int(float64(height)))
	b.io.SetDisplaySize(
		imgui.Vec2{
			X: float32(float64(width)),
			Y: float32(float64(height)),
		},
	)

}

func (b *BackendBridge) Run(f func()) {
	f()
}

func (b *BackendBridge) Refresh() {
	println("backend bridge refreshing!")
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

func (b *BackendBridge) DisplaySize() (width int32, height int32) {
	return
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
