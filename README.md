# cimgui-go-ebiten
[Ebiten](https://ebitengine.org/) backend for [cimgui-go](https://github.com/AllenDang/cimgui-go)!

![WIP](screenshot_wip.png)

*This is a work in progress!*

Credit for rendering code goes to https://github.com/gabstv/ebiten-imgui. Unfortunately the ebiten-imgui repository hasn't been updated in a while and the version of imgui is forked and outdated. This implementation uses the new backend abstraction to allow it to work with the latest versions of imgui, from the cimgui-go repository. 


## Important!
This repository is designed to work with https://github.com/AllenDang/cimgui-go. A caveat
of this, however, is that you _must_ include the following buildtags:
* exclude_cimgui_sdl
* exclude_cimgui_glfw

example:
```
go run -tags exclude_cimgui_glfw,exclude_cimgui_sdl examples/basic/main.go
```
*Failure to do so will result in a number of ld failures. Ebiten has its own GLFW 
implementation that conflicts with vanilla cimgui-go*


## TODO
    * Add gamepad support