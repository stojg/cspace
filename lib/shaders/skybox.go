package shaders

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	_ "github.com/go-gl/glfw/v3.2/glfw"
)

type Skybox struct {
	Program          uint32
	LocScreenTexture int32
	LocProjection    int32
	LocView          int32
	SkyboxVAO        uint32
}

func NewSkybox() *Skybox {
	c := buildShader("skybox", "skybox")
	shader := &Skybox{
		Program:          c,
		LocScreenTexture: loc(c, "skybox"),
		LocProjection:    loc(c, "projection"),
		LocView:          loc(c, "view"),
	}
	gl.GenVertexArrays(1, &shader.SkyboxVAO)
	gl.BindVertexArray(shader.SkyboxVAO)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(skyboxVertices), gl.Ptr(skyboxVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.BindVertexArray(0)

	return shader

}

var skyboxVertices []float32 = []float32{
	// Positions
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
}
