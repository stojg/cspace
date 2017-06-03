package main

import (
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/cspace/lib/shaders"
)

// NewSSAO returns the Screen Space Ambient Occlusion effect
// on macbook air this takes ~10ms to render
func NewSSAO(width, height int32) *SsaoFBO {
	ssao := &SsaoFBO{
		Width:  width,
		Height: height,
	}

	GLFramebuffer(&ssao.fbo)

	gl.GenTextures(1, &ssao.positionTexture)
	gl.BindTexture(gl.TEXTURE_2D, ssao.positionTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA16F, windowWidth/2, windowHeight/2, 0, gl.RGB, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, ssao.positionTexture, 0)

	gl.GenTextures(1, &ssao.occlusionTexture)
	gl.BindTexture(gl.TEXTURE_2D, ssao.occlusionTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, ssao.Width/2, ssao.Height/2, 0, gl.RGB, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gl.TEXTURE_2D, ssao.occlusionTexture, 0)

	gl.GenTextures(1, &ssao.blur)
	gl.BindTexture(gl.TEXTURE_2D, ssao.blur)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, ssao.Width/2, ssao.Height/2, 0, gl.RGB, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT2, gl.TEXTURE_2D, ssao.blur, 0)

	gl.GenTextures(1, &ssao.outTexture)
	gl.BindTexture(gl.TEXTURE_2D, ssao.outTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, windowWidth/2, windowHeight/2, 0, gl.RGB, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT3, gl.TEXTURE_2D, ssao.outTexture, 0)

	for i := range ssao.Kernel {
		smp := mgl32.Vec3{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()}
		smp = smp.Normalize()
		scale := float32(i) / 64.0
		//scale samples s.t. they're more aligned to center of kernel
		scale = Lerp(0.1, 1.0, scale*scale)
		smp = smp.Mul(scale)
		ssao.Kernel[i][0] = smp[0]
		ssao.Kernel[i][1] = smp[1]
		ssao.Kernel[i][2] = smp[2] // 0.1 - 1.0 hemisphereic
	}

	var noise [3 * 16]float32
	for i := 0; i < len(noise); i += 3 {
		noise[i] = rand.Float32()*2 - 1
		noise[i+1] = rand.Float32()*2 - 1
		noise[i+2] = 0
	}
	GLTextureRGB16F(&ssao.noiseTexture, 4, 4, gl.NEAREST, gl.REPEAT, gl.Ptr(&noise[0]))
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, 4, 4, 0, gl.RGB, gl.FLOAT, gl.Ptr(&noise[0]))

	chkFramebuffer()

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	ssao.depthResampler = shaders.NewSSAODepthResampler()
	ssao.shader = shaders.NewSSAO()
	ssao.blurShader = shaders.NewSSAOGaussian()
	ssao.passShader = shaders.NewPassthrough()
	return ssao
}

type SsaoFBO struct {
	fbo              uint32
	positionTexture  uint32
	blur             uint32
	occlusionTexture uint32
	noiseTexture     uint32
	outTexture       uint32
	depthResampler   *shaders.SSAODepthResampler
	shader           *shaders.SSAO
	blurShader       *shaders.SSAOGaussian
	passShader       *shaders.Passthrough

	Width  int32
	Height int32
	Kernel [64][3]float32
}

func (s *SsaoFBO) Render(gDepthTexture, gNormalTexture uint32) uint32 {
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.fbo)

	gl.Viewport(0, 0, s.Width/2, s.Height/2)

	// first downsample the gbuffer depth texture and calculate worldspace positions
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
	gl.UseProgram(s.depthResampler.Program)
	GLBindTexture(0, s.depthResampler.LocGDepth, gDepthTexture)
	gl.Uniform2f(s.depthResampler.LocScreenSize, float32(s.Width/2), float32(s.Height/2))
	renderQuad()

	// calculate the ambient occlusion
	gl.DrawBuffer(gl.COLOR_ATTACHMENT1)
	gl.UseProgram(s.shader.Program)
	if ssaoOn {
		gl.Uniform1i(s.shader.LocEnabled, 1)
	} else {
		gl.Uniform1i(s.shader.LocEnabled, 0)
	}
	// Send kernel (samples)
	for i, sample := range s.Kernel {
		gl.Uniform3f(s.shader.LocSamples[i], sample[0], sample[1], sample[2])
	}
	gl.Uniform2f(s.shader.LocScreenSize, float32(s.Width/2), float32(s.Height/2))

	GLBindTexture(0, s.shader.LocGDepth, s.positionTexture)
	GLBindTexture(1, s.shader.LocGNormal, gNormalTexture)
	GLBindTexture(2, s.shader.LocTexNoise, s.noiseTexture)
	renderQuad()

	// vertical gaussian blur pass
	gl.DrawBuffer(gl.COLOR_ATTACHMENT2)
	gl.UseProgram(s.blurShader.Program)
	gl.Uniform1i(s.blurShader.LocHorizontal, int32(0))
	GLBindTexture(0, s.blurShader.LocScreenTexture, s.occlusionTexture)
	renderQuad()

	// horisontal gaussian blur pass into final out texture
	gl.DrawBuffer(gl.COLOR_ATTACHMENT3)
	gl.UseProgram(s.blurShader.Program)
	gl.Uniform1i(s.blurShader.LocHorizontal, int32(1))
	GLBindTexture(0, s.blurShader.LocScreenTexture, s.blur)
	renderQuad()

	gl.Viewport(0, 0, s.Width, s.Height)
	return s.outTexture
}
