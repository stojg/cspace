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

	GLFramebuffer(&b.fbo)

	GLTextureRGB16F(&b.texture, windowWidth, windowHeight, gl.LINEAR, gl.CLAMP_TO_EDGE, nil)
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

	GLBindTexture(0, b.bloomSeparator.LocScreenTexture, inTexture)
	renderQuad()

	// blur the bright part
	const blurAmount = 10
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
		if firstIteration {
			GLBindTexture(0, b.gaussianShader.LocScreenTexture, b.texture)
			firstIteration = false
		} else {
			GLBindTexture(0, b.gaussianShader.LocScreenTexture, b.pingBuffers[horizontal].textures[0])
		}
		renderQuad()
	}
	gl.Viewport(0, 0, windowWidth, windowHeight)

	// combine the normal and blurry bright texture for a bloom effect
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.fbo)
	gl.UseProgram(b.bloomBlend.Program)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)

	GLBindTexture(0, b.bloomBlend.LocScreenTexture, inTexture)
	GLBindTexture(1, b.bloomBlend.LocBloomTexture, b.pingBuffers[1].textures[0])

	renderQuad()

	return b.texture
}
