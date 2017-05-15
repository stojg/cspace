// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"math/rand"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const logFile = "gl.log"
const windowWidth = 1440 * 0.9 // (1296)
const windowHeight = 900 * 0.9 // (81

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

	keys = make(map[glfw.Key]bool)
	cursor[0] = windowWidth / 2
	cursor[1] = windowHeight / 2

	if err := restartLog(); err != nil {
		return err
	}
	defer glLogln("Program stopped")

	window, err := initWindow(windowWidth, windowHeight)
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	if err := initGL(); err != nil {
		return err
	}

	scene := NewScene()
	//MaterialLevel(scene.graph)
	//ScifiLevel(scene.graph)
	ReferenceLevel(scene.graph)

	//var frame = 0
	for !window.ShouldClose() {
		//frame++
		glfw.PollEvents()
		scene.Render()
		fpsCounter(window)
		window.SwapBuffers()
		//if frame > 60*2 {
		//	window.SetShouldClose(true)
		//}
	}

	window.Destroy()

	return nil
}

var quadVAO uint32
var quadVBO uint32

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
		gl.GenBuffers(1, &quadVBO)
		gl.BindVertexArray(quadVAO)
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
