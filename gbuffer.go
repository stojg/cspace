package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Gbuffer struct {
	fbo uint32

	gPosition    uint32
	gNormal      uint32
	gAlbedoSpec  uint32
	gDepth       uint32
	finalTexture uint32
}

func NewGbuffer(SCR_WIDTH, SCR_HEIGHT int32) *Gbuffer {
	gbuffer := &Gbuffer{}
	gl.GenFramebuffers(1, &gbuffer.fbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gbuffer.fbo)

	// position
	gl.GenTextures(1, &gbuffer.gPosition)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.gPosition)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, SCR_WIDTH, SCR_HEIGHT, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, gbuffer.gPosition, 0)

	// - Normal color buffer
	gl.GenTextures(1, &gbuffer.gNormal)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.gNormal)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, SCR_WIDTH, SCR_HEIGHT, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gl.TEXTURE_2D, gbuffer.gNormal, 0)

	// - Color + Specular color buffer
	gl.GenTextures(1, &gbuffer.gAlbedoSpec)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.gAlbedoSpec)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, SCR_WIDTH, SCR_HEIGHT, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT2, gl.TEXTURE_2D, gbuffer.gAlbedoSpec, 0)

	// gDepth
	gl.GenTextures(1, &gbuffer.gDepth)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.gDepth)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH32F_STENCIL8, SCR_WIDTH, SCR_HEIGHT, 0, gl.DEPTH_STENCIL, gl.FLOAT_32_UNSIGNED_INT_24_8_REV, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.TEXTURE_2D, gbuffer.gDepth, 0)

	// final
	gl.GenTextures(1, &gbuffer.finalTexture)
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.finalTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, SCR_WIDTH, SCR_HEIGHT, 0, gl.RGB, gl.FLOAT, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT4, gl.TEXTURE_2D, gbuffer.finalTexture, 0)

	// - Finally check if framebuffer is complete
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, status: 0x%x\n", status))
	}

	// restore default FBO
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return gbuffer
}

func (g *Gbuffer) StartFrame() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, g.fbo)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT4)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (g *Gbuffer) BindForGeomPass() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, g.fbo)
	var attachments = [3]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}
	gl.DrawBuffers(3, &attachments[0])
}

func (g *Gbuffer) BindForStencilPass() {
	gl.DrawBuffer(gl.NONE)
}

func (g *Gbuffer) BindForLightPass(s GbufferLightShader) {
	gl.DrawBuffer(gl.COLOR_ATTACHMENT4)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(s.UniformPosLoc(), 0)
	gl.BindTexture(gl.TEXTURE_2D, g.gPosition)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.Uniform1i(s.UniformNormalLoc(), 1)
	gl.BindTexture(gl.TEXTURE_2D, g.gNormal)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.Uniform1i(s.UniformAlbedoSpecLoc(), 2)
	gl.BindTexture(gl.TEXTURE_2D, g.gAlbedoSpec)
}

func (g *Gbuffer) BindForFinalPass() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, g.fbo)
	gl.ReadBuffer(gl.COLOR_ATTACHMENT4)
}
