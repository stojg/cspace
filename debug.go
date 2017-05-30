package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var albedoDebugShader *DefaultShader
var vaoAlbedodoDebug uint32

func DisplayAlbedoTexBuffer(textureID uint32) {
	if vaoAlbedodoDebug == 0 {
		quadVertices := []float32{
			0.0, 1, 0.0, 0.0, 1.0,
			0.0, 0.5, 0.0, 0.0, 0.0,
			0.5, 1, 0.0, 1.0, 1.0,
			0.5, 0.5, 0.0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoAlbedodoDebug)
		gl.BindVertexArray(vaoAlbedodoDebug)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		albedoDebugShader = NewDefaultShader("fbo_debug", "fbo_debug")
	}

	albedoDebugShader.Use()

	GLBindTexture(0, 0, textureID)

	gl.BindVertexArray(vaoAlbedodoDebug)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

var vaoDebugNormalTextureRect uint32

func DisplayNormalBufferTexture(textureID uint32) {
	if vaoDebugNormalTextureRect == 0 {
		quadVertices := []float32{
			0.5, 1, 0.0, 0.0, 1.0,
			0.5, 0.5, 0.0, 0.0, 0.0,
			1, 1, 0.0, 1.0, 1.0,
			1, 0.5, 0.0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoDebugNormalTextureRect)
		gl.BindVertexArray(vaoDebugNormalTextureRect)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		albedoDebugShader = NewDefaultShader("fbo_debug", "fbo_debug")
	}

	gl.ActiveTexture(gl.TEXTURE0)
	albedoDebugShader.Use()
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoDebugNormalTextureRect)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

	gl.UseProgram(0)
}

var vaoDebugDepthTexturedRect uint32
var depthShaderTextureLoc int32
var depthShader *DefaultShader

func DisplayDepthbufferTexture(textureID uint32) {
	if vaoDebugDepthTexturedRect == 0 {
		quadVertices := []float32{
			-0.5, 1, 0.0, 0.0, 1.0,
			-0.5, 0.5, 0.0, 0.0, 0.0,
			0, 1, 0.0, 1.0, 1.0,
			0, 0.5, 0.0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoDebugDepthTexturedRect)
		gl.BindVertexArray(vaoDebugDepthTexturedRect)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		depthShader = NewDefaultShader("depth_debug", "depth_debug")
		depthShaderTextureLoc = uniformLocation(depthShader, "screenTexture")
	}

	depthShader.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(depthShaderTextureLoc, 0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoDebugDepthTexturedRect)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

	gl.UseProgram(0)
}

var ssaoDebugTextureLoc int32
var ssaoDebug *DefaultShader
var vaoSsaoTexturedRect uint32

