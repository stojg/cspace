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
	p.shader = &GbufferShader{
		Shader: NewDefaultShader("g_buffer", "g_buffer"),
	}
	p.shader.uniformDiffuseLoc = uniformLocation(p.shader.Shader, "mat.diffuse0")
	p.shader.uniformSpecularLoc = uniformLocation(p.shader.Shader, "mat.specular0")
	p.shader.uniformNormalLoc = uniformLocation(p.shader.Shader, "mat.normal0")
	p.shader.uniformModelLoc = uniformLocation(p.shader.Shader, "model")
	return p
}

type GBufferPipeline struct {
	buffer     *Gbuffer
	shader     *GbufferShader
	nullShader *DefaultShader
}

func (p *GBufferPipeline) Render(projection, view mgl32.Mat4, graph *Node) {

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, p.buffer.fbo)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT4)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	p.shader.UsePV(projection, view)

	// 1. render into the gBuffer
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, p.buffer.fbo)
	var attachments = [3]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}
	gl.DrawBuffers(3, &attachments[0])

	// Only the geometry pass updates the gDepth buffer
	gl.DepthMask(true)
	gl.ClearColor(0.0, 0.0, 0.0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	graph.Render(p.shader)

}
