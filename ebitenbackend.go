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

/*
Cannot use '&BackendBridge{}' (type *BackendBridge) as the type
	imgui.Backend[EbitenWindowFlags] Type does not implement
	'imgui.Backend[EbitenWindowFlags]' need the method:
	SetCloseCallback(WindowCloseCallback[BackendFlagsT]) have
	the method: SetCloseCallback(cbfun func(b BackendBridge))

WindowCloseCallback is defined as:
type WindowCloseCallback[BackendFlagsT ~int] func(b Backend[BackendFlagsT])

*/

type (
	BackendBridge struct {
		hookAfterCreateContext   func()
		hookBeforeDestroyContext func()
		hookLoop                 func()
		hookBeforeRender         func()

		hookAfterRender func()
		//beforeRender         func()
		beforeRender func()
		//afterRender          func()
		afterRenderFn func()
		//dropCB               imgui.DropCallback
		dropCBFn imgui.DropCallback
		//closeCB              func(b BackendBridge)
		closeCBFn func(b BackendBridge)
		//keyCb                imgui.KeyCallback
		keyCBFn imgui.KeyCallback
		//sizeChangeCallback   imgui.SizeChangeCallback
		sizeChangedCbFn imgui.SizeChangeCallback

		defaultFontTextureID imgui.TextureID
		fontAtlas            *imgui.FontAtlas
		io                   *imgui.IO
		ctx                  *imgui.Context

		game   ebiten.Game
		filter ebiten.Filter

		uiTx   *ebiten.Image
		gameTx *ebiten.Image

		lmask *ebiten.Image

		cache                     TextureCache
		windowFlags               EbitenWindowFlags
		width, height             float32
		screenWidth, screenHeight int

		bgColor imgui.Vec4

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

func NewBackend() *BackendBridge {
	b := &BackendBridge{
		// From what I understand, the default font is always 1
		defaultFontTextureID: imgui.TextureID{Data: uintptr(id1)},
		filter:               ebiten.FilterNearest,
	}

	runtime.SetFinalizer(b, (*BackendBridge[EbitenWindowFlags]).onfinalize)
	return b
}

func (b *BackendBridge) SetBgColor(color imgui.Vec4) {
	// TODO
}

func (b *BackendBridge) SetCloseCallback(cbfun BackendBridge) {
	// TODO
}

func (u *BackendBridge) SetWindowFlags(flag EbitenWindowFlags, value int) {
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

func (u *BackendBridge) SetAfterCreateContextHook(f func()) {
	u.hookAfterCreateContext = f
}

func (u *BackendBridge) SetBeforeDestroyContextHook(f func()) {
	u.hookBeforeDestroyContext = f
}

func (u *BackendBridge) SetBeforeRenderHook(f func()) {
	u.hookBeforeRender = f
}

func (u *BackendBridge) SetAfterRenderHook(f func()) {
	u.hookAfterRender = f
}

func (u *BackendBridge) SetKeyCallback(callback imgui.KeyCallback) {
	u.keyCBFn = callback
}

func (u *BackendBridge) SetSizeChangeCallback(callback imgui.SizeChangeCallback) {
	u.sizeChangedCbFn = callback
}

func (u *BackendBridge) SetDropCallback(callback imgui.DropCallback) {
	u.dropCBFn = callback
}

func (b *BackendBridge) SetGame(g ebiten.Game) *BackendBridge {
	b.game = g
	return b
}

func (u *BackendBridge) onfinalize() {
	runtime.SetFinalizer(u, nil)
	u.ctx.Destroy()
}

func (u *BackendBridge) Update() error {
	if u.hookLoop == nil {
		panic("UI Loop function not set!")
	}

	cx, cy := ebiten.CursorPosition()
	u.io.SetMousePos(imgui.Vec2{X: float32(cx), Y: float32(cy)})
	u.io.SetMouseButtonDown(0, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
	u.io.SetMouseButtonDown(1, ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
	u.io.SetMouseButtonDown(2, ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle))
	xoff, yoff := ebiten.Wheel()
	u.io.AddMouseWheelDelta(float32(xoff), float32(yoff))

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

	u.hookLoop()

	return nil
}

// Draw draws the pxlgame screen by one frame.
//
// The give argument represents a screen image. The updated content is adopted as the pxlgame screen.
//
// The frequency of Draw calls depends on the user's environment, especially the monitors refresh rate.
// For portability, you should not put your pxlgame logic in Draw in general.
func (u *BackendBridge) Draw(screen *ebiten.Image) {
	// TODO Consider different viewport modes.
	//   - UI over Game
	//       - Does this function properly with Imgui set with transparent background?
	//       - Is docking supported in this mode?
	//   - Game in UI viewport
	//       - Consider if we want to crop, fit, etc. This will likely affect mouse deltas
	u.screenWidth = screen.Bounds().Dx()
	u.screenHeight = screen.Bounds().Dy()

	u.game.Draw(screen)

	imgui.Render()

	if u.ClipMask {
		if u.lmask == nil {
			sz := screen.Bounds().Size()
			u.lmask = ebiten.NewImage(sz.X, sz.Y)
		} else {
			sz1 := screen.Bounds().Size()
			sz2 := u.lmask.Bounds().Size()
			if sz1.X != sz2.X || sz1.Y != sz2.Y {
				u.lmask.Dispose()
				u.lmask = ebiten.NewImage(sz1.X, sz1.Y)
			}
		}
		RenderMasked(screen, u.lmask, imgui.CurrentDrawData(), u.cache, u.filter)
	} else {
		Render(screen, imgui.CurrentDrawData(), u.cache, u.filter)
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
func (u *BackendBridge) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//TODO implement me
	u.width = float32(outsideWidth) * float32(ebiten.DeviceScaleFactor())
	u.height = float32(outsideHeight) * float32(ebiten.DeviceScaleFactor())

	u.io.SetDisplaySize(imgui.Vec2{X: u.width, Y: u.height})

	return int(u.width), int(u.height)
}

func (u *BackendBridge) CreateWindow(title string, width, height int) {
	// actually just sets up window. Run creates the window. This is
	// to satisfy the interface.
	u.ctx = imgui.CreateContext()
	u.io = imgui.CurrentIO()

	imgui.PlotCreateContext()
	imgui.ImNodesCreateContext()

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width*int(ebiten.DeviceScaleFactor()), height*int(ebiten.DeviceScaleFactor()))
}

func (u *BackendBridge) SetLoop(update func()) {
	u.hookLoop = update
}

func (u *BackendBridge) Game() ebiten.Game {
	return u.game
}

func (u *BackendBridge) Run(f func()) {
	f()
}

func (u *BackendBridge) Refresh() {
	// call refresh /update on ebiten pxlgame
	fmt.Println("refresh called")
}

func (u *BackendBridge) SetWindowPos(x, y int) {
	ebiten.SetWindowPosition(x, y)
}

func (u *BackendBridge) GetWindowPos() (x, y int32) {
	a, b := ebiten.WindowPosition()
	return int32(a), int32(b)
}

func (u *BackendBridge) SetWindowSize(width, height int) {
	ebiten.SetWindowSize(width, height)
}

func (u *BackendBridge) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	ebiten.SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (u *BackendBridge) SetWindowTitle(title string) {
	ebiten.SetWindowTitle(title)
}

func (u *BackendBridge) DisplaySize() (width, height int32) {
	return int32(u.width), int32(u.height)
}

func (u *BackendBridge) SetShouldClose(b bool) {
	ebiten.SetWindowClosingHandled(b)
}

func (u *BackendBridge) ContentScale() (xScale, yScale float32) {
	scale := ebiten.DeviceScaleFactor()
	return float32(scale), float32(scale)
}

func (u *BackendBridge) SetTargetFPS(fps uint) {
	ebiten.SetTPS(int(fps))
}

func (u *BackendBridge) SetIcons(icons ...image.Image) {
	ebiten.SetWindowIcon(icons)
}

func (u *BackendBridge) CreateTextureAs(texture imgui.TextureID, pixels unsafe.Pointer, width, height int) imgui.TextureID {
	pix := PremultiplyPixels(pixels, width, height)

	eimg := ebiten.NewImage(width, height)
	eimg.WritePixels(pix)

	u.cache.SetTexture(texture, eimg)
	return texture
}

func (u *BackendBridge) CreateTexture(pixels unsafe.Pointer, width, height int) imgui.TextureID {
	pix := PremultiplyPixels(pixels, width, height)

	eimg := ebiten.NewImage(width, height)
	eimg.WritePixels(pix)

	img := eimg.SubImage(eimg.Bounds())

	imgui.NewTextureFromRgba()

	u.cache.SetTexture(t, eimg)
	return t
}

func (u *BackendBridge) CreateTextureRgba(img *image.RGBA, width, height int) imgui.TextureID {
	fmt.Println("create texture rgba was called")
	tex := imgui.NewTextureFromRgba(img)
	// TODO Rework to store this in a texture cache.
	return tex.ID()
}

func (u *BackendBridge) DeleteTexture(id imgui.TextureID) {
	fmt.Println("delete texture was called")
	// TODO Rework to store this in a texture cache.
}
