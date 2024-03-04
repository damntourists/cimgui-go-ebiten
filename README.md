# cimgui-go-ebiten
Ebiten backend for [cimgui-go](https://github.com/AllenDang/cimgui-go)! 

*This is a work in progress!*

## Important!
You MUST use the following buildtags when building:
* exclude_cimgui_sdli
* exclude_cimgui_glfw

```
go run -tags exclude_cimgui_glfw,exclude_cimgui_sdli examples/basic/main.go
```
*Failure to do so will result in a number of ld failures. Ebiten has it's own GLFW 
implementation that conflicts with vanilla cimgui-go*
