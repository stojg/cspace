package main

type ModelShader interface {
	ModelUniform() int32
}

type TextureShader interface {
	Shader
	ModelUniform() int32
	TextureUniform(TextureType) int32
}

func NewTextureShader() TextureShader {
	shader := &GbufferTShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_t"),
	}
	shader.LocModel = uniformLocation(shader, "model")
	shader.LocAlbedo = uniformLocation(shader, "mat.albedo")
	shader.LocMetallic = uniformLocation(shader, "mat.metallic")
	shader.LocRoughness = uniformLocation(shader, "mat.roughness")
	shader.LocNormal = uniformLocation(shader, "mat.normal")
	return shader
}

type GbufferTShader struct {
	Shader
	LocModel     int32
	LocAlbedo    int32
	LocRoughness int32
	LocMetallic  int32
	LocNormal    int32
}

func (s *GbufferTShader) TextureUniform(t TextureType) int32 {
	if t == Albedo {
		return s.LocAlbedo
	}
	if t == Metallic {
		return s.LocMetallic
	}
	if t == Roughness {
		return s.LocRoughness
	}
	if t == Normal {
		return s.LocNormal
	}
	return -1
}

func (s *GbufferTShader) ModelUniform() int32 {
	return s.LocModel
}
