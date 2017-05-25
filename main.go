// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const logFile = "gl.log"

var keys map[glfw.Key]bool
var cursor [2]float64

func main() {
	rand.Seed(19)

	err := realMain()
	if err != nil {
		glError(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func realMain() error {

	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	// Set the working directory to the root of Go package, so that its assets can be accessed.
	dir, err := importPathToDir("github.com/stojg/cspace")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}

	if err = os.Chdir(dir); err != nil {
		log.Panicln("os.Chdir:", err)
	}

	keys = make(map[glfw.Key]bool)
	cursor[0] = float64(windowWidth / 2)
	cursor[1] = float64(windowHeight / 2)

	if err := restartLog(); err != nil {
		return err
	}
	defer glLogln("Program stopped")

	window, err := initWindow()
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	if err := initGL(); err != nil {
		return err
	}

	scene := NewScene()
	scene.Init()

	PBRLevel(scene.graph)

	for !window.ShouldClose() {
		glfw.PollEvents()
		scene.Render()
		fpsCounter(window)
		window.SwapBuffers()
	}

	window.Destroy()

	return nil
}

var quadVAO uint32

// renderQuad renders a full screen quad
func renderQuad() {
	if quadVAO == 0 {
		quadVertices := []float32{
			-1, 1, 0.0, 0.0, 1.0,
			-1, -1, 0.0, 0.0, 0.0,
			1, 1, 0.0, 1.0, 1.0,
			1, -1, 0.0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &quadVAO)
		gl.BindVertexArray(quadVAO)
		var quadVBO uint32
		gl.GenBuffers(1, &quadVBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	}
	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
}
