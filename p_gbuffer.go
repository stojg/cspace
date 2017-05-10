package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func NewGBufferPipeline() *GBufferPipeline {
	return &GBufferPipeline{
		shader:     NewGbufferShader(),
		nullShader: NewDefaultShader("null", "null"),
	}
}

type GBufferPipeline struct {
	shader     *GbufferShader
	nullShader *DefaultShader
}

func (p *GBufferPipeline) Render(gbuffer *Gbuffer, projection, view mgl32.Mat4, graph *Node) {

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gbuffer.fbo)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT4)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	p.shader.UsePV(projection, view)

	// 1. render into the gBuffer
	{

		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gbuffer.fbo)
		var attachments = [3]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}
		gl.DrawBuffers(3, &attachments[0])

		// Only the geometry pass updates the gDepth buffer
		gl.DepthMask(true)
		gl.ClearColor(0.0, 0.0, 0.0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)

		graph.Render(p.shader)
		// When we get here the gDepth buffer is already populated and the stencil pass depends on it, but it does not write to it.
		gl.DepthMask(false)
	}

	// We need stencil to be enabled in the stencil pass to get the stencil buffer updated and we also need it in the
	// light pass because we render the light only if the stencil passes.
	gl.Enable(gl.STENCIL_TEST)

}
