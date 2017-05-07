package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderI interface {
	Use()
	UsePV(projection, view mgl32.Mat4)
	Program() uint32
}

type LightShader struct {
	*Shader
	uniformPosLoc        int32
	uniformNormalLoc     int32
	uniformAlbedoSpecLoc int32
}

func NewLightShader(vertex, frag string) (*LightShader, error) {
	c, err := NewShader(vertex, frag)
	if err != nil {
		return nil, err
	}
	s := &LightShader{
		Shader:               c,
		uniformPosLoc:        uniformLocation(c, "gPosition"),
		uniformNormalLoc:     uniformLocation(c, "gNormal"),
		uniformAlbedoSpecLoc: uniformLocation(c, "gAlbedoSpec"),
	}
	return s, nil
}

type Shader struct {
	program uint32
}

func (s *Shader) Program() uint32 {
	return s.program
}

func (s *Shader) Use() {
	gl.UseProgram(s.program)
}

func (s *Shader) UsePV(projection, view mgl32.Mat4) {
	gl.UseProgram(s.program)
	setUniformMatrix4fv(s, "projection", projection)
	setUniformMatrix4fv(s, "view", view)
}

func NewShader(vertex, frag string) (*Shader, error) {
	shader := &Shader{}
	vertexShaderSource, err := loadVertexShader(vertex)
	if err != nil {
		return shader, err
	}
	fragmentShaderSource, err := loadFragShader(frag)
	if err != nil {
		return shader, err
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return shader, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return shader, err
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

		return shader, fmt.Errorf("failed to link program[%d]: %v", program, l)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	shader.program = program

	glLogShader(shader)

	return shader, nil
}

func loadVertexShader(name string) (string, error) {
	res, err := ioutil.ReadFile(filepath.Join("shaders", fmt.Sprintf("%s.vert", name)))
	return string(res) + "\x00", err
}

func loadFragShader(name string) (string, error) {
	res, err := ioutil.ReadFile(filepath.Join("shaders", fmt.Sprintf("%s.frag", name)))
	return string(res) + "\x00", err
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
