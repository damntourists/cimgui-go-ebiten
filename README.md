# cimgui-go-ebiten
[Ebiten](https://ebitengine.org/) backend for [cimgui-go](https://github.com/AllenDang/cimgui-go)!

![WIP](screenshot_wip.png)

*This is a work in progress!*

Credit for rendering code goes to https://github.com/gabstv/ebiten-imgui. Unfortunately the ebiten-imgui repository hasn't been updated in a while and the version of imgui is forked and outdated. This implementation uses the new backend abstraction to allow it to work with the latest versions of imgui, from the cimgui-go repository. 


## Important!
This repository is designed to work with https://github.com/AllenDang/cimgui-go. A caveat
of this, however, is that you _must_ include the following buildtags:
* exclude_cimgui_sdli
* exclude_cimgui_glfw

example:
```
go run -tags exclude_cimgui_glfw,exclude_cimgui_sdli examples/basic/main.go
```
*Failure to do so will result in a number of ld failures. Ebiten has it's own GLFW 
implementation that conflicts with vanilla cimgui-go*

### _**However**_

This repository has a `replace` line added to the `go.mod` file which redirects AllanDang's 
cimgui-go repository to a version that I strip of these third party dependencies. Using the 
lite version is totally optional and removes the need to include the build flags mentioned above.

```
require (
	github.com/AllenDang/cimgui-go v0.0.0-20240303223020-2c7d7a8d1731
	github.com/hajimehoshi/ebiten/v2 v2.6.6
)

replace github.com/AllenDang/cimgui-go => github.com/damntourists/cimgui-go-lite v1.0.0
```
