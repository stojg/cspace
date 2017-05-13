package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// FBO is a generic FBO with one texture
type FBO struct {
	id       uint32
	textures []uint32
}

func NewFBO() *FBO {
	frameBuffer := &FBO{
		textures: make([]uint32, 1),
	}
	setFBO(&frameBuffer.id, frameBuffer.textures)
	return frameBuffer
}

func setFBO(fbo *uint32, textures []uint32) {
	gl.GenFramebuffers(1, fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, *fbo)

	gl.GenTextures(int32(len(textures)), &textures[0])
	for i := 0; i < len(textures); i++ {
		gl.BindTexture(gl.TEXTURE_2D, textures[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+uint32(i), gl.TEXTURE_2D, textures[i], 0)
	}

	if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, s: 0x%x\n", s))
	}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
}
