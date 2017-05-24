package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewGBufferPipeline() *GBufferPipeline {
	p := &GBufferPipeline{
		buffer:     NewGbuffer(),
		nullShader: NewDefaultShader("null", "null"),
	}
	p.mShader = NewMaterialShader()

	p.tShader = NewTextureShader()
	blockIndex := gl.GetUniformBlockIndex(p.tShader.Program(), gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(p.tShader.Program(), blockIndex, 0)
	return p
}

type GBufferPipeline struct {
	buffer     *Gbuffer
	tShader    TextureShader
	mShader    MaterialShader
	nullShader *DefaultShader
}

type Gbuffer struct {
	fbo uint32

	gNormal      uint32
	gAlbedoSpec  uint32
	gDepth       uint32
	finalTexture uint32
}

func NewGbuffer() *Gbuffer {
	gbuffer := &Gbuffer{}
	gl.GenFramebuffers(1, &gbuffer.fbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gbuffer.fbo)

	// Normal + roughness texture buffer
	gl.GenTextures(1, &gbuffer.gNormal)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.gNormal)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, gbuffer.gNormal, 0)

	// Albedo + metallic texture buffer
	gl.GenTextures(1, &gbuffer.gAlbedoSpec)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.gAlbedoSpec)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gl.TEXTURE_2D, gbuffer.gAlbedoSpec, 0)

	//  Depth texture
	gl.GenTextures(1, &gbuffer.gDepth)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.gDepth)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT16, windowWidth, windowHeight, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, gbuffer.gDepth, 0)

	// Final output texture for this FBO
	gl.GenTextures(1, &gbuffer.finalTexture)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.finalTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, windowWidth, windowHeight, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT3, gl.TEXTURE_2D, gbuffer.finalTexture, 0)

	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, status: 0x%x\n", status))
	}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return gbuffer
}
