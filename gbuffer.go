package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Gbuffer struct {
	fbo uint32

	gPosition   uint32
	gNormal     uint32
	gAlbedoSpec uint32
	rboDepth    uint32
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

	// - Tell OpenGL which color attachments we'll use (of this framebuffer) for rendering
	var attachments = [3]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}
	gl.DrawBuffers(3, &attachments[0])

	// - Create and attach depth buffer (renderbuffer)
	gl.GenRenderbuffers(1, &gbuffer.rboDepth)
	gl.BindRenderbuffer(gl.RENDERBUFFER, gbuffer.rboDepth)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, SCR_WIDTH, SCR_HEIGHT)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, gbuffer.rboDepth)
	// - Finally check if framebuffer is complete

	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	if status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FB error, status: 0x%x\n", status))
	}
	fmt.Println(status)

	// restore default FBO
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return gbuffer
}

func (g *Gbuffer) BindForWriting() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, g.fbo)
}

func (g *Gbuffer) BindForReading() {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, g.fbo)
}

func (g *Gbuffer) SetReadBuffer(textureType uint32) {
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + uint32(textureType))
}
