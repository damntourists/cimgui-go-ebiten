package main

import (
	ebitenbackend "cimgui-go-ebiten"
	imgui "github.com/AllenDang/cimgui-go"
)

func main() {
	// TODO: Work in progress... currently not possible to cast to Backend
	backend, _ = imgui.CreateBackend(ebitenbackend.NewBackend())
}
