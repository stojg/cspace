package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func setUniformMatrix4fv(shader ShaderI, name string, matrix mgl32.Mat4) {
	location := uniformLocation(shader, name)
	gl.UniformMatrix4fv(location, 1, false, &matrix[0])
}

func uniformLocation(shader ShaderI, name string) int32 {
	location := gl.GetUniformLocation(shader.Program(), gl.Str(name+"\x00"))
	if location < 0 {
		glError(fmt.Errorf("uniform location for shader.Program '%d' and name '%s' not found", shader.Program(), name))
	}
	return location
}

func initWindow(width, height int) (*glfw.Window, error) {
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
	window, err = glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		return window, err
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(0) // disable vertical refresh (vsync)
	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			keys[key] = true
		} else if action == glfw.Release {
			keys[key] = false
		}
	})
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		cursor[0] = xpos
		cursor[1] = ypos
	})
	return window, nil
}

func initGL() error {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		return err
	}

	glLogGLParams()

	gl.Enable(gl.DEPTH_TEST)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	glLogln(fmt.Sprintf("OpenGL Version %s", version))
	return nil
}
