package main

import "github.com/go-gl/gl/v4.1-core/gl"

// FBO is a generic FBO with one texture
type FBO struct {
	id       uint32
	textures []uint32
}

func NewFBO(width, height int32) *FBO {
	frameBuffer := &FBO{
		textures: make([]uint32, 1),
	}
	setFBO(&frameBuffer.id, frameBuffer.textures, width, height)
	return frameBuffer
}

func setFBO(fbo *uint32, textures []uint32, width, height int32) {
	gl.GenFramebuffers(1, fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, *fbo)

	for i := 0; i < len(textures); i++ {
		GLTextureRGB16F(&textures[i], width, height, gl.LINEAR, gl.CLAMP_TO_EDGE, nil)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+uint32(i), gl.TEXTURE_2D, textures[i], 0)
	}

	chkFramebuffer()

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
}
