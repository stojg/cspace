package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/cspace/lib/shaders"
)

type AverageExposure struct {
	fbo        uint32
	passShader *shaders.Passthrough
	textures   [2]uint32
	exposure   float32
	firstFrame bool
}

func NewAverageExposure() *AverageExposure {

	a := &AverageExposure{
		firstFrame: true,
		passShader: shaders.NewPassthrough(),
		exposure:   1.0,
	}
	gl.GenFramebuffers(1, &a.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, a.fbo)

	for i := 0; i < len(a.textures); i++ {
		GLTextureRGB16F(&a.textures[i], windowWidth, windowHeight, gl.LINEAR, gl.CLAMP_TO_EDGE, nil)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+uint32(i), gl.TEXTURE_2D, a.textures[i], 0)
	}

	chkFramebuffer()

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	return a
}

// 0.5 Calculate average luminence of last frame's scene (pingpongColorBuffer[0] filled at end of loop (see last section of code))
func (a *AverageExposure) Exposure(inTexture uint32) float32 {
	// for the first frame, we do not have any previous frame to calculate from
	if a.firstFrame {
		a.firstFrame = false
		return 1.0
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, a.fbo)
	var attachments = [1]uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &attachments[0])

	gl.UseProgram(a.passShader.Program)
	GLBindTexture(0, a.passShader.LocScreenTexture, inTexture)
	renderQuad()

	texWidth := windowWidth
	texHeight := windowHeight

	readTexture := 0
	writeTexture := 1
	// Then pingpong between color buffers creating a smaller texture every time
	for texWidth > 1 {

		// first change texture size
		texWidth = texWidth / 2
		texHeight = texHeight / 2
		if texWidth < 1 {
			texWidth = 1
		}
		if texHeight < 1 {
			texHeight = 1
		}
		// @todo this is slow rought 5ms/frame on macbook pro
		gl.BindTexture(gl.TEXTURE_2D, a.textures[writeTexture])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, texWidth, texHeight, 0, gl.RGB, gl.FLOAT, nil)

		var attachments [1]uint32
		if readTexture == 0 {
			attachments[0] = gl.COLOR_ATTACHMENT1
		} else {
			attachments[0] = gl.COLOR_ATTACHMENT0
		}
		gl.DrawBuffers(1, &attachments[0])

		gl.Viewport(0, 0, texWidth, texHeight)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(a.passShader.LocScreenTexture, 0)
		gl.BindTexture(gl.TEXTURE_2D, a.textures[readTexture])
		renderQuad()

		writeTexture, readTexture = readTexture, writeTexture
	}

	// Once done read the luminescence value of 1x1 texture
	var luminescence [3]float32
	gl.ReadPixels(0, 0, 1, 1, gl.RGB, gl.FLOAT, gl.Ptr(&luminescence[0]))
	lum := 0.2126*luminescence[0] + 0.7152*luminescence[1] + 0.0722*luminescence[2]
	if lum < 0.1 {
		lum = 0.1
	} else if lum > 0.9 {
		lum = 0.9
	}
	a.exposure = Lerp(a.exposure, 0.5/lum, 0.2) // slowly adjust exposure based on average brightness

	for i := 0; i < 2; i++ {
		// Reset color buffer dimensions
		gl.BindTexture(gl.TEXTURE_2D, a.textures[i])
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, windowWidth, windowHeight, 0, gl.RGB, gl.FLOAT, nil)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return a.exposure
}
