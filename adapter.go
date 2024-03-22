package ebitenbackend

import (
	imgui "github.com/damntourists/cimgui-go-lite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"runtime"
)

var currentAdapter *EbitenAdapter

var keys = map[imgui.Key]int{
	/*
			imgui.KeyNone :
		imgui.Key_BEGIN
		imgui.KeyTab        :
		imgui.KeyLeftArrow  :
		imgui.KeyRightArrow :
		imgui.KeyUpArrow    :
		imgui.KeyDownArrow  :
		imgui.KeyPageUp     :
		imgui.KeyPageDown   :
		imgui.KeyHome       :
		imgui.KeyEnd        :
		imgui.KeyInsert     :
		imgui.KeyDelete     :
		imgui.KeyBackspace  :
		imgui.KeySpace      :
		imgui.KeyEnter      :
		imgui.KeyEscape     :
		imgui.KeyLeftCtrl   :
		imgui.KeyLeftShift  :
		imgui.KeyLeftAlt    :
		imgui.KeyLeftSuper  :
		imgui.KeyRightCtrl  :
		imgui.KeyRightShift :
		imgui.KeyRightAlt   :
		imgui.KeyRightSuper :
		imgui.KeyMenu       :



		int(KeyAltLeft),
		int(KeyAltRight),

		int(KeyArrowDown),
		int(KeyArrowLeft),
		int(KeyArrowRight),
		int(KeyArrowUp),

		int(KeyBackquote),
		int(KeyBackslash),
		int(KeyBackspace),

		int(KeyBracketLeft),
		int(KeyBracketRight),

		int(KeyCapsLock),
		int(KeyComma),
		int(KeyContextMenu),
		int(KeyControlLeft),
		int(KeyControlRight),

	*/
	imgui.KeyTab:        int(ebiten.KeyTab),
	imgui.KeyLeftArrow:  int(ebiten.KeyLeft),
	imgui.KeyRightArrow: int(ebiten.KeyRight),
	imgui.KeyUpArrow:    int(ebiten.KeyUp),
	imgui.KeyDownArrow:  int(ebiten.KeyDown),
	imgui.KeyPageUp:     int(ebiten.KeyPageUp),
	imgui.KeyPageDown:   int(ebiten.KeyPageDown),
	imgui.KeyHome:       int(ebiten.KeyHome),
	imgui.KeyEnd:        int(ebiten.KeyEnd),
	imgui.KeyInsert:     int(ebiten.KeyInsert),
	imgui.KeyDelete:     int(ebiten.KeyDelete),
	imgui.KeyBackspace:  int(ebiten.KeyBackspace),
	imgui.KeySpace:      int(ebiten.KeySpace),
	imgui.KeyEnter:      int(ebiten.KeyEnter),
	imgui.KeyEscape:     int(ebiten.KeyEscape),

	imgui.Key0:       int(ebiten.KeyDigit0),
	imgui.Key1:       int(ebiten.KeyDigit1),
	imgui.Key2:       int(ebiten.KeyDigit2),
	imgui.Key3:       int(ebiten.KeyDigit3),
	imgui.Key4:       int(ebiten.KeyDigit4),
	imgui.Key5:       int(ebiten.KeyDigit5),
	imgui.Key6:       int(ebiten.KeyDigit6),
	imgui.Key7:       int(ebiten.KeyDigit7),
	imgui.Key8:       int(ebiten.KeyDigit8),
	imgui.Key9:       int(ebiten.KeyDigit9),
	imgui.KeyA:       int(ebiten.KeyA),
	imgui.KeyB:       int(ebiten.KeyB),
	imgui.KeyC:       int(ebiten.KeyC),
	imgui.KeyD:       int(ebiten.KeyD),
	imgui.KeyE:       int(ebiten.KeyE),
	imgui.KeyF:       int(ebiten.KeyF),
	imgui.KeyG:       int(ebiten.KeyG),
	imgui.KeyH:       int(ebiten.KeyH),
	imgui.KeyI:       int(ebiten.KeyI),
	imgui.KeyJ:       int(ebiten.KeyJ),
	imgui.KeyK:       int(ebiten.KeyK),
	imgui.KeyL:       int(ebiten.KeyL),
	imgui.KeyM:       int(ebiten.KeyM),
	imgui.KeyN:       int(ebiten.KeyN),
	imgui.KeyO:       int(ebiten.KeyO),
	imgui.KeyP:       int(ebiten.KeyP),
	imgui.KeyQ:       int(ebiten.KeyQ),
	imgui.KeyR:       int(ebiten.KeyR),
	imgui.KeyS:       int(ebiten.KeyS),
	imgui.KeyT:       int(ebiten.KeyT),
	imgui.KeyU:       int(ebiten.KeyU),
	imgui.KeyV:       int(ebiten.KeyV),
	imgui.KeyW:       int(ebiten.KeyW),
	imgui.KeyX:       int(ebiten.KeyX),
	imgui.KeyY:       int(ebiten.KeyY),
	imgui.KeyZ:       int(ebiten.KeyZ),
	imgui.KeyKeypad0: int(ebiten.KeyNumpad0),
	imgui.KeyKeypad1: int(ebiten.KeyNumpad1),
	imgui.KeyKeypad2: int(ebiten.KeyNumpad2),
	imgui.KeyKeypad3: int(ebiten.KeyNumpad3),
	imgui.KeyKeypad4: int(ebiten.KeyNumpad4),
	imgui.KeyKeypad5: int(ebiten.KeyNumpad5),
	imgui.KeyKeypad6: int(ebiten.KeyNumpad6),
	imgui.KeyKeypad7: int(ebiten.KeyNumpad7),
	imgui.KeyKeypad8: int(ebiten.KeyNumpad8),
	imgui.KeyKeypad9: int(ebiten.KeyNumpad9),
	/*

		int(KeyF1),
		int(KeyF2),
		int(KeyF3),
		int(KeyF4),
		int(KeyF5),
		int(KeyF6),
		int(KeyF7),
		int(KeyF8),
		int(KeyF9),
		int(KeyF10),
		int(KeyF11),
		int(KeyF12),

	*/
	imgui.KeyF1:  nil,
	imgui.KeyF2:  nil,
	imgui.KeyF3:  nil,
	imgui.KeyF4:  nil,
	imgui.KeyF5:  nil,
	imgui.KeyF6:  nil,
	imgui.KeyF7:  nil,
	imgui.KeyF8:  nil,
	imgui.KeyF9:  nil,
	imgui.KeyF10: nil,
	imgui.KeyF11: nil,
	imgui.KeyF12: nil,
	imgui.KeyF13: nil,
	imgui.KeyF14: nil,
	imgui.KeyF15: nil,
	imgui.KeyF16: nil,
	imgui.KeyF17: nil,
	imgui.KeyF18: nil,
	imgui.KeyF19: nil,
	imgui.KeyF20: nil,
	imgui.KeyF21: nil,
	imgui.KeyF22: nil,
	imgui.KeyF23: nil,
	imgui.KeyF24: nil,

	/*

					   EBITEN:
						int(KeyDelete),
						int(KeyEnd),
						int(KeyEnter),
						int(KeyEqual),
						int(KeyEscape),

		nil,
		nil,



				imgui.KeyApostrophe :
				imgui.KeyComma :
				imgui.KeyMinus :
				imgui.KeyPeriod :
				imgui.KeySlash :
				imgui.KeySemicolon :
				imgui.KeyEqual :
				imgui.KeyLeftBracket :
				imgui.KeyBackslash :
				imgui.KeyRightBracket :
				imgui.KeyGraveAccent    :
				imgui.KeyCapsLock       :
				imgui.KeyScrollLock     :
				imgui.KeyNumLock        :
				imgui.KeyPrintScreen    :
				imgui.KeyPause          :



						int(KeyHome),
						int(KeyInsert),
						int(KeyMetaLeft),
						int(KeyMetaRight),
						int(KeyMinus),
						int(KeyNumLock),




				imgui.KeypadDecimal  :
				imgui.KeypadDivide   :
				imgui.KeypadMultiply :
				imgui.KeypadSubtract :
				imgui.KeypadAdd      :
				imgui.KeypadEnter    :
				imgui.KeypadEqual    :


						int(KeyNumpadAdd),
						int(KeyNumpadDecimal),
						int(KeyNumpadDivide),
						int(KeyNumpadEnter),
						int(KeyNumpadEqual),
						int(KeyNumpadMultiply),
						int(KeyNumpadSubtract),
						int(KeyPageDown),
						int(KeyPageUp),
						int(KeyPause),
						int(KeyPeriod),
						int(KeyPrintScreen),
						int(KeyQuote),
						int(KeyScrollLock),
						int(KeySemicolon),
						int(KeyShiftLeft),
						int(KeyShiftRight),
						int(KeySlash),
						int(KeySpace),
						int(KeyTab),
						int(KeyAlt),
						int(KeyControl),
						int(KeyShift),
						int(KeyMeta),
						int(KeyMax

					// Keys for backward compatibility.
					// Deprecated: as of v2.1.
						int(Key0),
						int(Key1),
						int(Key2),
						int(Key3),
						int(Key4),
						int(Key5),
						int(Key6),
						int(Key7),
						int(Key8),
						int(Key9),
						int(KeyApostrophe),
						int(KeyDown),
						int(KeyGraveAccent),
						int(KeyKP0),
						int(KeyKP1),
						int(KeyKP2),
						int(KeyKP3),
						int(KeyKP4),
						int(KeyKP5),
						int(KeyKP6),
						int(KeyKP7),
						int(KeyKP8),
						int(KeyKP9),
						int(KeyKPAdd),
						int(KeyKPDecimal),
						int(KeyKPDivide),
						int(KeyKPEnter),
						int(KeyKPEqual),
						int(KeyKPMultiply),
						int(KeyKPSubtract),
						int(KeyLeft),
						int(KeyLeftBracket),
						int(KeyMenu),
						int(KeyRight),
						int(KeyRightBracket),
						int(KeyUp),





							   IMGUI
				imgui.Keyboard/mouses. Often referred as "Browser Back"
				imgui.KeyAppBack    :
				imgui.KeyAppForward :
							   		// Menu (Xbox)      + (Switch)   Start/Options (PS)
				imgui.KeyGamepadStart :
							   		// View (Xbox)      - (Switch)   Share (PS)
				imgui.KeyGamepadBack :
							   		// X (Xbox)         Y (Switch)   Square (PS)        // Tap: Toggle Menu. Hold: Windowing mode (Focus/Move/Resize windows)
				imgui.KeyGamepadFaceLeft :
							   		// B (Xbox)         A (Switch)   Circle (PS)        // Cancel / Close / Exit
				imgui.KeyGamepadFaceRight :
				imgui.Keyboard
				imgui.KeyGamepadFaceUp :
							   		// A (Xbox)         B (Switch)   Cross (PS)         // Activate / Open / Toggle / Tweak
				imgui.KeyGamepadFaceDown :
							   		// D-pad Left                                       // Move / Tweak / Resize Window (in Windowing mode)
				imgui.KeyGamepadDpadLeft :
							   		// D-pad Right                                      // Move / Tweak / Resize Window (in Windowing mode)
				imgui.KeyGamepadDpadRight :
							   		// D-pad Up                                         // Move / Tweak / Resize Window (in Windowing mode)
				imgui.KeyGamepadDpadUp :
							   		// D-pad Down                                       // Move / Tweak / Resize Window (in Windowing mode)
				imgui.KeyGamepadDpadDown :
							   		// L Bumper (Xbox)  L (Switch)   L1 (PS)            // Tweak Slower / Focus Previous (in Windowing mode)
				imgui.KeyGamepadL1 :
							   		// R Bumper (Xbox)  R (Switch)   R1 (PS)            // Tweak Faster / Focus Next (in Windowing mode)
				imgui.KeyGamepadR1 :
							   		// L Trig. (Xbox)   ZL (Switch)  L2 (PS) [Analog]
				imgui.KeyGamepadL2 :
							   		// R Trig. (Xbox)   ZR (Switch)  R2 (PS) [Analog]
				imgui.KeyGamepadR2 :
							   		// L Stick (Xbox)   L3 (Switch)  L3 (PS)
				imgui.KeyGamepadL3 :
							   		// R Stick (Xbox)   R3 (Switch)  R3 (PS)
				imgui.KeyGamepadR3 :
							   		// [Analog]                                         // Move Window (in Windowing mode)
				imgui.KeyGamepadLStickLeft :
							   		// [Analog]                                         // Move Window (in Windowing mode)
				imgui.KeyGamepadLStickRight :
							   		// [Analog]                                         // Move Window (in Windowing mode)
				imgui.KeyGamepadLStickUp :
							   		// [Analog]                                         // Move Window (in Windowing mode)
				imgui.KeyGamepadLStickDown :
							   		// [Analog]
				imgui.KeyGamepadRStickLeft :
							   		// [Analog]
				imgui.KeyGamepadRStickRight :
							   		// [Analog]
				imgui.KeyGamepadRStickUp :
							   		// [Analog]
				imgui.KeyGamepadRStickDown   :
				imgui.KeyMouseLeft           :
				imgui.KeyMouseRight          :
				imgui.KeyMouseMiddle         :
				imgui.KeyMouseX1             :
				imgui.KeyMouseX2             :
				imgui.KeyMouseWheelX         :
				imgui.KeyMouseWheelY         :
				imgui.KeyReservedForModCtrl  :
				imgui.KeyReservedForModShift :
				imgui.KeyReservedForModAlt   :
				imgui.KeyReservedForModSuper :
				imgui.KeyCOUNT               :
				imgui.KeyBEGIN :
				imgui.KeyEND   :
				imgui.KeyCOUNT :
				imgui.Keys
				imgui.KeysDataSIZE :
				imgui.KeysData_OFFSET) index.
				imgui.KeysDataOFFSET :
							   	)





								   		ModNone                ),
								   		// Ctrl
								   		ModCtrl ),
								   		// Shift
								   		ModShift ),
								   		// Option/Menu
								   		ModAlt ),
								   		// Cmd/Super/Windows
								   		ModSuper ),
								   		// Alias for Ctrl (non-macOS) _or_ Super (macOS).
								   		ModShortcut ),
								   		// 5-bits
								   		ModMask          ),

	*/

}

