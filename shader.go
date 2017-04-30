package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	Program uint32
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

	shader.Program = program
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
