package main

import (
	ebitenbackend "cimgui-go-ebiten"
	imgui "github.com/AllenDang/cimgui-go"
)

func main() {
	backend, _ = imgui.CreateBackend(ebitenbackend.NewBackend())
	//replace CreateTextureAs imgui.TextureID{Data:uintptr(1)} (default font textureid)
}
