// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"go/build"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 600
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type pos [3]float32

func newSquare(index uint32, a, b pos) *square {
	s := &square{
		vbo: &index,
		vao: &index,
		vertices: []float32{
			a[0], a[1], 0.0,
			a[0], b[1], 0.0,
			b[0], a[1], 0.0,

			b[0], a[1], 0.0,
			b[0], b[1], 0.0,
			a[0], b[1], 0.0,
		},
	}
	s.setVertexBufferObject()
	s.setVertexArrayObject()
	return s
}

type square struct {
	vbo, vao *uint32
	vertices []float32
}

// vertex buffer object
// Configure the vertex data
// Now an unusual step. Most meshes will use a collection of one or more vertex buffer objects to hold vertex
// points, texture-coordinates, vertex normals, etc. In older GL implementations we would have to bind each one,
// and define their memory layout, every time that we draw the mesh. To simplify that, we have new thing called
// the vertex array object (VAO), which remembers all of the vertex buffers that you want to use, and the memory
// layout of each one. We set up the vertex array object once per mesh. When we want to draw, all we do then is
// bind the VAO and draw.
func (s *square) setVertexBufferObject() {
	gl.GenBuffers(1, s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(s.vertices)*4, gl.Ptr(s.vertices), gl.STATIC_DRAW)
}

func (s *square) setVertexArrayObject() {
	// Here we tell GL to generate a new VAO for us. It sets an unsigned integer to identify it with later.
	gl.GenVertexArrays(1, s.vao)
	// We bind it, to bring it in to focus in the state machine.
	gl.BindVertexArray(*s.vao)
	// This lets us enable the first attribute; 0. We are only using a single vertex buffer, so we know that it will
	// be attribute location 0
	gl.EnableVertexAttribArray(0)
	// The glVertexAttribPointer function defines the layout of our first vertex buffer; "0" means define the layout
	// for attribute number 0. "3" means that the variables are vec3 made from every 3 floats (GL_FLOAT) in the buffer.
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
}

func (s square) Count() int32 {
	return int32(len(s.vertices))
}

func (s square) Draw() {
	gl.BindVertexArray(*s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, s.Count())
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	square := newSquare(0, pos{0.1, 0.1}, pos{-0.1, -0.1})

	square2 := newSquare(1, pos{0.5, 0.5}, pos{0.6, 0.6})

	const vertex_shader = `
#version 410
in vec3 vp;
void main() {
	gl_Position = vec4(vp.x, vp.y, vp.z, 1.0);
}
` + "\x00"

	const fragment_shader = `
#version 410
out vec4 frag_colour;
void main() {
	frag_colour = vec4(0.2, 0.3, 0.4, 1.0);
}
` + "\x00"

	program, err := newProgram(vertex_shader, fragment_shader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.9, 0.9, 0.9, 1.0)

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
}

var prevSeconds float64
var frameCount int

func fpsCounter(window *glfw.Window) {
	currentSeconds := glfw.GetTime()
	elapsedSeconds := currentSeconds - prevSeconds
	if elapsedSeconds > 0.25 {
		prevSeconds = currentSeconds
		fps := float64(frameCount) / elapsedSeconds
		window.SetTitle(fmt.Sprintf("opengl @ fps: %.2f", fps))
		frameCount = 0
	}
	frameCount++
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

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
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

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
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
