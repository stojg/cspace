package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type GbufferTextureType int32

const (
	GBUFFER_TEXTURE_TYPE_POSITION GbufferTextureType = iota
	GBUFFER_TEXTURE_TYPE_DIFFUSE
	GBUFFER_TEXTURE_TYPE_NORMAL
	GBUFFER_TEXTURE_TYPE_TEXCOORD
	GBUFFER_NUM_TEXTURES
)

type Gbuffer struct {
	fbo          uint32
	textures     [GBUFFER_NUM_TEXTURES]uint32
	depthTexture uint32
}

func NewGbuffer(WindowWidth, WindowHeight int32) *Gbuffer {
	_ = mgl32.Scale3D(1, 1, 1)
	gbuffer := &Gbuffer{}
	gl.GenFramebuffers(1, &gbuffer.fbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, gbuffer.fbo)

	// Create the gbuffer textures
	gl.GenTextures(int32(GBUFFER_NUM_TEXTURES), &gbuffer.textures[0])
	gl.GenTextures(1, &gbuffer.depthTexture)

	for i := 0; i < len(gbuffer.textures); i++ {
		gl.BindTexture(gl.TEXTURE_2D, gbuffer.textures[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB32F, WindowWidth, WindowHeight, 0, gl.RGB, gl.FLOAT, nil)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+uint32(i), gl.TEXTURE_2D, gbuffer.textures[i], 0)
	}

	// depth
	gl.BindTexture(gl.TEXTURE_2D, gbuffer.depthTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32F, WindowWidth, WindowHeight, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, gbuffer.depthTexture, 0)

	DrawBuffers := []uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2, gl.COLOR_ATTACHMENT3}
	gl.DrawBuffers(int32(len(DrawBuffers)), &DrawBuffers[0])

	Status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	if Status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FB error, status: 0x%x\n", Status))
	}

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

func (g *Gbuffer) SetReadBuffer(textureType GbufferTextureType) {
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + uint32(textureType))
}
