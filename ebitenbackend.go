package ebitenbackend

import "C"

import (
	"fmt"
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"runtime"
	"unsafe"
)

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

//type WindowCloseCallback[EbitenWindowFlags ~int] func(b imgui.Backend[EbitenWindowFlags])
//type onClose WindowCloseCallback[EbitenWindowFlags]

type (
	BackendBridge struct {
		//imgui.Backend[EbitenWindowFlags]
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

		game   ebiten.Game
		filter ebiten.Filter

		uiTx   *ebiten.Image
		gameTx *ebiten.Image

		lmask *ebiten.Image

		cache                     TextureCache
		width, height             float32
		screenWidth, screenHeight int
		bgColor                   imgui.Vec4

		ClipMask bool

		/*
			To satisfy TextureManager interface we need to have these:
			type TextureManager interface {
				CreateTexture(pixels unsafe.Pointer, width, Height int) TextureID
				CreateTextureRgba(img *image.RGBA, width, height int) TextureID
				DeleteTexture(id TextureID)
			}
		*/
	}
)

/*
WindowCloseCallback is defined as:

*/

func NewBackend() imgui.Backend[EbitenWindowFlags] {
	b := &BackendBridge{
		cache:  NewCache(),
		filter: ebiten.FilterNearest,
	}

	runtime.SetFinalizer(b, (*BackendBridge).onfinalize)

	bb := (imgui.Backend[EbitenWindowFlags])(b)
	return bb
}

func (b *BackendBridge) SetCloseCallback(cb imgui.WindowCloseCallback[EbitenWindowFlags]) {}

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

func (b *BackendBridge) SetGame(g ebiten.Game) *BackendBridge {
	b.game = g
	return b
}

func (b *BackendBridge) onfinalize() {
	runtime.SetFinalizer(b, nil)
	b.ctx.Destroy()
}

func (b *BackendBridge) Update() error {
	if b.hookLoop == nil {
		panic("UI Loop function not set!")
	}

	cx, cy := ebiten.CursorPosition()
	b.io.SetMousePos(imgui.Vec2{X: float32(cx), Y: float32(cy)})
	b.io.SetMouseButtonDown(0, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
	b.io.SetMouseButtonDown(1, ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
	b.io.SetMouseButtonDown(2, ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle))
	xoff, yoff := ebiten.Wheel()
	b.io.AddMouseWheelDelta(float32(xoff), float32(yoff))

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

	b.hookLoop()

	return nil
}

// The frequency of Draw calls depends on the user's environment, especially the monitors
// refresh rate. For portability, you should not put your pxlgame logic in Draw in general.
func (b *BackendBridge) Draw(screen *ebiten.Image) {
	// TODO Consider different viewport modes.
	//   - UI over Game
	//       - Does this function properly with Imgui set with transparent background?
	//       - Is docking supported in this mode?
	//   - Game in UI viewport
	//       - Consider if we want to crop, fit, etc. This will likely affect mouse deltas
	b.screenWidth = screen.Bounds().Dx()
	b.screenHeight = screen.Bounds().Dy()

	b.game.Draw(screen)

	imgui.Render()

	if b.ClipMask {
		if b.lmask == nil {
			sz := screen.Bounds().Size()
			b.lmask = ebiten.NewImage(sz.X, sz.Y)
		} else {
			sz1 := screen.Bounds().Size()
			sz2 := b.lmask.Bounds().Size()
			if sz1.X != sz2.X || sz1.Y != sz2.Y {
				b.lmask.Dispose()
				b.lmask = ebiten.NewImage(sz1.X, sz1.Y)
			}
		}
		RenderMasked(screen, b.lmask, imgui.CurrentDrawData(), b.cache, b.filter)
	} else {
		Render(screen, imgui.CurrentDrawData(), b.cache, b.filter)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
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
func (b *BackendBridge) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//TODO implement me
	b.width = float32(outsideWidth) * float32(ebiten.DeviceScaleFactor())
	b.height = float32(outsideHeight) * float32(ebiten.DeviceScaleFactor())

	b.io.SetDisplaySize(imgui.Vec2{X: b.width, Y: b.height})

	return int(b.width), int(b.height)
}

func (b *BackendBridge) CreateWindow(title string, width, height int) {
	// actually just sets up window. Run creates the window. This is
	// to satisfy the interface.
	b.ctx = imgui.CreateContext()
	b.io = imgui.CurrentIO()

	imgui.PlotCreateContext()
	imgui.ImNodesCreateContext()

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width*int(ebiten.DeviceScaleFactor()), height*int(ebiten.DeviceScaleFactor()))
}

func (b *BackendBridge) SetLoop(update func()) {
	b.hookLoop = update
}

func (b *BackendBridge) Game() ebiten.Game {
	return b.game
}

func (b *BackendBridge) Run(f func()) {
	f()
}

func (b *BackendBridge) Refresh() {
	fmt.Println("refresh called")
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

	tid := imgui.TextureID{Data: uintptr(b.cache.NextId())}
	b.cache.SetTexture(tid, eimg)

	return tid
}

func (b *BackendBridge) CreateTextureRgba(img *image.RGBA, width, height int) imgui.TextureID {
	pix := img.Pix
	return b.CreateTexture(unsafe.Pointer(&pix), width, height)
}

func (b *BackendBridge) DeleteTexture(id imgui.TextureID) {
	b.cache.RemoveTexture(id)
}
