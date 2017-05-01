// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"go/build"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const logFile = "gl.log"
const windowWidth = 1440 / 1.2
const windowHeight = 900 / 1.2

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

		//glfw.SwapInterval(1) // enable vertical refresh

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

		glLogGLParams()

		//gl.Enable(gl.CULL_FACE)
		//gl.CullFace(gl.BACK)
		//gl.FrontFace(gl.CCW)

		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)

		// should be on by default, but just to make sure
		gl.Enable(gl.MULTISAMPLE)

		//gl.ClearColor(0.45, 0.5, 0.5, 1.0)
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)

		version := gl.GoStr(gl.GetString(gl.VERSION))
		glLogln(fmt.Sprintf("OpenGL Version %s", version))
	}

	cam := newCamera()

	// load mesh(es)
	cube := newCube(float32(rand.Float64()*20-10), 0, float32(rand.Float64()*20-10))

	diffuseTexture, err := newTexture("textures/crate0/crate0_diffuse.png")
	if err != nil {
		log.Fatalln(err)
	}
	specularTexture, err := newTexture("textures/specular.png")
	if err != nil {
		log.Fatalln(err)
	}

	ourShader, err := NewShader("material", "material")
	if err != nil {
		return err
	}

	whiteShader, err := NewShader("white", "white")
	if err != nil {
		return err
	}

	// this is pretty static for now. will need to be updated if window can change size
	projection := mgl32.Perspective(mgl32.DegToRad(67.0), float32(windowWidth)/windowHeight, 0.1, 100.0)

	positions := []mgl32.Vec3{
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
		{rand.Float32()*20 - 10, rand.Float32()*20 - 10, rand.Float32()*20 - 10},
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
		glfw.PollEvents()

		// update timers
		now := glfw.GetTime()
		elapsed := float32(now - previousTime)
		previousTime = now

		// update and get the camera view
		view := cam.View(elapsed)

		lightPositions := [][]float32{
			{-0.4, 1.4, -3.5},
			{0.7, 0.2, 2.0},
			{2.3, -3.3, -4.0},
			{-4.0, 2.0, -12.0},
		}
		for i := 0; i < 3; i++ {
			lightPositions[i][i] = lightPositions[i][i] + float32(math.Sin(glfw.GetTime())*2.0)
		}

		// draw the test meshes
		{
			ourShader.Use()
			setUniformMatrix4fv(ourShader, "projection", projection)
			setUniformMatrix4fv(ourShader, "view", view)

			viewPosLoc := uniformLocation(ourShader, "viewPos")
			gl.Uniform3f(viewPosLoc, cam.position[0], cam.position[1], cam.position[2])

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, diffuseTexture)
			gl.Uniform1i(uniformLocation(ourShader, "materialDiffuse"), 0)

			gl.ActiveTexture(gl.TEXTURE1)
			gl.BindTexture(gl.TEXTURE_2D, specularTexture)
			gl.Uniform1i(uniformLocation(ourShader, "materialSpecular"), 1)
			gl.Uniform1f(uniformLocation(ourShader, "materialShininess"), 32.0)

			for i := range lightPositions {
				name := fmt.Sprintf("lights[%d]", i)
				gl.Uniform4f(uniformLocation(ourShader, name+".vector"), lightPositions[i][0], lightPositions[i][1], lightPositions[i][2], 1)
				//gl.Uniform4f(uniformLocation(ourShader, "light.vector"), 1, 1, 0, 0)
				gl.Uniform3f(uniformLocation(ourShader, name+".ambient"), lightColours[i][0]/10, lightColours[i][1]/10, lightColours[i][2]/10)
				gl.Uniform3f(uniformLocation(ourShader, name+".diffuse"), lightColours[i][0], lightColours[i][1], lightColours[i][2])
				gl.Uniform3f(uniformLocation(ourShader, name+".specular"), 1.0, 1.0, 1.0)
				gl.Uniform1f(uniformLocation(ourShader, name+".constant"), 1.0)
				gl.Uniform1f(uniformLocation(ourShader, name+".linear"), 0.14)
				gl.Uniform1f(uniformLocation(ourShader, name+".quadratic"), 0.07)
			}

			for i := range positions {
				trans := mgl32.Translate3D(positions[i][0], positions[i][1], positions[i][2])
				trans = trans.Mul4(mgl32.HomogRotate3D(float32(i*20.0), mgl32.Vec3{0, 1, 0}))
				setUniformMatrix4fv(whiteShader, "transform", trans)

				gl.BindVertexArray(cube.vao)
				gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube.Vertices)))
			}
			// set back defaults, from the book of good practices
			//for i := range cube.Textures {
			//	gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
			//	gl.BindTexture(gl.TEXTURE_2D, 0)
			//}
		}

		// draw the lamp
		{
			whiteShader.Use()
			setUniformMatrix4fv(whiteShader, "projection", projection)
			setUniformMatrix4fv(whiteShader, "view", view)

			for i := range lightPositions {
				trans := mgl32.Translate3D(lightPositions[i][0], lightPositions[i][1], lightPositions[i][2])
				trans = trans.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
				setUniformMatrix4fv(whiteShader, "transform", trans)

				gl.Uniform3f(uniformLocation(whiteShader, "emissive"), lightColours[i][0], lightColours[i][1], lightColours[i][2])

				gl.BindVertexArray(cube.vao)
				gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube.Vertices)))
			}
		}

		window.SwapBuffers()
	}
	return nil
}

func setUniformMatrix4fv(shader *Shader, name string, matrix mgl32.Mat4) {
	location := uniformLocation(shader, name)
	gl.UniformMatrix4fv(location, 1, false, &matrix[0])
}

func uniformLocation(shader *Shader, name string) int32 {
	location := gl.GetUniformLocation(shader.Program, gl.Str(name+"\x00"))
	if location < 0 {
		glError(fmt.Errorf("uniform location for shader.Program '%d' and name '%s' not found", shader.Program, name))
	}
	return location
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
