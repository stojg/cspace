package main

import "github.com/go-gl/gl/v4.1-core/gl"

func NewBloomEffect() *BloomEffect {
	b := &BloomEffect{
		bloomFbo:         NewBloom(),
		pingBuffers:      [2]*FBO{NewFBO(), NewFBO()},
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
	bloomFbo    *BloomFBO
	pingBuffers [2]*FBO

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
		gl.GenVertexArrays(1, &b.quadVAO)
		gl.GenBuffers(1, &b.quadVAO)
		gl.BindVertexArray(b.quadVAO)
		gl.BindBuffer(gl.ARRAY_BUFFER, b.quadVAO)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	}

	// separate the brightest colours into a separate texture
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.bloomFbo.fbo)

	var attachments = [2]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1}
	gl.DrawBuffers(2, &attachments[0])
	b.separationShader.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(b.locSeparationScreenTexture, 0)

	gl.BindTexture(gl.TEXTURE_2D, inTexture)
	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

	// blur the bright part
	const blurAmount = 6
	horizontal := 0
	firstIteration := true

	// ping-pong
	for i := 0; i < blurAmount; i++ {
		b.gaussianBlender.Use()
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, b.pingBuffers[horizontal].id)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(b.locGaussianScreenTexture, 0)
		gl.Uniform1i(b.locGaussianHorizontal, int32(horizontal))
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
		gl.BindVertexArray(quadVAO)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		gl.BindVertexArray(0)
	}

	// combine the normal and blurry bright texture for a bloom effect
	gl.BindFramebuffer(gl.FRAMEBUFFER, b.bloomFbo.fbo)
	b.blendShader.Use()
	gl.DrawBuffer(gl.COLOR_ATTACHMENT1)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(uniformLocation(b.blendShader, "screenTexture"), 0)
	gl.BindTexture(gl.TEXTURE_2D, b.bloomFbo.textures[0])
	gl.ActiveTexture(gl.TEXTURE1)
	gl.Uniform1i(uniformLocation(b.blendShader, "bloomTexture"), 1)
	gl.BindTexture(gl.TEXTURE_2D, b.pingBuffers[1].textures[0])
	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
	return b.bloomFbo.textures[1]
}
