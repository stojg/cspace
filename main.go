// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"go/build"
	"os"

	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const logFile = "gl.log"
const windowWidth = 320
const windowHeight = 200

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

	uniformBlockIndex := gl.GetUniformBlockIndex(ourShader.Program, gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(ourShader.Program, uniformBlockIndex, 0)

	//// this is pretty static for now. will need to be updated if window can change size
	projection := mgl32.Perspective(mgl32.DegToRad(67.0), float32(windowWidth)/windowHeight, 0.1, 100.0)

	cam := newCamera()

	positions := []mgl32.Vec3{
		{0, 0, 0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}

	lightPositions := [][]float32{
		{-0.4, 1.4, -3.5},
		{0.7, 0.2, 2.0},
		{2.3, -3.3, -4.0},
		{-4.0, 2.0, -12.0},
	}

	lightColours := [][]float32{
		{0.8, 0.5, 0.5},
		{0.5, 0.8, 0.5},
		{0.5, 0.5, 0.8},
		{1.000, 0.749, 0.000},
	}

	previousTime := glfw.GetTime()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

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
		setLights(floorShader, lightPositions, lightColours)
		for i := range positions {
			trans := mgl32.Translate3D(positions[i][0], positions[i][1], positions[i][2])
			trans = trans.Mul4(mgl32.HomogRotate3D(sin+float32(i*20.0), mgl32.Vec3{1, 1, 1}.Normalize()))
			setUniformMatrix4fv(ourShader, "transform", trans)
			cubeMesh.Draw(ourShader)
		}

		// draw the floor
		floorShader.UsePV(projection, view)
		gl.Uniform3f(uniformLocation(floorShader, "viewPos"), cam.position[0], cam.position[1], cam.position[2])
		setLights(floorShader, lightPositions, lightColours)
		trans := mgl32.Translate3D(0, -5, 0)
		trans = trans.Mul4(mgl32.Scale3D(100, 0.1, 100))
		setUniformMatrix4fv(floorShader, "transform", trans)
		floor.Draw(floorShader)

		// draw the lamps
		lampShader.UsePV(projection, view)
		//gl.Uniform3f(uniformLocation(lampShader, "viewPos"), cam.position[0], cam.position[1], cam.position[2])
		for i := range lightPositions {
			trans := mgl32.Translate3D(lightPositions[i][0], lightPositions[i][1], lightPositions[i][2])
			trans = trans.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
			setUniformMatrix4fv(lampShader, "transform", trans)

			gl.Uniform3f(uniformLocation(lampShader, "emissive"), lightColours[i][0], lightColours[i][1], lightColours[i][2])

			lightMesh.Draw(lampShader)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

func setLights(shader *Shader, pos, color [][]float32) {
	for i := range pos {
		name := fmt.Sprintf("lights[%d]", i)
		gl.Uniform4f(uniformLocation(shader, name+".vector"), pos[i][0], pos[i][1], pos[i][2], 1)
		gl.Uniform3f(uniformLocation(shader, name+".ambient"), color[i][0]/10, color[i][1]/10, color[i][2]/10)
		gl.Uniform3f(uniformLocation(shader, name+".diffuse"), color[i][0], color[i][1], color[i][2])
		gl.Uniform3f(uniformLocation(shader, name+".specular"), 1.0, 1.0, 1.0)
		gl.Uniform1f(uniformLocation(shader, name+".constant"), 1.0)
		gl.Uniform1f(uniformLocation(shader, name+".linear"), 0.14)
		gl.Uniform1f(uniformLocation(shader, name+".quadratic"), 0.07)
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
