package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader interface {
	Use()
	UsePV(projection, view mgl32.Mat4)
	Program() uint32
}

type TextureShader interface {
	Shader
	ModelUniform() int32
	TextureUniform(TextureType, int) int32
}

type MaterialShader interface {
	Shader
	ModelUniform() int32
	DiffuseUniform() int32
	SpecularExpUniform() int32
}

type GbufferMShader struct {
	Shader
	uniformModelLoc int32
	locDiffuse      int32
	locSpecularExp  int32
}

func (s *GbufferMShader) ModelUniform() int32 {
	return s.uniformModelLoc
}

func (s *GbufferMShader) DiffuseUniform() int32 {
	return s.locDiffuse
}
func (s *GbufferMShader) SpecularExpUniform() int32 {
	return s.locSpecularExp
}

type GbufferLightShader interface {
	Shader
	UniformNormalLoc() int32
	UniformAlbedoSpecLoc() int32
	UniformDepthLoc() int32
}

type GbufferTShader struct {
	Shader
	uniformModelLoc    int32
	uniformDiffuseLoc  int32
	uniformNormalLoc   int32
	uniformSpecularLoc int32
}

func (s *GbufferTShader) TextureUniform(t TextureType, num int) int32 {
	if t == Diffuse {
		return s.uniformDiffuseLoc
	}
	if t == Specular {
		return s.uniformSpecularLoc
	}
	if t == Normal {
		return s.uniformNormalLoc
	}
	return -1
}

func (s *GbufferTShader) ModelUniform() int32 {
	return s.uniformModelLoc
}

type PassthroughShader struct {
	Shader
	uniformScreenTextureLoc int32
}

func NewPassthroughShader() *PassthroughShader {

	s := &PassthroughShader{
		Shader: NewDefaultShader("fx", "fx_text_pass"),
	}
	s.uniformScreenTextureLoc = uniformLocation(s.Shader, "screenTexture")
	return s

}

type PointLightShader struct {
	*DefaultShader
	uniformNormalLoc     int32
	uniformAlbedoSpecLoc int32
	uniformDepthLoc      int32

	uniformLightPosLoc              int32
	uniformLightColorLoc            int32
	uniformLightDiffuseIntensityLoc int32
	uniformLightLinearLoc           int32
	uniformLightQuadraticLoc        int32

	locProjMatrixInv int32
	locViewMatrixInv int32
}

func (s *PointLightShader) UniformNormalLoc() int32 {
	return s.uniformNormalLoc
}

func (s *PointLightShader) UniformAlbedoSpecLoc() int32 {
	return s.uniformAlbedoSpecLoc
}

func (s *PointLightShader) UniformDepthLoc() int32 {
	return s.uniformDepthLoc
}

func (s *PointLightShader) SetLight(light *PointLight) {
	gl.Uniform3f(s.uniformLightPosLoc, light.Position[0], light.Position[1], light.Position[2])
	gl.Uniform3f(s.uniformLightColorLoc, light.Color[0], light.Color[1], light.Color[2])
	gl.Uniform1f(s.uniformLightLinearLoc, light.Linear)
	gl.Uniform1f(s.uniformLightQuadraticLoc, light.Exp)
}

func NewPointLightShader(vertex, frag string) *PointLightShader {
	c := NewDefaultShader(vertex, frag)
	s := &PointLightShader{
		DefaultShader:        c,
		uniformNormalLoc:     uniformLocation(c, "gNormal"),
		uniformAlbedoSpecLoc: uniformLocation(c, "gAlbedoSpec"),
		uniformDepthLoc:      uniformLocation(c, "gDepth"),

		uniformLightPosLoc:       uniformLocation(c, "pointLight.Position"),
		uniformLightColorLoc:     uniformLocation(c, "pointLight.Color"),
		uniformLightLinearLoc:    uniformLocation(c, "pointLight.Linear"),
		uniformLightQuadraticLoc: uniformLocation(c, "pointLight.Quadratic"),

		locProjMatrixInv: uniformLocation(c, "projMatrixInv"),
		locViewMatrixInv: uniformLocation(c, "viewMatrixInv"),
	}
	return s
}

type DirLightShader struct {
	*DefaultShader
	uniformDepthLoc      int32
	uniformNormalLoc     int32
	uniformAlbedoSpecLoc int32

	uniformLightDirectionLoc        int32
	uniformLightColorLoc            int32
	uniformLightDiffuseIntensityLoc int32

	locProjMatrixInv int32
	locViewMatrixInv int32
}

func (s *DirLightShader) UniformNormalLoc() int32 {
	return s.uniformNormalLoc
}

func (s *DirLightShader) UniformAlbedoSpecLoc() int32 {
	return s.uniformAlbedoSpecLoc
}

func (s *DirLightShader) UniformDepthLoc() int32 {
	return s.uniformDepthLoc
}

func (s *DirLightShader) SetLight(light *DirectionalLight) {
	gl.Uniform3f(s.uniformLightDirectionLoc, light.Direction[0], light.Direction[1], light.Direction[2])
	gl.Uniform3f(s.uniformLightColorLoc, light.Color[0], light.Color[1], light.Color[2])
}

func NewDirLightShader(vertex, frag string) *DirLightShader {
	c := NewDefaultShader(vertex, frag)
	s := &DirLightShader{
		DefaultShader:            c,
		uniformDepthLoc:          uniformLocation(c, "gDepth"),
		uniformNormalLoc:         uniformLocation(c, "gNormal"),
		uniformAlbedoSpecLoc:     uniformLocation(c, "gAlbedoSpec"),
		uniformLightDirectionLoc: uniformLocation(c, "dirLight.Direction"),
		uniformLightColorLoc:     uniformLocation(c, "dirLight.Color"),

		locProjMatrixInv: uniformLocation(c, "projMatrixInv"),
		locViewMatrixInv: uniformLocation(c, "viewMatrixInv"),
	}
	return s
}

type DefaultShader struct {
	program    uint32
	projection int32
	view       int32
}

func (s *DefaultShader) Program() uint32 {
	return s.program
}

func (s *DefaultShader) Use() {
	gl.UseProgram(s.program)
}

func (s *DefaultShader) UsePV(projection, view mgl32.Mat4) {
	gl.UseProgram(s.program)
	gl.UniformMatrix4fv(s.projection, 1, false, &projection[0])
	gl.UniformMatrix4fv(s.view, 1, false, &view[0])
}

func NewDefaultShader(vertex, frag string) *DefaultShader {
	shader := &DefaultShader{}
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

	shader.program = program

	// try to find these pretty standard uniforms
	shader.view = gl.GetUniformLocation(shader.Program(), gl.Str("view\x00"))
	shader.projection = gl.GetUniformLocation(shader.Program(), gl.Str("projection\x00"))

	glLogShader(shader, vertex, frag)

	return shader
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
