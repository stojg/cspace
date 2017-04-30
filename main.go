// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"go/build"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const logFile = "gl.log"
const windowWidth = 800
const windowHeight = 600

var keys map[glfw.Key]bool
var cursor [2]float64

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	err := realMain()
	if err != nil {
		glError(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func realMain() error {

	keys = make(map[glfw.Key]bool)
	cursor[0] = windowWidth / 2
	cursor[1] = windowHeight / 2

	if err := restartLog(); err != nil {
		return err
	}
	defer glLogln("Program stopped")

	// initialise and setup a window with user inputs
	var window *glfw.Window
	{
		if err := glfw.Init(); err != nil {
			return fmt.Errorf("failed to initialize glfw: %s", err)
		}
		defer glfw.Terminate()

		glfw.WindowHint(glfw.Resizable, glfw.False)
		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		glfw.WindowHint(glfw.Samples, 16)

		var err error
		window, err = glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
		if err != nil {
			return err
		}
		window.MakeContextCurrent()
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
	}

	// Initialize Glow
	{
		if err := gl.Init(); err != nil {
			return err
		}
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
		gl.FrontFace(gl.CCW)

		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)

		glLogGLParams()

		version := gl.GoStr(gl.GetString(gl.VERSION))
		glLogln(fmt.Sprintf("OpenGL Version %s", version))
	}

	testShader, err := NewShader("test", "test")
	if err != nil {
		return err
	}

	// setup projection and model (world)
	{
		gl.UseProgram(testShader.Program)
		projection := mgl32.Perspective(mgl32.DegToRad(67.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
		projectionUniform := gl.GetUniformLocation(testShader.Program, gl.Str("projection\x00"))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

		model := mgl32.Ident4()
		modelUniform := gl.GetUniformLocation(testShader.Program, gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	}

	cam := newCamera()
	cam.Draw(testShader)

	// load meshes
	square := newCube(1)
	square.Texture, err = newTexture("square_running.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	gl.ClearColor(0.45, 0.5, 0.5, 1.0)
	previousTime := glfw.GetTime()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		fpsCounter(window)
		glfw.PollEvents()

		// Update
		time := glfw.GetTime()
		elapsed := float32(time - previousTime)
		previousTime = time

		cam.Update(testShader, elapsed)

		// Render
		square.Draw(testShader)
		//square2.Draw()

		// Maintenance
		window.SwapBuffers()

	}
	return nil
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/stojg/cspace")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
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
