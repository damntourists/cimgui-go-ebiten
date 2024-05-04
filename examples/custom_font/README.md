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