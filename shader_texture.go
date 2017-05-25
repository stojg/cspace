package main

type ModelShader interface {
	ModelUniform() int32
}

func NewTextureShader() *GbufferTShader {
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
	switch t {
	case Albedo:
		return s.LocAlbedo
	case Metallic:
		return s.LocMetallic
	case Roughness:
		return s.LocRoughness
	case Normal:
		return s.LocNormal
	default:
		return -1
	}
}
