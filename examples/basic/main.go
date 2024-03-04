package main

import (
	imgui "github.com/AllenDang/cimgui-go"
	ebitenbackend "github.com/damntourists/cimgui-go-ebiten"
)

func main() {
	// You MUST use the following buildtags when building:
	// * exclude_cimgui_sdli
	// * exclude_cimgui_glfw

	// TODO: Work in progress... currently not possible to cast to Backend
	b, _ := imgui.CreateBackend(ebitenbackend.NewBackend())
	b.SetWindowTitle("hello")

}
