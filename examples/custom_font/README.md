# Fonts in imgui
By default, imgui comes with a 13px font called ProggyClean. If a user has not set up their own fonts, 
then the builtin font is used. Unfortunately this font is limited to 13px and will be pixelated when the 
monitor's scale factor is set higher than `1`.

# HiDPI?
Please refer to https://github.com/ocornut/imgui/blob/master/docs/FAQ.md#q-how-should-i-handle-dpi-in-my-application if you have any questions about how HiDPI is handled in imgui. 
This is typically up to the user to implement and is implemented in this example. 

This library will also scale the current imgui style during window creation.


# Building 
Please be sure to include these tags when building to disable sdl and 
glfw backends in imgui:
* `exclude_cimgui_sdl`
* `exclude_cimgui_glfw`

# Sharpie font
The font used in this example is Sharpie from https://www.fontshare.com/. Their license, 
however, does not allow for me to host the files here on github to the public. 
You can download them from Fontshare's  site at:
https://www.fontshare.com/fonts/sharpie

You can also download and decompress the font by running the following:
```shell
wget https://api.fontshare.com/v2/fonts/download/sharpie -O Sharpie_Complete.zip && unzip Sharpie_Complete.zip
```