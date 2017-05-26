package main

import (
	"fmt"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/cspace/lib/shaders"
)

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

func NewSSAO() *SsaoFBO {
	ssao := &SsaoFBO{
		Width:  windowWidth,
		Height: windowHeight,
	}

	gl.GenFramebuffers(1, &ssao.fbo)

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
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, ssao.Width, ssao.Height, 0, gl.RGB, gl.UNSIGNED_INT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.BindFramebuffer(gl.FRAMEBUFFER, ssao.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, ssao.texture, 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT1, gl.TEXTURE_2D, ssao.outTexture, 0)

	for i := range ssao.Kernel {
		smp := mgl32.Vec3{rand.Float32()*2 - 1, rand.Float32()*2 - 1, rand.Float32()}
		smp = smp.Normalize()
		scale := float32(i) / 64.0
		//scale samples s.t. they're more aligned to center of kernel
		scale = lerp(0.1, 1.0, scale*scale)
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
	gl.GenTextures(1, &ssao.noiseTexture)
	gl.BindTexture(gl.TEXTURE_2D, ssao.noiseTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB16F, 4, 4, 0, gl.RGB, gl.FLOAT, gl.Ptr(&noise[0]))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

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

	ssao.shader = shaders.NewSSAO()
	ssao.blurShader = shaders.NewBlur()
	chkError("Shader?")
	return ssao
}

func lerp(a, b, f float32) float32 {
	return a + f*(b-a)
}
