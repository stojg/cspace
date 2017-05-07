// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const logFile = "gl.log"
const windowWidth = 1200
const windowHeight = 750

var keys map[glfw.Key]bool
var cursor [2]float64

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

	window, err := initWindow(windowWidth, windowHeight)
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	if err := initGL(); err != nil {
		return err
	}

	scene := NewScene(windowWidth, windowHeight)

	crate := NewCrateModel()
	for i := 0; i < 10; i++ {
		t := mgl32.Translate3D(rand.Float32()*20-10, 0.5, rand.Float32()*20-10)
		t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
		scene.graph.Add(crate, t)
	}
	grass := NewGrassMesh()
	for x := 0; x < 10; x++ {
		for z := 0; z < 10; z++ {
			grassT := mgl32.Translate3D(float32(x)*3-15, 0, float32(z)*3-15)
			grassT = grassT.Mul4(mgl32.Scale3D(3, 0.1, 3))
			scene.graph.Add(grass, grassT)
		}
	}

	for !window.ShouldClose() {
		glfw.PollEvents()
		scene.Render()
		fpsCounter(window)
		window.SwapBuffers()
	}
	return nil
}

var quadVAO uint32 = 0
var quadVBO uint32 = 0

func renderQuad() {
	if quadVAO == 0 {
		quadVertices := []float32{
			// Positions        // Texture Coords
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

func drawScene(shader *Shader, cubeMesh, floorMesh *Mesh) {

	trans := mgl32.Translate3D(0, -0.5, 0)
	trans = trans.Mul4(mgl32.Scale3D(25, 0.1, 25))
	setUniformMatrix4fv(shader, "model", trans)
	floorMesh.Render(shader)

	model := mgl32.Translate3D(0, 1.5, 0)
	setUniformMatrix4fv(shader, "model", model)
	cubeMesh.Render(shader)

	model = mgl32.Translate3D(2.0, 0.0, 1)
	setUniformMatrix4fv(shader, "model", model)
	cubeMesh.Render(shader)

	model = mgl32.Translate3D(-1.0, 0.0, 2)
	model = model.Mul4(mgl32.HomogRotate3D(45, mgl32.Vec3{1, 0, 1}.Normalize()))
	model = model.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

	setUniformMatrix4fv(shader, "model", model)
	cubeMesh.Render(shader)

}

func setDirectionalLight(shader *Shader, direction, color [3]float32) {
	name := fmt.Sprint("lights[0]")
	gl.Uniform4f(uniformLocation(shader, name+".vector"), direction[0], direction[1], direction[2], 0)
	gl.Uniform3f(uniformLocation(shader, name+".diffuse"), color[0], color[1], color[2])
	gl.Uniform3f(uniformLocation(shader, name+".ambient"), color[0]/10, color[1]/10, color[2]/10)
	gl.Uniform3f(uniformLocation(shader, name+".specular"), 1.0, 1.0, 1.0)
}

func setLights(shader *Shader, lights []*PointLight) {
	for i, light := range lights {
		name := fmt.Sprintf("lights[%d]", i)
		gl.Uniform3f(uniformLocation(shader, name+".Direction"), light.Position[0], light.Position[1], light.Position[2])
		gl.Uniform3f(uniformLocation(shader, name+".Color"), light.Color[0], light.Color[1], light.Color[2])
		gl.Uniform1f(uniformLocation(shader, name+".Radius"), light.Radius())
		gl.Uniform1f(uniformLocation(shader, name+".Linear"), light.Linear)
		gl.Uniform1f(uniformLocation(shader, name+".Quadratic"), light.Exp)
		//gl.Uniform1f(uniformLocation(shader, name+".AmbientIntensity"), light.AmbientIntensity)
		gl.Uniform1f(uniformLocation(shader, name+".DiffuseIntensity"), light.DiffuseIntensity)
	}
}
