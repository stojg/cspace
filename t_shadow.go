package main

import "github.com/go-gl/gl/v4.1-core/gl"

type ShadowFBO struct {
	fbo                 uint32
	depthMap            uint32
	shader              *DefaultShader
	locLightSpaceMatrix int32
}

func NewShadow() *ShadowFBO {
	shadow := &ShadowFBO{}

	const shadowWidth int32 = 1024
	const shadowHeight int32 = 1024

	gl.GenFramebuffers(1, &shadow.fbo)

	gl.BindTexture(gl.TEXTURE_2D, shadow.depthMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, shadowWidth, shadowHeight, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	borderColor := [4]float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	//gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, shadow.depthMap, 0)
	//gl.BindFramebuffer(gl.FRAMEBUFFER, shadow.fbo)

	//if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
	//	panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, s: 0x%x\n", s))
	//}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	shadow.shader = NewDefaultShader("shadow", "shadow")
	shadow.locLightSpaceMatrix = uniformLocation(shadow.shader, "lightSpaceMatrix")

	return shadow
}
