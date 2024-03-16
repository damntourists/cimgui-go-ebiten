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
	imgui.KeyA:          int(ebiten.KeyA),
	imgui.KeyC:          int(ebiten.KeyC),
	imgui.KeyV:          int(ebiten.KeyV),
	imgui.KeyX:          int(ebiten.KeyX),
	imgui.KeyY:          int(ebiten.KeyY),
	imgui.KeyZ:          int(ebiten.KeyZ),
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

	gameScreen := ebiten.NewImage(16, 16)

	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, gameScreen)

	a.game = &GameProxy{
		game:                g,
		filter:              ebiten.FilterNearest,
		gameScreen:          gameScreen,
		gameScreenTextureID: tid,
	}
}

func (a *EbitenAdapter) ScreenTextureID() imgui.TextureID {
	return a.ScreenTextureID()
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

func (a *EbitenAdapter) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	/*
		io := cimgui.GetIO()
		for imguiKey, nativeKey := range keys {
			// io.KeyMap(int(imguiKey), nativeKey)
		}
	*/
}

func (a *EbitenAdapter) SetGameScreenSize(avail imgui.Vec2) {
	a.Game().
}
