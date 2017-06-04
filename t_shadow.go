package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShadowFBO struct {
	fbo                 uint32
	depthMap            uint32
	shader              *ShadowShader
	locLightSpaceMatrix int32
	Width               int32
	Height              int32
	View                mgl32.Mat4
	Projection          mgl32.Mat4
}

func NewShadow(light *DirectionalLight) *ShadowFBO {
	shadow := &ShadowFBO{
		Width:  1024 * 2,
		Height: 1024 * 2,
	}

	gl.GenFramebuffers(1, &shadow.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, shadow.fbo)

	gl.GenTextures(1, &shadow.depthMap)
	gl.BindTexture(gl.TEXTURE_2D, shadow.depthMap)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT16, shadow.Width, shadow.Height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	borderColor := [4]float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, shadow.depthMap, 0)

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
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	shadow.shader = &ShadowShader{
		DefaultShader: NewDefaultShader("shadow", "shadow"),
	}
	shadow.locLightSpaceMatrix = uniformLocation(shadow.shader, "lightSpaceMatrix")
	shadow.shader.uniformModelLoc = uniformLocation(shadow.shader.DefaultShader, "model")

	shadow.Projection = mgl32.Ortho(-44, 40, -25, 25, -45, 40)
	shadow.View = mgl32.LookAt(light.Direction[0], light.Direction[1], light.Direction[2], 0, 0, 0, 0, 1, 0)

	return shadow
}

// Render the directional lights shadow mask and push that into a shadow depth texture
func (s *ShadowFBO) Render(graph SceneNode) uint32 {
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.fbo)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	gl.UseProgram(s.shader.program)
	lightSpaceMatrix := s.Projection.Mul4(s.View)
	gl.UniformMatrix4fv(s.locLightSpaceMatrix, 1, false, &lightSpaceMatrix[0])

	gl.Viewport(0, 0, s.Width, s.Height)
	graph.SimpleRender(s.shader)
	gl.Viewport(0, 0, windowWidth, windowHeight)

	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return s.depthMap
}

type ShadowShader struct {
	*DefaultShader
	uniformModelLoc int32
}

func (s *ShadowShader) ModelUniform() int32 {
	return s.uniformModelLoc
}
