package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/cspace/lib/shaders"
)

func NewToneMap(width, height int32) *ToneMap {
	t := &ToneMap{
		width:  width,
		height: height,
		shader: shaders.NewHDR(),
	}
	GLFramebuffer(&t.fbo)
	GLTextureRGB8(&t.texture, width, height, gl.LINEAR, gl.CLAMP_TO_EDGE, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, t.texture, 0)
	chkFramebuffer()
	return t
}

type ToneMap struct {
	width, height int32
	fbo           uint32
	texture       uint32
	shader        *shaders.HDR
}

func (t *ToneMap) Render(inTexture uint32) uint32 {

	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)

	gl.UseProgram(t.shader.Program)

	GLBindTexture(0, t.shader.LocScreenTexture, inTexture)
	gl.Uniform1f(t.shader.LocExposure, 1.2)
	renderQuad()

	GLUnbindTexture(0)
	return t.texture

}
