package main

func NewGBufferPipeline() *GBufferPipeline {
	p := &GBufferPipeline{
		buffer:     NewGbuffer(),
		nullShader: NewDefaultShader("null", "null"),
	}
	p.mShader = &GbufferMShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_m"),
	}
	p.mShader.uniformModelLoc = uniformLocation(p.mShader.Shader, "model")
	p.mShader.locDiffuse = uniformLocation(p.mShader.Shader, "mat.diffuse")
	p.mShader.locSpecularExp = uniformLocation(p.mShader.Shader, "mat.specularExp")

	p.tShader = &GbufferTShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_t"),
	}
	p.tShader.uniformDiffuseLoc = uniformLocation(p.tShader.Shader, "mat.diffuse0")
	p.tShader.uniformSpecularLoc = uniformLocation(p.tShader.Shader, "mat.specular0")
	p.tShader.uniformNormalLoc = uniformLocation(p.tShader.Shader, "mat.normal0")
	p.tShader.uniformModelLoc = uniformLocation(p.tShader.Shader, "model")
	return p
}

type GBufferPipeline struct {
	buffer     *Gbuffer
	tShader    *GbufferTShader
	mShader    *GbufferMShader
	nullShader *DefaultShader
}
