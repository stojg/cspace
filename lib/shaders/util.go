package shaders

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func Use(program uint32) {
	gl.UseProgram(program)
}

func SetMatrix4f(loc int32, mat mgl32.Mat4) {
	gl.UniformMatrix4fv(loc, 1, false, &mat[0])
}

func buildShader(vertex, frag string) uint32 {
	vertexShaderSource, err := loadVertexShader(vertex)
	if err != nil {
		panic(err)
	}
	fragmentShaderSource, err := loadFragShader(frag)
	if err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
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

		panic(fmt.Errorf("failed to link program[%d]: %v", program, l))
	}

	gl.DetachShader(program, vertexShader)
	gl.DeleteShader(vertexShader)
	gl.DetachShader(program, fragmentShader)
	gl.DeleteShader(fragmentShader)

	//glLogShader(program, vertex, frag)
	return program
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

func loc(program uint32, name string) int32 {
	location := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	if location < 0 {
		panic(fmt.Sprintf("uniform location for tShader.Program '%d' and name '%s' not found", program, name))
	}
	return location
}
