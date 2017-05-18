package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShadowFBO struct {
	fbo                 uint32
	eh                  uint32
	depthMap            uint32
	shader              *ShadowShader
	locLightSpaceMatrix int32
	Width               int32
	Height              int32
}

func NewShadow() *ShadowFBO {
	shadow := &ShadowFBO{
		Width:  1024,
		Height: 1024,
	}

	gl.GenFramebuffers(1, &shadow.fbo)

	gl.GenTextures(1, &shadow.depthMap)
	gl.BindTexture(gl.TEXTURE_2D, shadow.depthMap)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, shadow.Width, shadow.Height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	borderColor := [4]float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	gl.BindFramebuffer(gl.FRAMEBUFFER, shadow.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, shadow.depthMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)

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

	shadow.shader = &ShadowShader{
		DefaultShader: NewDefaultShader("shadow", "shadow"),
	}
	shadow.locLightSpaceMatrix = uniformLocation(shadow.shader, "lightSpaceMatrix")
	shadow.shader.uniformModelLoc = uniformLocation(shadow.shader.DefaultShader, "model")
	return shadow
}

type ShadowShader struct {
	*DefaultShader
	uniformModelLoc int32
}

func (s *ShadowShader) ModelUniform() int32 {
	return s.uniformModelLoc
}
