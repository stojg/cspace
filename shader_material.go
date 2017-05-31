package main

import "github.com/go-gl/gl/v4.1-core/gl"

func NewMaterialShader() *GbufferMShader {
	shader := &GbufferMShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_m"),
	}

	blockIndex := gl.GetUniformBlockIndex(shader.Program(), gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(shader.Program(), blockIndex, 0)

	shader.LocModel = uniformLocation(shader, "model")
	shader.LocAlbedo = uniformLocation(shader, "mat.albedo")
	shader.LocMetallic = uniformLocation(shader, "mat.metallic")
	shader.LocRoughness = uniformLocation(shader, "mat.roughness")
	return shader
}

type GbufferMShader struct {
	Shader
	LocModel     int32
	LocAlbedo    int32
	LocMetallic  int32
	LocRoughness int32
}
