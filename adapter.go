package ebitenbackend

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"runtime"
)


// Adapter should proxy calls to backend.
type Adapter interface {
	SetAfterCreateContextHook(func())
	SetBeforeDestroyContextHook(func())
	SetBeforeRenderHook(func())
	SetAfterRenderHook(func())

	SetBgColor(color Vec4)
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

	//SetDropCallback(imgui.DropCallback)
	//SetCloseCallback(imgui.WindowCloseCallback[BackendFlagsT])
	//SetKeyCallback(KeyCallback)
	//SetSizeChangeCallback(SizeChangeCallback)
	// SetWindowFlags selected hint to specified value.
	// ATTENTION: This method is able to set only one flag per call.
	//SetWindowFlags(flag BackendFlagsT, value int)
	SetIcons(icons ...image.Image)
	CreateWindow(title string, width, height int)



	//CreateWindow(title string, width, height int)
	//SetWindowPos(x, y int)
	//Run(func())
	Backend() *imgui.Backend[EbitenWindowFlags]
	SetGame(ebiten.Game)
	SetUILoop(func())
	UILoop()
	Game() ebiten.Game
	Update(float32)
	finalize()
}

type EbitenAdapter struct {
	backend imgui.Backend[EbitenWindowFlags]
	game    ebiten.Game
	loop    func()
}

func (a *EbitenAdapter) SetBeforeDestroyContextHook(f func()) {
	a.
}

func (a *EbitenAdapter) SetBeforeRenderHook(f func()) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetAfterRenderHook(f func()) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetBgColor(color interface{}) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) Refresh() {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) GetWindowPos() (x, y int32) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetWindowSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetWindowSizeLimits(minWidth, minHeight, maxWidth, maxHeight int) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetWindowTitle(title string) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) DisplaySize() (width, height int32) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetShouldClose(b bool) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) ContentScale() (xScale, yScale float32) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetTargetFPS(fps uint) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) SetIcons(icons ...image.Image) {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) Backend() *imgui.Backend[EbitenWindowFlags] {
	//TODO implement me
	panic("implement me")
}

func (a *EbitenAdapter) UILoop() {
	//TODO implement me
	panic("implement me")
}

func NewEbitenAdapter() *EbitenAdapter {
	b := &BackendBridge{
		ctx:    nil,
		cache:  NewCache(),
		filter: ebiten.FilterNearest,
	}

	b.ctx = imgui.CreateContext()

	fonts := imgui.CurrentIO().Fonts()
	_, _, _, _ = fonts.GetTextureDataAsRGBA32() // call this to force imgui to build the font atlas cache

	txid := imgui.TextureID{Data: uintptr(1)}
	fonts.SetTexID(txid)
	b.cache.SetFontAtlasTextureID(txid)

	bb := (imgui.Backend[EbitenWindowFlags])(b)

	createdBackend, _ := imgui.CreateBackend(bb)

	a := EbitenAdapter{
		backend: createdBackend,
	}

	a.backend.SetBeforeRenderHook(func() {
		println("before render hook!")
	})
	a.backend.SetAfterRenderHook(func() {
		println("after render hook!")
	})

	runtime.SetFinalizer(&a, (*EbitenAdapter).finalize)

	return &a
}

func (a *EbitenAdapter) finalize() {
	runtime.SetFinalizer(a, nil)
}

func (a *EbitenAdapter) SetGame(g ebiten.Game) {
	//a.game = g
	a.game = g
}

func (a *EbitenAdapter) Game() ebiten.Game {
	println("ebiten adapter Game() called. ")
	return a.game
}

func (a *EbitenAdapter) SetUILoop(f func()) {
	a.loop = f
}

func (a *EbitenAdapter) UILoop() func() {
	return a.loop
}

func (a *EbitenAdapter) Update(delta float32) {
	println("adapter update!", delta)
	io := imgui.CurrentIO()
	io.SetDeltaTime(delta)

	_ = a.game.Update()
	//if err != nil {
	//	return
	//}

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
