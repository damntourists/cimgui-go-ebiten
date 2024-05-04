package ebitenbackend

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"runtime"
)

var CurrentAdapter *EbitenAdapter

var keys = map[imgui.Key]int{
	imgui.KeyTab:            int(ebiten.KeyTab),
	imgui.KeyCapsLock:       int(ebiten.KeyCapsLock),
	imgui.KeyComma:          int(ebiten.KeyComma),
	imgui.KeyMenu:           int(ebiten.KeyContextMenu),
	imgui.KeyLeftArrow:      int(ebiten.KeyArrowLeft),
	imgui.KeyRightArrow:     int(ebiten.KeyArrowRight),
	imgui.KeyUpArrow:        int(ebiten.KeyArrowUp),
	imgui.KeyDownArrow:      int(ebiten.KeyArrowDown),
	imgui.KeyBackslash:      int(ebiten.KeyBackslash),
	imgui.KeyBackspace:      int(ebiten.KeyBackspace),
	imgui.KeyLeftCtrl:       int(ebiten.KeyControlLeft),
	imgui.KeyRightCtrl:      int(ebiten.KeyControlRight),
	imgui.KeyLeftAlt:        int(ebiten.KeyAltLeft),
	imgui.KeyRightAlt:       int(ebiten.KeyAltRight),
	imgui.KeyLeftShift:      int(ebiten.KeyShiftLeft),
	imgui.KeyRightShift:     int(ebiten.KeyShiftRight),
	imgui.KeyLeftSuper:      int(ebiten.KeyMetaLeft),
	imgui.KeyRightSuper:     int(ebiten.KeyMetaRight),
	imgui.KeyLeftBracket:    int(ebiten.KeyBracketLeft),
	imgui.KeyRightBracket:   int(ebiten.KeyBracketRight),
	imgui.KeyPageUp:         int(ebiten.KeyPageUp),
	imgui.KeyPageDown:       int(ebiten.KeyPageDown),
	imgui.KeyEnd:            int(ebiten.KeyEnd),
	imgui.KeyHome:           int(ebiten.KeyHome),
	imgui.KeyInsert:         int(ebiten.KeyInsert),
	imgui.KeyDelete:         int(ebiten.KeyDelete),
	imgui.KeySpace:          int(ebiten.KeySpace),
	imgui.KeyEnter:          int(ebiten.KeyEnter),
	imgui.KeyEscape:         int(ebiten.KeyEscape),
	imgui.KeyEqual:          int(ebiten.KeyEqual),
	imgui.Key0:              int(ebiten.KeyDigit0),
	imgui.Key1:              int(ebiten.KeyDigit1),
	imgui.Key2:              int(ebiten.KeyDigit2),
	imgui.Key3:              int(ebiten.KeyDigit3),
	imgui.Key4:              int(ebiten.KeyDigit4),
	imgui.Key5:              int(ebiten.KeyDigit5),
	imgui.Key6:              int(ebiten.KeyDigit6),
	imgui.Key7:              int(ebiten.KeyDigit7),
	imgui.Key8:              int(ebiten.KeyDigit8),
	imgui.Key9:              int(ebiten.KeyDigit9),
	imgui.KeyA:              int(ebiten.KeyA),
	imgui.KeyB:              int(ebiten.KeyB),
	imgui.KeyC:              int(ebiten.KeyC),
	imgui.KeyD:              int(ebiten.KeyD),
	imgui.KeyE:              int(ebiten.KeyE),
	imgui.KeyF:              int(ebiten.KeyF),
	imgui.KeyG:              int(ebiten.KeyG),
	imgui.KeyH:              int(ebiten.KeyH),
	imgui.KeyI:              int(ebiten.KeyI),
	imgui.KeyJ:              int(ebiten.KeyJ),
	imgui.KeyK:              int(ebiten.KeyK),
	imgui.KeyL:              int(ebiten.KeyL),
	imgui.KeyM:              int(ebiten.KeyM),
	imgui.KeyN:              int(ebiten.KeyN),
	imgui.KeyO:              int(ebiten.KeyO),
	imgui.KeyP:              int(ebiten.KeyP),
	imgui.KeyQ:              int(ebiten.KeyQ),
	imgui.KeyR:              int(ebiten.KeyR),
	imgui.KeyS:              int(ebiten.KeyS),
	imgui.KeyT:              int(ebiten.KeyT),
	imgui.KeyU:              int(ebiten.KeyU),
	imgui.KeyV:              int(ebiten.KeyV),
	imgui.KeyW:              int(ebiten.KeyW),
	imgui.KeyX:              int(ebiten.KeyX),
	imgui.KeyY:              int(ebiten.KeyY),
	imgui.KeyZ:              int(ebiten.KeyZ),
	imgui.KeyKeypad0:        int(ebiten.KeyNumpad0),
	imgui.KeyKeypad1:        int(ebiten.KeyNumpad1),
	imgui.KeyKeypad2:        int(ebiten.KeyNumpad2),
	imgui.KeyKeypad3:        int(ebiten.KeyNumpad3),
	imgui.KeyKeypad4:        int(ebiten.KeyNumpad4),
	imgui.KeyKeypad5:        int(ebiten.KeyNumpad5),
	imgui.KeyKeypad6:        int(ebiten.KeyNumpad6),
	imgui.KeyKeypad7:        int(ebiten.KeyNumpad7),
	imgui.KeyKeypad8:        int(ebiten.KeyNumpad8),
	imgui.KeyKeypad9:        int(ebiten.KeyNumpad9),
	imgui.KeyF1:             int(ebiten.KeyF1),
	imgui.KeyF2:             int(ebiten.KeyF2),
	imgui.KeyF3:             int(ebiten.KeyF3),
	imgui.KeyF4:             int(ebiten.KeyF4),
	imgui.KeyF5:             int(ebiten.KeyF5),
	imgui.KeyF6:             int(ebiten.KeyF6),
	imgui.KeyF7:             int(ebiten.KeyF7),
	imgui.KeyF8:             int(ebiten.KeyF8),
	imgui.KeyF9:             int(ebiten.KeyF9),
	imgui.KeyF10:            int(ebiten.KeyF10),
	imgui.KeyF11:            int(ebiten.KeyF11),
	imgui.KeyF12:            int(ebiten.KeyF12),
	imgui.KeyApostrophe:     int(ebiten.KeyApostrophe),
	imgui.KeyMinus:          int(ebiten.KeyMinus),
	imgui.KeyPeriod:         int(ebiten.KeyPeriod),
	imgui.KeySlash:          int(ebiten.KeySlash),
	imgui.KeySemicolon:      int(ebiten.KeySemicolon),
	imgui.KeyGraveAccent:    int(ebiten.KeyGraveAccent),
	imgui.KeyScrollLock:     int(ebiten.KeyScrollLock),
	imgui.KeyNumLock:        int(ebiten.KeyNumLock),
	imgui.KeyPrintScreen:    int(ebiten.KeyPrintScreen),
	imgui.KeyPause:          int(ebiten.KeyPause),
	imgui.KeyKeypadDecimal:  int(ebiten.KeyKPDecimal),
	imgui.KeyKeypadDivide:   int(ebiten.KeyKPDivide),
	imgui.KeyKeypadMultiply: int(ebiten.KeyKPMultiply),
	imgui.KeyKeypadSubtract: int(ebiten.KeyKPSubtract),
	imgui.KeyKeypadAdd:      int(ebiten.KeyKPAdd),
	imgui.KeyKeypadEnter:    int(ebiten.KeyKPEnter),
	imgui.KeyKeypadEqual:    int(ebiten.KeyKPEqual),
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

	CurrentAdapter = &a

	return &a
}

func (a *EbitenAdapter) finalize() {
	runtime.SetFinalizer(a, nil)
}

func (a *EbitenAdapter) SetGame(g ebiten.Game) {
	// Create game wrapper
	a.game = &GameProxy{
		game:   g,
		filter: ebiten.FilterNearest,

		// Init at 1px so ebiten doesn't panic.
		width:  1,
		height: 1,

		screenHeight: 1,
		screenWidth:  1,

		clipRegion: imgui.Vec2{X: 1, Y: 1},
	}
}

func (a *EbitenAdapter) SetGameRenderDestination(dest *ebiten.Image) {
	// Cache gamescreen texture
	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, dest)
	a.game.gameScreenTextureID = tid
	a.game.gameScreen = dest
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
