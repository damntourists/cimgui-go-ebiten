package ebitenbackend

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
	"runtime"
)

type Adapter interface {
	Backend() *imgui.Backend[EbitenWindowFlags]
	SetGame(ebiten.Game)
	Game() ebiten.Game
	finalize()
}

type EbitenAdapter struct {
	backend *BackendBridge
	game    ebiten.Game
}

func NewEbitenAdapter() *EbitenAdapter {
	return &EbitenAdapter{
		backend: nil,
	}
}

func (a *EbitenAdapter) Backend() imgui.Backend[EbitenWindowFlags] {
	b := &BackendBridge{
		cache:  NewCache(),
		filter: ebiten.FilterNearest,
	}

	runtime.SetFinalizer(a, (*EbitenAdapter).finalize)

	bb := (imgui.Backend[EbitenWindowFlags])(b)
	return bb
}

func (a *EbitenAdapter) finalize() {
	if a.backend != nil {
		a.backend.ctx.Destroy()
	}
	runtime.SetFinalizer(a, nil)

}

func (a *EbitenAdapter) SetGame(g ebiten.Game) {
	a.game = g
}

func (a *EbitenAdapter) Game() ebiten.Game {
	return a.game
}
