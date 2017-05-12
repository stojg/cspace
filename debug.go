package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var shaderDisplayFBOOutput *DefaultShader
var vaoDebugTexturedRect uint32

func DisplayFramebufferTexture(textureID uint32) {
	if vaoDebugTexturedRect == 0 {
		quadVertices := []float32{
			0.5, 1, 0.0, 0.0, 1.0,
			0.5, 0.5, 0.0, 0.0, 0.0,
			1, 1, 0.0, 1.0, 1.0,
			1, 0.5, 0.0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoDebugTexturedRect)
		gl.GenBuffers(1, &vaoDebugTexturedRect)
		gl.BindVertexArray(vaoDebugTexturedRect)
		gl.BindBuffer(gl.ARRAY_BUFFER, vaoDebugTexturedRect)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		shaderDisplayFBOOutput = NewDefaultShader("fbo_debug", "fbo_debug")
	}

	gl.ActiveTexture(gl.TEXTURE0)
	shaderDisplayFBOOutput.Use()
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoDebugTexturedRect)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

	//gl.DrawArrays(gl.TRIANGLES, 0, 12)
	//gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func restartLog() error {
	f, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s cspace log file\n", time.Now().Format("15:04:05.000000000"))
	return err
}

func glLogln(s string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s %s\n", time.Now().Format("15:04:05.000000000"), s)
	if err != nil {
		panic(err)
	}
}

func glLogf(format string, a ...interface{}) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	now := time.Now().Format("15:04:05.000000000")
	args := append([]interface{}{now}, a...)
	_, err = fmt.Fprintf(f, "%s "+format, args...)
	if err != nil {
		panic(err)
	}
}

func glError(inError error) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s %v\n", time.Now().Format("15:04:05.000000000"), inError)
	fmt.Fprintf(os.Stderr, "%s %v\n", time.Now().Format("15:04:05.000000000"), inError)
	if err != nil {
		panic(err)
	}
}

var fpsPrevSeconds float64
var fpsFrameCount int

func fpsCounter(window *glfw.Window) {
	currentSeconds := glfw.GetTime()
	elapsedSeconds := currentSeconds - fpsPrevSeconds
	if elapsedSeconds > 0.25 {
		fpsPrevSeconds = currentSeconds
		fps := (float64(fpsFrameCount) / elapsedSeconds)
		ms := 1000 / fps
		window.SetTitle(fmt.Sprintf("cspace @ %.1fFPS, %.2fms", fps, ms))
		fpsFrameCount = 0
	}
	fpsFrameCount++
}

func glLogShader(shader Shader, vertex, frag string) {

	program := shader.Program()
	glLogf("------- info shader programme %d | %s / %s -------\n", program, vertex, frag)

	var params int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &params)
	glLogf("gl.LINK_STATUS = %d\n", params)

	gl.GetProgramiv(program, gl.ATTACHED_SHADERS, &params)
	glLogf("%d gl.ATTACHED_SHADERS\n", params)

	gl.GetProgramiv(program, gl.ACTIVE_UNIFORM_BLOCKS, &params)
	glLogf("%d gl.ACTIVE_UNIFORM_BLOCKS\n", params)
	for i := int32(0); i < params; i++ {
		var nameLength int32
		var size int32
		gl.GetActiveUniformBlockiv(program, uint32(i), gl.UNIFORM_BLOCK_NAME_LENGTH, &nameLength)
		name := make([]byte, nameLength)
		gl.GetActiveUniformBlockName(program, uint32(i), nameLength, &size, &name[0])
		glLogf("\t%d) %s\n", i, name)
	}

	gl.GetProgramiv(program, gl.ACTIVE_ATTRIBUTES, &params)
	glLogf("%d gl.ACTIVE_ATTRIBUTES\n", params)

	for i := int32(0); i < params; i++ {
		var maxLength int32 = 64
		var actualLength int32
		var size int32
		var xType uint32
		name := make([]byte, maxLength)

		gl.GetActiveAttrib(program, uint32(i), maxLength, &actualLength, &size, &xType, &name[0])
		if size > 1 {
			for j := int32(0); j < size; j++ {
				longName := []byte(fmt.Sprintf("%s[%d]", name, j))
				location := gl.GetAttribLocation(program, &longName[0])
				glLogf("\t%d) %s ", i, glTypeToString(xType), bytes.Trim(longName, "\x00"), location)
			}
		} else {
			location := gl.GetAttribLocation(program, &name[0])
			glLogf("\t%d) %s %s @ location %d\n", i, glTypeToString(xType), bytes.Trim(name, "\x00"), location)
		}
	}

	gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &params)
	glLogf("%d gl.ACTIVE_UNIFORMS\n", params)
	for i := int32(0); i < params; i++ {
		var actualLength int32
		var size int32
		var xtype uint32
		var maxLength int32 = 64
		name := make([]byte, maxLength)
		gl.GetActiveUniform(program, uint32(i), maxLength, &actualLength, &size, &xtype, &name[0])
		if size > 1 {
			for j := int32(0); j < size; j++ {
				longName := []byte(fmt.Sprintf("%s[%d]", name, j))
				location := gl.GetAttribLocation(program, &longName[0])
				glLogf("\t%d) %s ", i, glTypeToString(xtype), bytes.Trim(longName, "\x00"), location)
			}
		} else {
			location := uniformLocation(shader, fmt.Sprintf("%s\n", name))
			glLogf("\t%d) %s %s @ location %d\n", i, glTypeToString(xtype), bytes.Trim(name, "\x00"), location)
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
		"gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS",
		"gl.MAX_CUBE_MAP_TEXTURE_SIZE",
		"gl.MAX_DRAW_BUFFERS",
		"gl.MAX_FRAGMENT_UNIFORM_COMPONENTS",
		"gl.MAX_TEXTURE_IMAGE_UNITS",
		"gl.MAX_VERTEX_ATTRIBS",
	}

	glLogln("GL Context Params:\n")

	for i := 0; i < 6; i++ {
		var v int32
		gl.GetIntegerv(params[i], &v)
		glLogln(fmt.Sprintf("%s %d", names[i], v))
	}
}

//func CaptureRGBA(im *image.RGBA) {
//	b := im.Bounds()
//	gl.ReadBuffer(gl.BACK_LEFT)
//	gl.ReadPixels(0, 0, b.Dx(), b.Dy(), gl.RGBA, gl.UNSIGNED_BYTE, im.Pix)
//}
//
//// Note: You may want to call ClearAlpha(1) first..
//func CaptureToPng(filename string) {
//	w, h := GetViewportWH()
//	im := image.NewRGBA(image.Rect(0, 0, w, h))
//	CaptureRGBA(im)
//
//	fd, err := os.Create(filename)
//	if err != nil {
//		log.Panic("Err: ", err)
//	}
//	defer fd.Close()
//
//	png.Encode(fd, im)
//}
