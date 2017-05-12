package main

import "github.com/go-gl/gl/v4.1-core/gl"

func NewBloomEffect() *BloomEffect {
	return &BloomEffect{
		bloomFbo:    NewBloom(),
		pingBuffers: [2]*FBO{NewFBO(), NewFBO()},
	}

}

type BloomEffect struct {
	bloomFbo    *BloomFBO
	pingBuffers [2]*FBO
}

func (b *BloomEffect) Render(inTexture uint32) uint32 {
	// divide the brightest colours into a new buffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.bloomFbo.fbo)

	var attachments = [2]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1}
	gl.DrawBuffers(2, &attachments[0])
	bloomColShader.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(uniformLocation(bloomColShader, "screenTexture"), 0)
	gl.BindTexture(gl.TEXTURE_2D, inTexture)
	renderQuad()

	// blur the bright part
	const blurAmount = 2
	horizontal := 0
	firstIteration := true

	// @todo cache these outside the render loop
	textLoc := uniformLocation(shaderBlur, "screenTexture")
	horisontalLoc := uniformLocation(shaderBlur, "horizontal")

	// pingpong
	for i := 0; i < blurAmount; i++ {
		shaderBlur.Use()
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, b.pingBuffers[horizontal].id)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(textLoc, 0)
		gl.Uniform1i(horisontalLoc, int32(horizontal))
		if horizontal == 0 {
			horizontal = 1
		} else {
			horizontal = 0
		}
		if firstIteration {
			gl.BindTexture(gl.TEXTURE_2D, b.bloomFbo.textures[1])
			firstIteration = false
		} else {
			gl.BindTexture(gl.TEXTURE_2D, b.pingBuffers[horizontal].textures[0])
		}
		renderQuad()
	}

	// combine the normal and blurry bright texture for a bloom effect
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.bloomFbo.fbo)
	bloomBlender.Use()
	gl.DrawBuffer(gl.COLOR_ATTACHMENT1)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(uniformLocation(bloomBlender, "screenTexture"), 0)
	gl.BindTexture(gl.TEXTURE_2D, b.bloomFbo.textures[0])
	gl.ActiveTexture(gl.TEXTURE1)
	gl.Uniform1i(uniformLocation(bloomBlender, "bloomTexture"), 1)
	gl.BindTexture(gl.TEXTURE_2D, b.pingBuffers[1].textures[0])
	renderQuad()
	return b.bloomFbo.textures[1]
}
