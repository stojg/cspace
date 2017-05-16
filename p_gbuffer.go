package main

func NewGBufferPipeline() *GBufferPipeline {
	p := &GBufferPipeline{
		buffer:     NewGbuffer(),
		nullShader: NewDefaultShader("null", "null"),
	}
	p.mShader = NewMaterialShader()

	p.tShader = NewTextureShader()
	return p
}

type GBufferPipeline struct {
	buffer     *Gbuffer
	tShader    TextureShader
	mShader    MaterialShader
	nullShader *DefaultShader
}