func sendInput(io *imgui.IO, inputChars []rune) []rune {
	// Ebiten hides the LeftAlt RightAlt implementation (inside the uiDriver()), so
	// here only the left alt is sent
	if ebiten.IsKeyPressed(ebiten.KeyAlt) {
		io.SetKeyAlt(true)
	} else {
		io.SetKeyAlt(false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		io.SetKeyShift(true)
	} else {
		io.SetKeyShift(false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		io.SetKeyCtrl(true)
	} else {
		io.SetKeyCtrl(false)
	}
	// TODO: get KeySuper somehow (GLFW: KeyLeftSuper    = Key(343), R: 347)
	inputChars = ebiten.AppendInputChars(inputChars)
	if len(inputChars) > 0 {
		io.AddInputCharactersUTF8(string(inputChars))
		inputChars = inputChars[:0]
	}
	for ik, iv := range keys {
		if inpututil.IsKeyJustPressed(ebiten.Key(iv)) {
			io.AddKeyEvent(ik, true)
		}
		if inpututil.IsKeyJustReleased(ebiten.Key(iv)) {
			io.AddKeyEvent(ik, false)
		}
	}
	return inputChars
}

// Adapter should proxy calls to backend.
type Adapter interface {
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
	Game() ebiten.Game
	//Update(float32)
	finalize()
}

type EbitenAdapter struct {
	backend imgui.Backend[EbitenWindowFlags]
	game    *GameProxy
	loop    func()

	ClipMask   bool
	lmask      *ebiten.Image
	cliptxt    string
	inputChars []rune
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

	bb := (imgui.Backend[EbitenWindowFlags])(b)
	createdBackend, _ := imgui.CreateBackend(bb)

	a := EbitenAdapter{
		backend:    createdBackend,
		ClipMask:   true,
		inputChars: make([]rune, 0, 256),
	}
	a.setKeyMapping()

	runtime.SetFinalizer(&a, (*EbitenAdapter).finalize)

	currentAdapter = &a

	return &a
}

func (a *EbitenAdapter) finalize() {
	runtime.SetFinalizer(a, nil)
}

func (a *EbitenAdapter) SetGame(g ebiten.Game) {
	// init a 1px
	gameScreen := ebiten.NewImage(1, 1)

	// Cache gamescreen texture
	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, gameScreen)

	// Create game wrapper
	a.game = &GameProxy{
		game:                g,
		filter:              ebiten.FilterNearest,
		gameScreen:          gameScreen,
		gameScreenTextureID: tid,

		// Init at 1px so ebiten doesn't panic.
		width:        1,
		height:       1,
		screenHeight: 1,
		screenWidth:  1,
	}
}

func (a *EbitenAdapter) ScreenTextureID() imgui.TextureID {
	return a.game.ScreenTextureID()
}

func (a *EbitenAdapter) Game() ebiten.Game {
	return a.game
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

func (a *EbitenAdapter) setKeyMapping() {
	// TODO
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	/*
		io := cimgui.GetIO()
		for imguiKey, nativeKey := range keys {
			// io.KeyMap(int(imguiKey), nativeKey)
		}
	*/
}

func (a *EbitenAdapter) SetGameScreenSize(avail imgui.Vec2) {
	a.game.SetGameScreenSize(avail)
}
