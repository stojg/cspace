package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func NewGBufferPipeline() *GBufferPipeline {
	p := &GBufferPipeline{
		buffer:     NewGbuffer(windowWidth, windowHeight),
		nullShader: NewDefaultShader("null", "null"),
	}
	p.mShader = &GbufferMShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer_m"),
	}
	p.mShader.uniformModelLoc = uniformLocation(p.mShader.Shader, "model")
	p.mShader.locDiffuse = uniformLocation(p.mShader.Shader, "mat.diffuse")
	p.mShader.locSpecularExp = uniformLocation(p.mShader.Shader, "mat.specularExp")

	p.tShader = &GbufferTShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer"),
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

func (p *GBufferPipeline) Render(projection, view mgl32.Mat4, graph *Node) {

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, p.buffer.fbo)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT4)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// 1. render into the gBuffer
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, p.buffer.fbo)
	var attachments = [3]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}
	gl.DrawBuffers(3, &attachments[0])

	// Only the geometry pass updates the gDepth buffer
	gl.DepthMask(true)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	graph.Render(projection, view, p.tShader, p.mShader)

}
