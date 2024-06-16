package ebitenbackend

import "C"

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"runtime"
	"slices"
	"unsafe"
)

var Cache TextureCache

type WindowFlags int
type WindowCloseCallback[B ~int] func(b imgui.Backend[B])

const (
	WindowFlagsNone = WindowFlags(iota)

	WindowResizeingMode
	WindowResizingModeDisabled
	WindowResizingModeOnlyFullscreenEnabled
	WindowResizingModeEnabled

	//WindowFlagsResizeMode
	//WindowFlagsResizable
	WindowFlagsMaximized
	WindowFlagsMinimized
	WindowFlagsDecorated
	WindowFlagsFloating
	WindowFlagsMousePassthrough

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

var _ imgui.Backend[WindowFlags] = &Bridge{}

type Set [256]bool

func (s *Set) Contains(b byte) bool {
	return s[b]
}

func (s *Set) Add(b byte) {
	s[b] = true
}

func (s *Set) Remove(b byte) {
	s[b] = false
}

type (
	Bridge struct {
		dropCBFn        imgui.DropCallback
		closeCBFn       imgui.WindowCloseCallback[WindowFlags]
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

func (b *Bridge) SetSwapInterval(interval WindowFlags) error {
	//TODO - I was unable to find any reference to this in ebiten.
	panic("Not Implemented.")
}

func (b *Bridge) SetCursorPos(x, y float64) {
	//TODO As of this moment, I am unable to find a method in ebiten for setting cursor
	//	position. :/
	panic("Not Implemented.")
}

func (b *Bridge) SetInputMode(mode WindowFlags, value WindowFlags) {
	//TODO implement me
	panic("Not Implemented.")
}

func (b *Bridge) SetCloseCallback(cb imgui.WindowCloseCallback[WindowFlags]) {
	b.closeCBFn = cb
}

func (b *Bridge) SetBgColor(color imgui.Vec4) {
	b.bgColor = color
}

func (b *Bridge) SetWindowFlags(flag WindowFlags, value WindowFlags) {
	//TODO implement me
	switch flag {
	case WindowResizeingMode:
		f := []WindowFlags{
			WindowResizingModeEnabled,
			WindowResizingModeDisabled,
			WindowResizingModeOnlyFullscreenEnabled,
		}
		if slices.Contains(f, value); !ok {
			ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		}

	case WindowFlagsMaximized:
		fallthrough
	case WindowFlagsMinimized:
		fallthrough
	case WindowFlagsDecorated:
		fallthrough
	case WindowFlagsFloating:
		fallthrough
	case WindowFlagsMousePassthrough:
		fallthrough
	default: // WindowFlagsNone
		//
	}

}

func (b *Bridge) SetAfterCreateContextHook(_ func()) {
	// noop
}

func (b *Bridge) SetBeforeDestroyContextHook(_ func()) {
	// noop
}

func (b *Bridge) SetBeforeRenderHook(_ func()) {
	// noop
}

func (b *Bridge) SetAfterRenderHook(_ func()) {
	// noop
}

func (b *Bridge) SetKeyCallback(callback imgui.KeyCallback) {
	b.keyCBFn = callback
}

func (b *Bridge) SetSizeChangeCallback(callback imgui.SizeChangeCallback) {
	b.sizeChangedCbFn = callback
}

func (b *Bridge) SetDropCallback(callback imgui.DropCallback) {
	b.dropCBFn = callback
}

func (b *Bridge) onfinalize() {
	runtime.SetFinalizer(b, nil)
	b.ctx.Destroy()
}

func (b *Bridge) CreateWindow(title string, width, height int) {
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

func (b *Bridge) Run(f func()) {
	f()
}

func (b *Bridge) Refresh() {
	println("backend bridge refreshing!")
}

func (b *Bridge) SetWindowPos(x, y int) {
	ebiten.SetWindowPosition(x, y)
}

func (b *Bridge) GetWindowPos() (x, y int32) {
	xx, yy := ebiten.WindowPosition()

	return int32(xx), int32(yy)
}

func (b *Bridge) SetWindowSize(width, height int) {
	ebiten.SetWindowSize(width, height)
}

func (b *Bridge) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	ebiten.SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (b *Bridge) SetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func (b *Bridge) DisplaySize() (width int32, height int32) {
	return
}

func (b *Bridge) SetShouldClose(shouldClose bool) {
	ebiten.SetWindowClosingHandled(shouldClose)
}

func (b *Bridge) ContentScale() (xScale, yScale float32) {
	scale := ebiten.DeviceScaleFactor()
	return float32(scale), float32(scale)
}

func (b *Bridge) SetTargetFPS(fps uint) {
	ebiten.SetTPS(int(fps))
}

func (b *Bridge) SetIcons(icons ...image.Image) {
	ebiten.SetWindowIcon(icons)
}

func (b *Bridge) CreateTexture(pixels unsafe.Pointer, width, height int) imgui.TextureID {
	eimg := ebiten.NewImage(width, height)
	eimg.WritePixels(premultiplyPixels(pixels, width, height))

	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, eimg)
	return tid
}

func (b *Bridge) CreateTextureRgba(img *image.RGBA, width, height int) imgui.TextureID {
	pix := img.Pix
	return b.CreateTexture(unsafe.Pointer(&pix), width, height)
}

func (b *Bridge) DeleteTexture(id imgui.TextureID) {
	Cache.RemoveTexture(id)
}
