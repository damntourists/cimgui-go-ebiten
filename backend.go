package ebitenbackend

import "C"

import (
	"errors"
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"runtime"
	"unsafe"
)

var Cache TextureCache

type WindowFlags int
type WindowResizingMode int
type WindowCloseCallback[B ~int] func(b imgui.Backend[B])

const (
	EBTrue  = WindowFlags(1)
	EBFalse = WindowFlags(0)
)

const (
	//SwapIntervalImmediate    = SDLWindowFlags(0)
	SwapIntervalVsync = WindowFlags(1)
	//SwapIntervalAdaptiveSync = SDLWindowFlags(-1)
)

const (
	WindowFlagsNone = WindowFlags(iota)

	WindowFlagsMaximized
	WindowFlagsMinimized
	WindowFlagsDecorated
	WindowFlagsFloating
	WindowFlagsMousePassthrough
	/*
	     // Window flags
	     const (
	   	SDLWindowFlagsFullScreen        = SDLWindowFlags(C.SDL_WINDOW_FULLSCREEN)
	   	SDLWindowFlagsOpengl            = SDLWindowFlags(C.SDL_WINDOW_OPENGL)
	   	SDLWindowFlagsTransparent       = SDLWindowFlags(C.SDL_WINDOW_HIDDEN)
	   	SDLWindowFlagsVisible           = SDLWindowFlags(C.SDL_WINDOW_BORDERLESS)
	   	SDLWindowFlagsResizable         = SDLWindowFlags(C.SDL_WINDOW_RESIZABLE)
	   	SDLWindowFlagsMouseGrabbed      = SDLWindowFlags(C.SDL_WINDOW_MOUSE_GRABBED)
	   	SDLWindowFlagsInputFocus        = SDLWindowFlags(C.SDL_WINDOW_INPUT_FOCUS)
	   	SDLWindowFlagsMouseFocus        = SDLWindowFlags(C.SDL_WINDOW_MOUSE_FOCUS)
	   	SDLWindowFlagsFullscreenDesktop = SDLWindowFlags(C.SDL_WINDOW_FULLSCREEN_DESKTOP)
	   	SDLWindowFlagsWindowForeign     = SDLWindowFlags(C.SDL_WINDOW_FOREIGN)
	   	SDLWindowFlagsAllowHighDPI      = SDLWindowFlags(C.SDL_WINDOW_ALLOW_HIGHDPI)
	   	SDLWindowFlagsMouseCapture      = SDLWindowFlags(C.SDL_WINDOW_MOUSE_CAPTURE)
	   	SDLWindowFlagsAlwaysOnTop       = SDLWindowFlags(C.SDL_WINDOW_ALWAYS_ON_TOP)
	   	SDLWindowFlagsSkipTaskbar       = SDLWindowFlags(C.SDL_WINDOW_SKIP_TASKBAR)
	   	SDLWindowFlagsUtility           = SDLWindowFlags(C.SDL_WINDOW_UTILITY)
	   	SDLWindowFlagsTooltip           = SDLWindowFlags(C.SDL_WINDOW_TOOLTIP)
	   	SDLWindowFlagsPopupMenu         = SDLWindowFlags(C.SDL_WINDOW_POPUP_MENU)
	   	SDLWindowFlagsKeyboardGrabbed   = SDLWindowFlags(C.SDL_WINDOW_KEYBOARD_GRABBED)
	   	SDLWindowFlagsWindowVulkan      = SDLWindowFlags(C.SDL_WINDOW_VULKAN)
	   	SDLWindowFlagsWindowMetal       = SDLWindowFlags(C.SDL_WINDOW_METAL)
	     )

	      // SetWindowHint applies to next CreateWindow call
	      // so use it before CreateWindow call ;-)
	      // default applied flags: SDLWindowFlagsOpengl | SDLWindowFlagsResizable | SDLWindowFlagsAllowHighDPI
	      // set flag if value is 1, clear flag if value is 0
	      func (b *SDLBackend) SetWindowFlags(flag SDLWindowFlags, value int) {
	   	C.igSDLWindowHint(C.SDL_WindowFlags(flag), C.int(value))
	      }

	*/
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
	var err error
	switch interval {
	case SwapIntervalVsync:
		ebiten.SetVsyncEnabled(true)
	default:
		err = errors.New("ebiten: invalid swapping interval")
	}
	return err
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

func (b *Bridge) SetWindowResizingMode(mode ebiten.WindowResizingModeType) {
	ebiten.SetWindowResizingMode(mode)
}

func (b *Bridge) SetWindowFlags(flag WindowFlags, value int) {
	//TODO implement me
	switch flag {
	case WindowFlagsMaximized:
		switch value {
		case int(EBTrue):
			ebiten.MaximizeWindow()
		case int(EBFalse):
			ebiten.RestoreWindow()
		default:
			panic("Invalid value for WindowFlagsMaximized.")
		}
	case WindowFlagsMinimized:
		switch value {
		case int(EBTrue):
			ebiten.MinimizeWindow()
		case int(EBFalse):
			ebiten.RestoreWindow()
		default:
			panic("Invalid value for WindowFlagsMinimized.")
		}
	case WindowFlagsDecorated:
		ebiten.SetWindowDecorated(value == int(EBTrue))
	case WindowFlagsFloating:
		ebiten.SetWindowFloating(value == int(EBTrue))
	case WindowFlagsMousePassthrough:
		ebiten.SetWindowMousePassthrough(value == int(EBTrue))
	default: // WindowFlagsNone
		panic("Invalid flag for SetWindowFlags.")
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
	scale := ebiten.Monitor().DeviceScaleFactor()
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
