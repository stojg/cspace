package main

type MaterialShader interface {
	Shader
	ModelUniform() int32
	AlbedoUniform() int32
	MetallicUniform() int32
	RoughnessUniform() int32
}

func NewMaterialShader() MaterialShader {
	shader := &GbufferMShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_m"),
	}
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

func (s *GbufferMShader) ModelUniform() int32 {
	return s.locModel
}

func (s *GbufferMShader) AlbedoUniform() int32 {
	return s.locAlbedo
}

func (s *GbufferMShader) MetallicUniform() int32 {
	return s.locMetallic
}

func (s *GbufferMShader) RoughnessUniform() int32 {
	return s.locRoughness
}
