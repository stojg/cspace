package main

import (
	"fmt"
	"go/build"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func uniformLocation(shader Shader, name string) int32 {
	location := gl.GetUniformLocation(shader.Program(), gl.Str(name+"\x00"))
	if location < 0 {
		glError(fmt.Errorf("uniform location for tShader.Program '%d' and name '%s' not found", shader.Program(), name))
	}
	return location
}

func pUniformLocation(program uint32, name string) int32 {
	location := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	if location < 0 {
		glError(fmt.Errorf("uniform location for tShader.Program '%d' and name '%s' not found", program, name))
	}
	return location
}

func initWindow() (*glfw.Window, error) {
	// initialise and setup a window with user inputs
	var window *glfw.Window
	if err := glfw.Init(); err != nil {
		return window, fmt.Errorf("failed to initialize glfw: %s", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 0)

	var err error
	window, err = glfw.CreateWindow(int(windowWidth), int(windowHeight), "cspace", nil, nil)
	if err != nil {
		return window, err
	}
	window.MakeContextCurrent()
	// disable or enable vertical refresh (vsync)
	glfw.SwapInterval(0)
	mouseCaptured := true
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Release && key == glfw.KeySpace {
			if !mouseCaptured {
				window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
				mouseCaptured = true
			} else {
				window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
				mouseCaptured = false
			}
		}

		if action == glfw.Press {
			keys[key] = true
		} else if action == glfw.Release {
			keys[key] = false
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		if mouseCaptured {
			cursor[0] = xpos
			cursor[1] = ypos
		}
	})

	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		fboWidth, fboHeight := w.GetFramebufferSize()
		viewPortWidth = int32(fboWidth)
		viewPortHeight = int32(fboHeight)
		windowWidth = int32(width)
		windowHeight = int32(height)
	})

	w, h := window.GetFramebufferSize()
	viewPortWidth = int32(w)
	viewPortHeight = int32(h)
	return window, nil
}

func initGL() error {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		return err
	}

	gl.Disable(gl.MULTISAMPLE)
	glLogGLParams()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	glLogln(fmt.Sprintf("OpenGL Version %s", version))

	//gl.Enable(gl.FRAMEBUFFER_SRGB)
	return nil
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
