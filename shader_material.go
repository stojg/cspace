package main

import "github.com/go-gl/gl/v4.1-core/gl"

func NewMaterialShader() *GbufferMShader {
	shader := &GbufferMShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_m"),
	}

	blockIndex := gl.GetUniformBlockIndex(shader.Program(), gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(shader.Program(), blockIndex, 0)

	shader.locModel = uniformLocation(shader, "model")
	shader.locAlbedo = uniformLocation(shader, "mat.albedo")
	shader.locMetallic = uniformLocation(shader, "mat.metallic")
	shader.locRoughness = uniformLocation(shader, "mat.roughness")
	return shader
}

type GbufferMShader struct {
	Shader
	locModel     int32
	locAlbedo    int32
	locMetallic  int32
	locRoughness int32
}
