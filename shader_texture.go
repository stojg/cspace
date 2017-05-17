package main

type ModelShader interface {
	ModelUniform() int32
}

type TextureShader interface {
	Shader
	ModelUniform() int32
	TextureUniform(TextureType, int) int32
}

func NewTextureShader() TextureShader {
	shader := &GbufferTShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_t"),
	}
	shader.uniformDiffuseLoc = uniformLocation(shader, "mat.diffuse0")
	shader.uniformSpecularLoc = uniformLocation(shader, "mat.specular0")
	shader.uniformNormalLoc = uniformLocation(shader, "mat.normal0")
	shader.uniformModelLoc = uniformLocation(shader, "model")
	return shader
}

type GbufferTShader struct {
	Shader
	uniformModelLoc    int32
	uniformDiffuseLoc  int32
	uniformNormalLoc   int32
	uniformSpecularLoc int32
}

func (s *GbufferTShader) TextureUniform(t TextureType, num int) int32 {
	if t == Diffuse {
		return s.uniformDiffuseLoc
	}
	if t == Specular {
		return s.uniformSpecularLoc
	}
	if t == Normal {
		return s.uniformNormalLoc
	}
	return -1
}

func (s *GbufferTShader) ModelUniform() int32 {
	return s.uniformModelLoc
}