func DisplaySsaoTexture(textureID uint32) {
	if vaoSsaoTexturedRect == 0 {
		quadVertices := []float32{
			-1, 1, 0.0, 0.0, 1.0,
			-1, 0.5, 0.0, 0.0, 0.0,
			-0.5, 1, 0.0, 1.0, 1.0,
			-0.5, 0.5, 0.0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoSsaoTexturedRect)
		gl.BindVertexArray(vaoSsaoTexturedRect)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		ssaoDebug = NewDefaultShader("depth_debug", "depth_debug")
		ssaoDebugTextureLoc = uniformLocation(ssaoDebug, "screenTexture")
	}

	ssaoDebug.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(ssaoDebugTextureLoc, 0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoSsaoTexturedRect)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

var vaoRoughnessTexturedRect uint32
var roughnessDebug *DefaultShader
var roughnessDebugTextureLoc int32

func DisplayRoughnessTexture(textureID uint32) {
	if vaoRoughnessTexturedRect == 0 {
		quadVertices := []float32{
			-1.0, 0.5, 0.0, 0.0, 1.0,
			-1.0, 0.0, 0.0, 0.0, 0.0,
			-0.5, 0.5, 0, 1.0, 1.0,
			-0.5, 0.0, 0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoRoughnessTexturedRect)
		gl.BindVertexArray(vaoRoughnessTexturedRect)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		roughnessDebug = NewDefaultShader("depth_debug", "fbo_debug_alpha")
		roughnessDebugTextureLoc = uniformLocation(roughnessDebug, "screenTexture")
	}

	roughnessDebug.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(roughnessDebugTextureLoc, 0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoRoughnessTexturedRect)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

var vaoMetallicTexturedRect uint32
var metallicDebug *DefaultShader
var metallicDebugTextureLoc int32

func DisplayMetallicTexture(textureID uint32) {
	if vaoMetallicTexturedRect == 0 {
		quadVertices := []float32{
			-1.0, 0.0, 0.0, 0.0, 1.0,
			-1.0, -0.5, 0.0, 0.0, 0.0,
			-0.5, 0.0, 0, 1.0, 1.0,
			-0.5, -0.5, 0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoMetallicTexturedRect)
		gl.BindVertexArray(vaoMetallicTexturedRect)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		metallicDebug = NewDefaultShader("depth_debug", "fbo_debug_alpha")
		metallicDebugTextureLoc = uniformLocation(metallicDebug, "screenTexture")
	}

	metallicDebug.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(metallicDebugTextureLoc, 0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoMetallicTexturedRect)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

var vaoshadowTexturedRect uint32
var shadowDebug *DefaultShader
var shadowDebugTextureLoc int32

func DisplayShadowTexture(textureID uint32) {
	if vaoshadowTexturedRect == 0 {
		quadVertices := []float32{
			0.5, 0.5, 0.0, 0.0, 1.0,
			0.5, 0.0, 0.0, 0.0, 0.0,
			1.0, 0.5, 0, 1.0, 1.0,
			1.0, 0.0, 0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoshadowTexturedRect)
		gl.BindVertexArray(vaoshadowTexturedRect)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		shadowDebug = NewDefaultShader("depth_debug", "depth_debug")
		shadowDebugTextureLoc = uniformLocation(shadowDebug, "screenTexture")
	}

	shadowDebug.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(shadowDebugTextureLoc, 0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoshadowTexturedRect)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

var vaoBloomTexturedRect uint32
var BloomDebug *DefaultShader
var BloomDebugTextureLoc int32

func DisplayBloomTexture(textureID uint32) {
	if vaoBloomTexturedRect == 0 {
		quadVertices := []float32{
			0.5, 0.0, 0.0, 0.0, 1.0,
			0.5, -0.5, 0.0, 0.0, 0.0,
			1.0, 0.0, 0, 1.0, 1.0,
			1.0, -0.5, 0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &vaoBloomTexturedRect)
		gl.BindVertexArray(vaoBloomTexturedRect)
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		BloomDebug = NewDefaultShader("fbo_debug", "fbo_debug")
		BloomDebugTextureLoc = uniformLocation(BloomDebug, "screenTexture")
	}

	BloomDebug.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(BloomDebugTextureLoc, 0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.BindVertexArray(vaoBloomTexturedRect)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func chkError(name string) {
	err := gl.GetError()
	if err == 0 {
		return
	}
	switch err {
	case gl.INVALID_OPERATION:
		fmt.Printf("GL Error: INVALID_OPERATION 0x0%x\n", err)
	case gl.INVALID_ENUM:
		fmt.Printf("GL Error: INVALID_ENUM 0x0%x\n", err)
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		fmt.Printf("GL Error: INVALID_FRAMEBUFFER_OPERATION 0x0%x\n", err)
	default:
		fmt.Printf("GL Error: 0x0%x\n", err)
	}
	panic(name)
}

func chkFramebuffer() {
	if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
		switch s {
		case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
			panic("Framebuffer incomplete: No image is attached to FBO")
		case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
			panic("Framebuffer incomplete: Attachment is NOT complete")
		case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
			panic("Framebuffer incomplete: Draw buffer")
		case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
			panic("Framebuffer incomplete: Read buffer")
		default:
			panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, s: 0x%x\n", s))
		}
	}
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
	if elapsedSeconds > 4 {
		fpsPrevSeconds = currentSeconds
		fps := (float64(fpsFrameCount) / elapsedSeconds)
		ms := 1000 / fps
		msg := fmt.Sprintf("cspace @ %.2fms / %.0f", ms, fps)
		window.SetTitle(msg)
		fmt.Println(msg)
		fpsFrameCount = 0
	}
	fpsFrameCount++
}

func glLogShader(program uint32, vertex, frag string) {

	glLogf("------- info tShader programme %d | %s / %s -------\n", program, vertex, frag)

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

	//gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &params)
	//glLogf("%d gl.ACTIVE_UNIFORMS\n", params)
	//for i := int32(0); i < params; i++ {
	//	var actualLength int32
	//	var size int32
	//	var xtype uint32
	//	var maxLength int32 = 64
	//	name := make([]byte, maxLength)
	//	gl.GetActiveUniform(program, uint32(i), maxLength, &actualLength, &size, &xtype, &name[0])
	//	if size > 1 {
	//		for j := int32(0); j < size; j++ {
	//			longName := []byte(fmt.Sprintf("%s[%d]", name, j))
	//			location := gl.GetAttribLocation(program, &longName[0])
	//			glLogf("\t%d) %s ", i, glTypeToString(xtype), bytes.Trim(longName, "\x00"), location)
	//		}
	//	} else {
	//		location := pUniformLocation(program, fmt.Sprintf("%s\n", name))
	//		glLogf("\t%d) %s %s @ location %d\n", i, glTypeToString(xtype), bytes.Trim(name, "\x00"), location)
	//	}
	//}

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
