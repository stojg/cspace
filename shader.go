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

type GbufferLightShader interface {
	ShaderI
	UniformPosLoc() int32
	UniformNormalLoc() int32
	UniformAlbedoSpecLoc() int32
}

type PointLightShader struct {
	*Shader
	uniformPosLoc        int32
	uniformNormalLoc     int32
	uniformAlbedoSpecLoc int32

	uniformLightPosLoc              int32
	uniformLightColorLoc            int32
	uniformLightDiffuseIntensityLoc int32
	uniformLightLinearLoc           int32
	uniformLightQuadraticLoc        int32
}

func (s *PointLightShader) UniformPosLoc() int32 {
	return s.uniformPosLoc
}

func (s *PointLightShader) UniformNormalLoc() int32 {
	return s.uniformNormalLoc
}

func (s *PointLightShader) UniformAlbedoSpecLoc() int32 {
	return s.uniformAlbedoSpecLoc
}

func (s *PointLightShader) SetLight(light *PointLight) {
	gl.Uniform3f(s.uniformLightPosLoc, light.Position[0], light.Position[1], light.Position[2])
	gl.Uniform3f(s.uniformLightColorLoc, light.Color[0], light.Color[1], light.Color[2])
	gl.Uniform1f(s.uniformLightDiffuseIntensityLoc, light.DiffuseIntensity)
	gl.Uniform1f(s.uniformLightLinearLoc, light.Linear)
	gl.Uniform1f(s.uniformLightQuadraticLoc, light.Exp)
}

func NewPointLightShader(vertex, frag string) (*PointLightShader, error) {
	c, err := NewShader(vertex, frag)
	if err != nil {
		return nil, err
	}
	s := &PointLightShader{
		Shader:               c,
		uniformPosLoc:        uniformLocation(c, "gPosition"),
		uniformNormalLoc:     uniformLocation(c, "gNormal"),
		uniformAlbedoSpecLoc: uniformLocation(c, "gAlbedoSpec"),

		uniformLightPosLoc:              uniformLocation(c, "pointLight.Position"),
		uniformLightColorLoc:            uniformLocation(c, "pointLight.Color"),
		uniformLightDiffuseIntensityLoc: uniformLocation(c, "pointLight.DiffuseIntensity"),
		uniformLightLinearLoc:           uniformLocation(c, "pointLight.Linear"),
		uniformLightQuadraticLoc:        uniformLocation(c, "pointLight.Quadratic"),
	}
	return s, nil
}

type DirLightShader struct {
	*Shader
	uniformPosLoc        int32
	uniformNormalLoc     int32
	uniformAlbedoSpecLoc int32

	uniformLightDirectionLoc        int32
	uniformLightColorLoc            int32
	uniformLightDiffuseIntensityLoc int32
}

func (s *DirLightShader) UniformPosLoc() int32 {
	return s.uniformPosLoc
}

func (s *DirLightShader) UniformNormalLoc() int32 {
	return s.uniformNormalLoc
}

func (s *DirLightShader) UniformAlbedoSpecLoc() int32 {
	return s.uniformAlbedoSpecLoc
}

func (s *DirLightShader) SetLight(light *DirectionalLight) {
	gl.Uniform3f(s.uniformLightDirectionLoc, light.Direction[0], light.Direction[1], light.Direction[2])
	gl.Uniform3f(s.uniformLightColorLoc, light.Color[0], light.Color[1], light.Color[2])
	gl.Uniform1f(s.uniformLightDiffuseIntensityLoc, light.DiffuseIntensity)
}

func NewDirLightShader(vertex, frag string) (*DirLightShader, error) {
	c, err := NewShader(vertex, frag)
	if err != nil {
		return nil, err
	}
	s := &DirLightShader{
		Shader:                          c,
		uniformPosLoc:                   uniformLocation(c, "gPosition"),
		uniformNormalLoc:                uniformLocation(c, "gNormal"),
		uniformAlbedoSpecLoc:            uniformLocation(c, "gAlbedoSpec"),
		uniformLightDirectionLoc:        uniformLocation(c, "dirLight.Direction"),
		uniformLightColorLoc:            uniformLocation(c, "dirLight.Color"),
		uniformLightDiffuseIntensityLoc: uniformLocation(c, "dirLight.DiffuseIntensity"),
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
