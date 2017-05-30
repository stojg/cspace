package main

import (
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/cspace/lib/shaders"
)

func NewSSAO(width, height int32) *SsaoFBO {
	ssao := &SsaoFBO{
		Width:  width,
		Height: height,
	}

	GLFramebuffer(&ssao.fbo)

	gl.GenTextures(1, &ssao.texture)
	gl.BindTexture(gl.TEXTURE_2D, ssao.texture)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, ssao.Width, ssao.Height, 0, gl.RGB, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	borderColor := [4]float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	gl.GenTextures(1, &ssao.outTexture)
	gl.BindTexture(gl.TEXTURE_2D, ssao.outTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, windowWidth, windowHeight, 0, gl.RGB, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, ssao.texture, 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gl.TEXTURE_2D, ssao.outTexture, 0)

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

	ssao.shader = shaders.NewSSAO()
	ssao.blurShader = shaders.NewBlur()
	return ssao
}

type SsaoFBO struct {
	fbo          uint32
	texture      uint32
	noiseTexture uint32
	outTexture   uint32
	shader       *shaders.SSAO
	blurShader   *shaders.Blur

	Width  int32
	Height int32
	Kernel [64][3]float32
}

func (s *SsaoFBO) Render(gDepthTexture, gNormalTexture uint32) uint32 {
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.fbo)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0)
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
	gl.Uniform2f(s.shader.LocScreenSize, float32(s.Width), float32(s.Height))

	GLBindTexture(0, s.shader.LocGDepth, gDepthTexture)
	GLBindTexture(1, s.shader.LocGNormal, gNormalTexture)
	GLBindTexture(2, s.shader.LocTexNoise, s.noiseTexture)

	renderQuad()

	gl.DrawBuffer(gl.COLOR_ATTACHMENT1)
	gl.UseProgram(s.blurShader.Program)

	GLBindTexture(0, s.blurShader.LocScreenTexture, s.texture)
	renderQuad()

	return s.texture
}
