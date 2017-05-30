package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func GLBindTexture(pos int, loc int32, textureID uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(pos))
	gl.Uniform1i(loc, int32(pos))
	gl.BindTexture(gl.TEXTURE_2D, textureID)
}

func GLBindCubeMap(pos int, loc int32, textureID uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(pos))
	gl.Uniform1i(loc, int32(pos))
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, textureID)
}

func GLFramebuffer(fboID *uint32) {
	gl.GenFramebuffers(1, fboID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, *fboID)
}

func GLTextureRGB16F(textureID *uint32, width, height int32, filter, wrap int32, pixels unsafe.Pointer) {
	gl.GenTextures(1, textureID)
	gl.BindTexture(gl.TEXTURE_2D, *textureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, width, height, 0, gl.RGB, gl.FLOAT, pixels)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, wrap)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, wrap)
}
