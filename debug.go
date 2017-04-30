package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

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

func glLogProgramme(program uint32) {
	glLogf("------- info shader programme %d -------\n", program)

	var params int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &params)
	glLogf("GL_LINK_STATUS = %d\n", params)

	gl.GetProgramiv(program, gl.ATTACHED_SHADERS, &params)
	glLogf("GL_ATTACHED_SHADERS = %d\n", params)

	gl.GetProgramiv(program, gl.ACTIVE_ATTRIBUTES, &params)
	glLogf("GL_ACTIVE_ATTRIBUTES = %d\n", params)

	for i := int32(0); i < params; i++ {
		var maxLength int32 = 64
		var actualLength int32
		var size int32
		var xtype uint32
		var name [64]byte

		gl.GetActiveAttrib(program, uint32(i), maxLength, &actualLength, &size, &xtype, &name[0])
		if size > 1 {
			for j := int32(0); j < size; j++ {
				longName := []byte(fmt.Sprintf("%s[%d]", name, j))
				location := gl.GetAttribLocation(program, &longName[0])
				glLogf("    %d) type: %s ", i, glTypeToString(xtype), longName, location)
			}
		} else {
			location := gl.GetAttribLocation(program, &name[0])
			glLogf("    %d) type: %s name: %s, location: %d\n", i, glTypeToString(xtype), name, location)
		}
	}

	gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &params)
	glLogf("GL_ACTIVE_UNIFORMS = %d\n", params)
	for i := int32(0); i < params; i++ {
		var maxLength int32 = 64
		var actualLength int32
		var size int32
		var xtype uint32
		var name [64]byte
		gl.GetActiveUniform(program, uint32(i), maxLength, &actualLength, &size, &xtype, &name[0])
		if size > 1 {
			for j := int32(0); j < size; j++ {
				longName := []byte(fmt.Sprintf("%s[%d]", name, j))
				location := gl.GetUniformLocation(program, &longName[0])
				glLogf("    %d) type: %s ", i, glTypeToString(xtype), longName, location)
			}
		} else {
			location := gl.GetUniformLocation(program, &name[0])
			glLogf("    %d) type: %s name: %s, location: %d\n", i, glTypeToString(xtype), name, location)
		}
	}

	glLogln("---------------------------------------")
}

func glTypeToString(xtype uint32) string {
	switch xtype {
	case gl.BOOL:
		return "bool"
	case gl.INT:
		return "int"
	case gl.FLOAT:
		return "float"
	case gl.FLOAT_VEC2:
		return "vec2"
	case gl.FLOAT_VEC3:
		return "vec3"
	case gl.FLOAT_VEC4:
		return "vec4"
	case gl.FLOAT_MAT2:
		return "mat2"
	case gl.FLOAT_MAT3:
		return "mat3"
	case gl.FLOAT_MAT4:
		return "mat4"
	case gl.SAMPLER_2D:
		return "sampler2d"
	case gl.SAMPLER_3D:
		return "sampler3d"
	case gl.SAMPLER_CUBE:
		return "samplerCube"
	case gl.SAMPLER_2D_SHADOW:
		return "sampler2DShadow"
	default:
		return "unknown"
	}
}

func glLogGLParams() {

	params := []uint32{
		gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS,
		gl.MAX_CUBE_MAP_TEXTURE_SIZE,
		gl.MAX_DRAW_BUFFERS,
		gl.MAX_FRAGMENT_UNIFORM_COMPONENTS,
		gl.MAX_TEXTURE_IMAGE_UNITS,
		gl.MAX_TEXTURE_SIZE,
		gl.MAX_VARYING_FLOATS,
		gl.MAX_VERTEX_ATTRIBS,
		gl.MAX_VERTEX_TEXTURE_IMAGE_UNITS,
		gl.MAX_VERTEX_UNIFORM_COMPONENTS,
		gl.MAX_VIEWPORT_DIMS,
		gl.STEREO,
	}

	names := []string{
		"GL_MAX_COMBINED_TEXTURE_IMAGE_UNITS",
		"GL_MAX_CUBE_MAP_TEXTURE_SIZE",
		"GL_MAX_DRAW_BUFFERS",
		"GL_MAX_FRAGMENT_UNIFORM_COMPONENTS",
		"GL_MAX_TEXTURE_IMAGE_UNITS",
		"GL_MAX_TEXTURE_SIZE",
		"GL_MAX_VARYING_FLOATS",
		"GL_MAX_VERTEX_ATTRIBS",
		"GL_MAX_VERTEX_TEXTURE_IMAGE_UNITS",
		"GL_MAX_VERTEX_UNIFORM_COMPONENTS",
		"GL_MAX_VIEWPORT_DIMS",
		"GL_STEREO",
	}

	glLogln("GL Context Params:\n")

	for i := 0; i < 10; i++ {
		var v int32
		gl.GetIntegerv(params[i], &v)
		glLogln(fmt.Sprintf("%s %d", names[i], v))
	}

	var v int32
	gl.GetIntegerv(params[10], &v)
	glLogln(fmt.Sprintf("%s %d", names[10], v))

	var b bool
	gl.GetBooleanv(params[11], &b)
	glLogln(fmt.Sprintf("%s %t", names[11], b))
}
