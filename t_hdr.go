package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type HDRFBO struct {
	fbo     uint32
	texture uint32
}

func NewHDRFBO() *HDRFBO {
	hdr := &HDRFBO{}

	gl.GenFramebuffers(1, &hdr.fbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, hdr.fbo)

	gl.GenTextures(1, &hdr.texture)
	gl.BindTexture(gl.TEXTURE_2D, hdr.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, windowWidth, windowHeight, 0, gl.RGB, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, hdr.texture, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
		switch s {
		case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
			panic("Framebuffer incomplete: No image is attached to FBO")
		case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
			panic("Framebuffer incomplete: Attachment is NOT complete")
		case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
			panic("Framebuffer incomplete: Draw buffer")
		case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
			panic("Framebuffer incomplete: Read buffer")
		default:
			panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, s: 0x%x\n", s))
		}

	}
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	return hdr
}
