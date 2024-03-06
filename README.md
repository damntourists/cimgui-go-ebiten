# cimgui-go-ebiten
[Ebiten](https://ebitengine.org/) backend for [cimgui-go](https://github.com/AllenDang/cimgui-go)!

![WIP](screenshot_wip.png)

*This is a work in progress!*
Currently there are issues with cliprects not being respected for scrolling regions, so some elements may render outside of the window or tables. I am still trying to figure out the best solution for this, but if anyone else wants to take a stab at fixing, the issue happens here https://github.com/damntourists/cimgui-go-ebiten/blob/main/render.go#L60.

Credit for rendering code goes to https://github.com/gabstv/ebiten-imgui. Unfortunately the ebiten-imgui repository hasn't been updated in a while and the version of imgui is forked and outdated. This implementation uses the new backend abstraction to allow it to work with the latest versions of imgui, from the cimgui-go repository. 


## Important!
You MUST use the following buildtags when building.
* exclude_cimgui_sdli
* exclude_cimgui_glfw

```
go run -tags exclude_cimgui_glfw,exclude_cimgui_sdli examples/basic/main.go
```
*Failure to do so will result in a number of ld failures. Ebiten has it's own GLFW 
implementation that conflicts with vanilla cimgui-go*
