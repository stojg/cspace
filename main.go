// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"go/build"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const logFile = "gl.log"
const windowWidth = 600
const windowHeight = 600

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
	if err := restartLog(); err != nil {
		return err
	}
	defer glLogln("Program stopped")

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

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		return err
	}

	glLogGLParams()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	glLogln(fmt.Sprintf("OpenGL Version %s", version))

	square := newSquare(0, pos{0.1, 0.1}, pos{-0.1, -0.1})
	square2 := newSquare(1, pos{0.5, 0.5}, pos{0.6, 0.6})

	program, err := loadTestShader()
	if err != nil {
		glError(err)
	}

	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key != glfw.KeySpace || action != glfw.Press {
			return
		}
		newProgram, err := loadTestShader()
		if err != nil {
			glError(err)
			return
		}
		program = newProgram
		fmt.Println("shaders reloaded")
	})

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.6, 0.6, 0.8, 1.0)

	//gl.PolygonMode(gl.FRONT, gl.LINE)

	for !window.ShouldClose() {
		fpsCounter(window)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(program)
		square.Draw()
		square2.Draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()

	}
	return nil
}

func loadTestShader() (uint32, error) {
	vertex_shader, err := loadVertexShader("test")
	if err != nil {
		return 0, err
	}
	fragment_shader, err := loadFragShader("test")
	if err != nil {
		return 0, err
	}
	program, err := newProgram(vertex_shader, fragment_shader)
	if err != nil {
		return 0, err
	}
	gl.UseProgram(program)
	glLogProgramme(program)
	colorUniform := gl.GetUniformLocation(program, gl.Str("inputColor\x00"))
	gl.Uniform4f(colorUniform, 0, 0.8, 0, 4)
	return program, nil
}

func loadVertexShader(name string) (string, error) {
	res, err := ioutil.ReadFile(filepath.Join("shaders", fmt.Sprintf("%s.vert", name)))
	return string(res) + "\x00", err
}

func loadFragShader(name string) (string, error) {
	res, err := ioutil.ReadFile(filepath.Join("shaders", fmt.Sprintf("%s.frag", name)))
	return string(res) + "\x00", err
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		l := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(l))

		return 0, fmt.Errorf("failed to link program[%d]: %v", program, l)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		l := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(l))
		return 0, fmt.Errorf("failed to compile \n%v \n%v", l, source)
	}
	return shader, nil
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
