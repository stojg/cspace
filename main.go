// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"math"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const logFile = "gl.log"
const windowWidth = 1000
const windowHeight = 800

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

	// this is pretty static for now. will need to be updated if window can change size
	projection := mgl32.Perspective(mgl32.DegToRad(67.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	cam := newCamera()

	// load Mesh(es)
	cubeMesh := newCrateMesh()
	lightMesh := newLightMesh()
	floor := newPlaneMesh()

	ourShader, err := NewShader("material", "material")
	if err != nil {
		return err
	}

	lampShader, err := NewShader("simple", "emissive")
	if err != nil {
		return err
	}

	floorShader, err := NewShader("material", "floor")
	if err != nil {
		return err
	}

	positions := []mgl32.Vec3{
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}

	useNormalMapping := true

	directionalColor := [3]float32{1.000, 0.99, 0.9}

	previousTime := glfw.GetTime()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if keys[glfw.Key1] {
			useNormalMapping = true
		}
		if keys[glfw.Key2] {
			useNormalMapping = false
		}

		fpsCounter(window)

		// update timers
		now := glfw.GetTime()
		elapsed := float32(now - previousTime)
		previousTime = now

		// update and get the camera view
		view := cam.View(elapsed)

		sin := float32(math.Sin(now))

		// draw the test meshes
		ourShader.UsePV(projection, view)
		gl.Uniform3f(uniformLocation(ourShader, "viewPos"), cam.position[0], cam.position[1], cam.position[2])
		setDirectionalLight(ourShader, [3]float32{0, -5, 0}, directionalColor)
		//setLights(floorShader, lightPositions, lightColours)
		for i := range positions {
			trans := mgl32.Translate3D(positions[i][0], positions[i][1], positions[i][2])
			trans = trans.Mul4(mgl32.HomogRotate3D(sin+float32(i*20.0), mgl32.Vec3{1, 1, 1}.Normalize()))
			setUniformMatrix4fv(ourShader, "transform", trans)
			location := uniformLocation(ourShader, "useNormalMapping")
			if useNormalMapping {
				gl.Uniform1f(location, 1.0)
			} else {
				gl.Uniform1f(location, 0.0)
			}
			cubeMesh.Draw(ourShader)
		}

		// draw the floor
		floorShader.UsePV(projection, view)
		gl.Uniform3f(uniformLocation(floorShader, "viewPos"), cam.position[0], cam.position[1], cam.position[2])
		setDirectionalLight(ourShader, [3]float32{0, -5, 0}, directionalColor)
		trans := mgl32.Translate3D(0, -5, 0)
		trans = trans.Mul4(mgl32.Scale3D(100, 0.1, 100))
		setUniformMatrix4fv(floorShader, "transform", trans)
		floor.Draw(floorShader)

		// draw the lamps
		//_ = lampShader
		//_ = lightMesh
		lampShader.UsePV(projection, view)
		trans = mgl32.Translate3D(0, 5, 0)
		trans = trans.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
		setUniformMatrix4fv(lampShader, "transform", trans)
		gl.Uniform3f(uniformLocation(lampShader, "emissive"), 1, 1, 1)
		lightMesh.Draw(lampShader)

		window.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

func setDirectionalLight(shader *Shader, direction, color [3]float32) {
	name := fmt.Sprint("lights[0]")
	gl.Uniform4f(uniformLocation(shader, name+".vector"), direction[0], direction[1], direction[2], 0)
	gl.Uniform3f(uniformLocation(shader, name+".diffuse"), color[0], color[1], color[2])
	gl.Uniform3f(uniformLocation(shader, name+".ambient"), color[0]/10, color[1]/10, color[2]/10)
	gl.Uniform3f(uniformLocation(shader, name+".specular"), 1.0, 1.0, 1.0)
}

func setLights(shader *Shader, pos, color [][]float32) {
	for i := range pos {
		name := fmt.Sprintf("lights[%d]", i)
		gl.Uniform4f(uniformLocation(shader, name+".vector"), pos[i][0], pos[i][1], pos[i][2], 1)
		gl.Uniform3f(uniformLocation(shader, name+".diffuse"), color[i][0], color[i][1], color[i][2])
		gl.Uniform3f(uniformLocation(shader, name+".ambient"), color[i][0]/10, color[i][1]/10, color[i][2]/10)
		gl.Uniform3f(uniformLocation(shader, name+".specular"), color[i][0], color[i][1], color[i][2])
		gl.Uniform1f(uniformLocation(shader, name+".constant"), 1.0)
		gl.Uniform1f(uniformLocation(shader, name+".linear"), 0.14)
		gl.Uniform1f(uniformLocation(shader, name+".quadratic"), 0.07)
	}
}
