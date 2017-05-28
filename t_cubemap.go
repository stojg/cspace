package main

import (
	"fmt"

	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/cspace/lib/shaders"
)

// FBO is a generic FBO with one texture
type CubeMap struct {
	width, height                  int32
	fbo                            uint32
	rbo                            uint32
	envCubemap                     uint32
	equirectangularToCubemapShader *shaders.HDRCube
}

func NewCubeMap(width, height int32) *CubeMap {
	cube := &CubeMap{
		width:  width,
		height: height,
		equirectangularToCubemapShader: shaders.NewHDRCube(),
	}

	gl.GenFramebuffers(1, &cube.fbo)
	gl.GenRenderbuffers(1, &cube.rbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, cube.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, cube.width, cube.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, cube.rbo)

	gl.GenTextures(1, &cube.envCubemap)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cube.envCubemap)
	for i := 0; i < 6; i++ {
		// note that we store each face with 16 bit floating point values
		gl.TexImage2D(
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			0,
			gl.RGB16F,
			cube.width,
			cube.height,
			0,
			gl.RGB,
			gl.FLOAT,
			nil,
		)
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, s: 0x%x\n", s))
	}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	return cube
}

func (box *CubeMap) Update(texture *Texture) {
	fovy := (90 * math.Pi) / 180.0
	captureProjection := mgl32.Perspective(float32(fovy), 1, 0.1, 10)
	captureViews := []mgl32.Mat4{
		mgl32.LookAt(0, 0, 0, 1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, -1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 1, 0, 0, 0, 1),
		mgl32.LookAt(0, 0, 0, 0, -1, 0, 0, 0, -1),
		mgl32.LookAt(0, 0, 0, 0, 0, 1, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 0, -1, 0, -1, 0),
	}

	_ = captureProjection
	_ = captureViews

	// convert HDR equirectangular environment map to cubemap equivalent
	gl.UseProgram(box.equirectangularToCubemapShader.Program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(box.equirectangularToCubemapShader.LocEquirectangularMap, 0)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)

	gl.UniformMatrix4fv(box.equirectangularToCubemapShader.LocProjection, 1, false, &captureProjection[0])
	gl.Viewport(0, 0, box.width, box.height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, box.fbo)

	for i := 0; i < 6; i++ {
		gl.UniformMatrix4fv(box.equirectangularToCubemapShader.LocView, 1, false, &captureViews[i][0])
		gl.FramebufferTexture2D(
			gl.FRAMEBUFFER,
			gl.COLOR_ATTACHMENT0,
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			box.envCubemap,
			0,
		)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		renderCube()
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
