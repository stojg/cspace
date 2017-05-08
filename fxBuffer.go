package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type FXFbo struct {
	id       uint32
	textures [4]uint32
}

func NewFXFbo() *FXFbo {
	fxGbo := &FXFbo{}
	gl.GenFramebuffers(1, &fxGbo.id)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fxGbo.id)

	gl.GenTextures(4, &fxGbo.textures[0])
	for i := 0; i < len(fxGbo.textures); i++ {
		gl.BindTexture(gl.TEXTURE_2D, fxGbo.textures[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16F, windowWidth, windowHeight, 0, gl.RGBA, gl.FLOAT, nil)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	// Attach first fx texture to framebuffer
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fxGbo.textures[0], 0)

	// - Finally check if framebuffer is complete
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, status: 0x%x\n", status))
	}

	// restore default FBO
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return fxGbo
}
