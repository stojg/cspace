package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type BloomFBO struct {
	fbo      uint32
	textures []uint32
}

func NewBloom() *BloomFBO {
	frameBuffer := &BloomFBO{
		textures: make([]uint32, 1),
	}
	gl.GenFramebuffers(1, &frameBuffer.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer.fbo)

	gl.GenTextures(int32(len(frameBuffer.textures)), &frameBuffer.textures[0])
	for i := 0; i < len(frameBuffer.textures); i++ {
		gl.BindTexture(gl.TEXTURE_2D, frameBuffer.textures[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, windowWidth, windowHeight, 0, gl.RGB, gl.FLOAT, nil)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+uint32(i), gl.TEXTURE_2D, frameBuffer.textures[i], 0)
	}

	if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, s: 0x%x\n", s))
	}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return frameBuffer
}

func NewBloomEffect(width, height int32) *BloomEffect {
	b := &BloomEffect{
		width:            width,
		height:           height,
		bloomFbo:         NewBloom(),
		pingBuffers:      [2]*FBO{NewFBO(width, height), NewFBO(width, height)},
		separationShader: NewDefaultShader("fx", "fx_brigthness_sep"),
		blendShader:      NewDefaultShader("fx", "fx_bloom_blender"),
		gaussianBlender:  NewDefaultShader("fx", "fx_guassian_blur"),
	}

	b.locSeparationScreenTexture = uniformLocation(b.separationShader, "screenTexture")
	b.locGaussianScreenTexture = uniformLocation(b.gaussianBlender, "screenTexture")
	b.locGaussianHorizontal = uniformLocation(b.gaussianBlender, "horizontal")
	return b
}

type BloomEffect struct {
	width, height int32
	bloomFbo      *BloomFBO
	pingBuffers   [2]*FBO

	separationShader *DefaultShader
	gaussianBlender  *DefaultShader
	blendShader      *DefaultShader

	locSeparationScreenTexture int32
	locGaussianScreenTexture   int32
	locGaussianHorizontal      int32

	quadVAO uint32
}

func (b *BloomEffect) Render(inTexture uint32) uint32 {

	if b.quadVAO == 0 {
		quadVertices := []float32{
			-1, 1, 0.0, 0.0, 1.0,
			-1, -1, 0.0, 0.0, 0.0,
			1, 1, 0.0, 1.0, 1.0,
			1, -1, 0.0, 1.0, 0.0,
		}
		var vbo uint32
		gl.GenVertexArrays(1, &b.quadVAO)
		gl.BindVertexArray(b.quadVAO)
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	}

	// separate the brightest colours into a separate texture
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.bloomFbo.fbo)

	var attachments = [1]uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &attachments[0])
	b.separationShader.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(b.locSeparationScreenTexture, 0)
	gl.BindTexture(gl.TEXTURE_2D, inTexture)
	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

	// blur the bright part
	const blurAmount = 4
	horizontal := 0
	firstIteration := true

	// ping-pong
	gl.Viewport(0, 0, b.width, b.height)
	b.gaussianBlender.Use()
	for i := 0; i < blurAmount; i++ {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, b.pingBuffers[horizontal].id)

		gl.Uniform1i(b.locGaussianHorizontal, int32(horizontal))
		if horizontal == 0 {
			horizontal = 1
		} else {
			horizontal = 0
		}
		gl.Uniform1i(b.locGaussianScreenTexture, 0)
		gl.ActiveTexture(gl.TEXTURE0)
		if firstIteration {
			gl.BindTexture(gl.TEXTURE_2D, b.bloomFbo.textures[0])
			firstIteration = false
		} else {
			gl.BindTexture(gl.TEXTURE_2D, b.pingBuffers[horizontal].textures[0])
		}
		gl.BindVertexArray(quadVAO)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		gl.BindVertexArray(0)
	}
	gl.Viewport(0, 0, windowWidth, windowHeight)

	// combine the normal and blurry bright texture for a bloom effect
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.bloomFbo.fbo)
	b.blendShader.Use()
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(uniformLocation(b.blendShader, "screenTexture"), 0)
	gl.BindTexture(gl.TEXTURE_2D, inTexture)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.Uniform1i(uniformLocation(b.blendShader, "bloomTexture"), 1)
	gl.BindTexture(gl.TEXTURE_2D, b.pingBuffers[1].textures[0])
	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

	return b.bloomFbo.textures[0]
}
