package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/cspace/lib/shaders"
)

type Fxaa struct {
	width, height int32
	fbo           uint32
	texture       uint32

	ShowEdges     int32
	LumaThreshold float32
	MulReduce     float32
	MinReduce     float32
	MaxSpan       float32

	shader *shaders.Fxaa
}

func NewFxaa(width, height int32) *Fxaa {
	c := shaders.NewFxaa()
	fxaa := &Fxaa{
		shader:        c,
		width:         width,
		height:        height,
		LumaThreshold: 0.6,
		MulReduce:     1 / 8.0,
		MinReduce:     1 / 128.0,
		MaxSpan:       8.0,
	}

	GLFramebuffer(&fxaa.fbo)
	GLTextureRGB8(&fxaa.texture, width, height, gl.LINEAR, gl.CLAMP_TO_EDGE, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fxaa.texture, 0)

	chkFramebuffer()

	return fxaa
}

func (fxaa *Fxaa) Render(inTexture uint32) uint32 {
	gl.Disable(gl.DEPTH_TEST)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fxaa.fbo)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
	gl.UseProgram(fxaa.shader.Program)
	chkError("here")

	if showDebug {
		gl.Uniform1i(fxaa.shader.LocShowEdges, 1)
	} else {
		gl.Uniform1i(fxaa.shader.LocShowEdges, 0)
	}
	gl.Uniform1f(fxaa.shader.LocLumaThreshold, fxaa.LumaThreshold)
	gl.Uniform1f(fxaa.shader.LocMinReduce, fxaa.MinReduce)
	gl.Uniform1f(fxaa.shader.LocMulReduce, fxaa.MulReduce)
	gl.Uniform1f(fxaa.shader.LocMaxSpan, fxaa.MaxSpan)
	GLBindTexture(0, fxaa.shader.LocInTexture, inTexture)

	renderQuad()

	return fxaa.texture
}
