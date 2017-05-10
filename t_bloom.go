package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Bloom struct {
	fbo      uint32
	textures [2]uint32
}

func NewBloom() *Bloom {
	frameBuffer := &Bloom{}
	gl.GenFramebuffers(1, &frameBuffer.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer.fbo)

	gl.GenTextures(2, &frameBuffer.textures[0])
	for i := 0; i < len(frameBuffer.textures); i++ {
		gl.BindTexture(gl.TEXTURE_2D, frameBuffer.textures[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	// Attach first fx texture to framebuffer
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, frameBuffer.textures[0], 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gl.TEXTURE_2D, frameBuffer.textures[1], 0)

	// - Finally check if framebuffer is complete
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, status: 0x%x\n", status))
	}

	// restore default FBO
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return frameBuffer
}
