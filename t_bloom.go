package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/cspace/lib/shaders"
)

func NewBloomEffect(width, height int32) *BloomEffect {
	b := &BloomEffect{
		width:          width,
		height:         height,
		pingBuffers:    [2]*FBO{NewFBO(width, height), NewFBO(width, height)},
		bloomSeparator: shaders.NewBloomSeparator(),
		bloomBlend:     shaders.NewBloomBlend(),
		gaussianShader: shaders.NewGaussian(),
	}

	gl.GenFramebuffers(1, &b.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.fbo)

	gl.GenTextures(1, &b.texture)

	gl.BindTexture(gl.TEXTURE_2D, b.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, windowWidth, windowHeight, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, b.texture, 0)

	chkFramebuffer()

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return b
}

type BloomEffect struct {
	fbo           uint32
	texture       uint32
	width, height int32
	pingBuffers   [2]*FBO

	bloomSeparator *shaders.BloomSeparator
	gaussianShader *shaders.Gaussian
	bloomBlend     *shaders.BloomBlend
}

func (b *BloomEffect) Render(inTexture uint32) uint32 {

	// separate the brightest colours into a separate texture
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.fbo)
	var attachments = [1]uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &attachments[0])

	gl.UseProgram(b.bloomSeparator.Program)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(b.bloomSeparator.LocScreenTexture, 0)
	gl.BindTexture(gl.TEXTURE_2D, inTexture)
	renderQuad()

	// blur the bright part
	const blurAmount = 2
	horizontal := 0
	firstIteration := true

	// ping-pong
	gl.Viewport(0, 0, b.width, b.height)
	gl.UseProgram(b.gaussianShader.Program)
	for i := 0; i < blurAmount; i++ {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, b.pingBuffers[horizontal].id)

		gl.Uniform1i(b.gaussianShader.LocHorizontal, int32(horizontal))
		if horizontal == 0 {
			horizontal = 1
		} else {
			horizontal = 0
		}
		gl.Uniform1i(b.gaussianShader.LocScreenTexture, 0)
		gl.ActiveTexture(gl.TEXTURE0)
		if firstIteration {
			gl.BindTexture(gl.TEXTURE_2D, b.texture)
			firstIteration = false
		} else {
			gl.BindTexture(gl.TEXTURE_2D, b.pingBuffers[horizontal].textures[0])
		}
		renderQuad()
	}
	gl.Viewport(0, 0, windowWidth, windowHeight)

	// combine the normal and blurry bright texture for a bloom effect
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.fbo)
	gl.UseProgram(b.bloomBlend.Program)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(b.bloomBlend.LocScreenTexture, 0)
	gl.BindTexture(gl.TEXTURE_2D, inTexture)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.Uniform1i(b.bloomBlend.LocBloomTexture, 1)
	gl.BindTexture(gl.TEXTURE_2D, b.pingBuffers[1].textures[0])

	renderQuad()

	return b.texture
}
