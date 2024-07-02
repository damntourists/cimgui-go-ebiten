package ebitenbackend


import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"runtime"
)

var CurrentAdapter *Adapter



// AdapterType should proxy calls to backend.
type AdapterType interface {
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

type Adapter struct {
	backend imgui.Backend[EbitenWindowFlags]
	game    *GameProxy
	loop    func()

	ClipMask   bool
	lmask      *ebiten.Image
	cliptxt    string
	inputChars []rune
}

func (a *Adapter) SetBeforeDestroyContextHook(f func()) {
	a.backend.SetBeforeDestroyContextHook(f)
}

func (a *Adapter) SetBeforeRenderHook(f func()) {
	a.backend.SetBeforeRenderHook(f)
}
		fontAtlas *imgui.FontAtlas

func (a *Adapter) SetAfterRenderHook(f func()) {
	a.backend.SetAfterRenderHook(f)
}

func (a *Adapter) SetBgColor(color imgui.Vec4) {
	a.backend.SetBgColor(color)
}

func (a *Adapter) Refresh() {
	a.backend.Refresh()
}

func (a *Adapter) GetWindowPos() (x, y int32) {
	return a.backend.GetWindowPos()
}

func (a *Adapter) SetWindowSize(width, height int) {
	a.backend.SetWindowSize(width, height)
}

func (a *Adapter) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	a.backend.SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight)
}

func (a *Adapter) SetWindowTitle(title string) {
	a.backend.SetWindowTitle(title)
}

func (a *Adapter) DisplaySize() (width, height int32) {
	return a.backend.DisplaySize()
}

func (a *Adapter) SetShouldClose(b bool) {
	a.backend.SetShouldClose(b)
}

func (a *Adapter) ContentScale() (xScale, yScale float32) {
	return a.backend.ContentScale()
}

func (a *Adapter) SetTargetFPS(fps uint) {
	a.backend.SetTargetFPS(fps)
}

func (a *Adapter) SetIcons(icons ...image.Image) {
	a.backend.SetIcons(icons...)
}

func (a *Adapter) Backend() *imgui.Backend[EbitenWindowFlags] {
	return &a.backend
}

func NewEbitenAdapter() *Adapter {
	b := &Bridge{
		ctx: nil,
	}

	Cache = NewCache()

	b.ctx = imgui.CreateContext()
	//imgui.ImNodesCreateContext()

	bb := (imgui.Backend[EbitenWindowFlags])(b)
	createdBackend, _ := imgui.CreateBackend(bb)

	a := Adapter{
		backend:    createdBackend,
		ClipMask:   true,
		inputChars: make([]rune, 0, 256),
	}
	a.setKeyMapping()

	runtime.SetFinalizer(&a, (*Adapter).finalize)

	CurrentAdapter = &a

	return &a
}

func (a *Adapter) finalize() {
	runtime.SetFinalizer(a, nil)
}

func (a *Adapter) SetGame(g ebiten.Game) {
	// Create game wrapper
	a.game = &GameProxy{
		game:       g,
		filter:     ebiten.FilterNearest,
		clipRegion: imgui.Vec2{X: 1, Y: 1},
	}
}

func (a *Adapter) SetGameRenderDestination(dest *ebiten.Image) {
	// Cache gamescreen texture
	tid := imgui.TextureID{Data: uintptr(Cache.NextId())}
	Cache.SetTexture(tid, dest)
	a.game.gameScreenTextureID = tid
	a.game.gameScreen = dest
	a.SetGameScreenSize(imgui.Vec2{
		X: float32(dest.Bounds().Size().X),
		Y: float32(dest.Bounds().Size().Y),
	})
}

func (a *Adapter) ScreenTextureID() imgui.TextureID {
	return a.game.ScreenTextureID()
}

func (a *Adapter) Game() *GameProxy {
	return a.game
}

func (a *Adapter) SetWindowPos(x, y int) {
	a.backend.SetWindowPos(x, y)
}

func (a *Adapter) CreateWindow(title string, width, height int) {
	a.backend.CreateWindow(title, width, height)
}

func (a *Adapter) Run(f func()) {
	a.backend.Run(f)
}

func (a *Adapter) setKeyMapping() {
	// TODO
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	/*
		io := cimgui.GetIO()
		for imguiKey, nativeKey := range keys {
			// io.KeyMap(int(imguiKey), nativeKey)
		}
	*/
}

func (a *Adapter) SetGameScreenSize(size imgui.Vec2) {
	if a.game.gameScreen == nil {
		dest := ebiten.NewImage(int(size.X), int(size.Y))
		a.SetGameRenderDestination(dest)
	}
	a.game.SetGameScreenSize(size)
}
