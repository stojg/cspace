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

	glLog("GL Context Params:\n")

	for i := 0; i < 10; i++ {
		var v int32
		gl.GetIntegerv(params[i], &v)
		glLog(fmt.Sprintf("%s %d", names[i], v))
	}

	var v int32
	gl.GetIntegerv(params[10], &v)
	glLog(fmt.Sprintf("%s %d", names[10], v))

	var b bool
	gl.GetBooleanv(params[11], &b)
	glLog(fmt.Sprintf("%s %t", names[11], b))

	//gl_log("%s %i %i\n", names[10], v[0], v[1]);
	//unsigned char s = 0;
	//glGetBooleanv(params[11], &s);
	//gl_log(m"%s %u\n", names[11], (unsigned int)s);
	//gl_log("-----------------------------\n");

}
