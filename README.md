# UPDATE 2024-09-18
**This backend is [now available](https://github.com/AllenDang/cimgui-go/pull/296) in [cimgui-go](https://github.com/AllenDang/cimgui-go)! 🎉🎉**

I will be archiving this repositoring soon since it is now redundant. 


# cimgui-go-ebiten
[Ebiten](https://ebitengine.org/) backend for [cimgui-go](https://github.com/AllenDang/cimgui-go)!

![WIP](screenshot_wip.png)

*This is a work in progress!*

Credit for rendering code goes to https://github.com/gabstv/ebiten-imgui. Unfortunately the ebiten-imgui repository hasn't been updated in a while and the version of imgui is forked and outdated. This implementation uses the new backend abstraction to allow it to work with the latest versions of imgui, from the cimgui-go repository. 


## Important!
This repository is designed to work with https://github.com/AllenDang/cimgui-go; however,
the cimgui-go repository has their own sdl and glfw implementations. Luckly,these can be 
disabled at build that we can use ebiten's window by supplying the following build tags:
* exclude_cimgui_sdl
* exclude_cimgui_glfw

example:
```
go run -tags exclude_cimgui_glfw,exclude_cimgui_sdl examples/ui_over_game/main.go
```
*Trying to build without the cimgui exclude build tags will cause the compiler to throw
a bunch of ld errors.*


## TODO
* Add gamepad support
