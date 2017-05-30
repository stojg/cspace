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
	equirectangularToCubemapShader *shaders.EquiRectToCubeMap
	envCubeMap                     uint32
	irradianceShader               *shaders.IBLIrradiance
	irradianceMap                  uint32
	prefilterMap                   uint32
	prefilterShader                *shaders.IBLPreFilter
	brdfLUTTexture                 uint32
	brdfShader                     *shaders.IBLBrdf
}

func NewCubeMap(width, height int32) *CubeMap {
	cube := &CubeMap{
		width:  width,
		height: height,
		equirectangularToCubemapShader: shaders.NewEquiRectToCubeMap(),
		irradianceShader:               shaders.NewIBLIrradiance(),
		prefilterShader:                shaders.NewIBLPrefilter(),
		brdfShader:                     shaders.NewIBLBrdf(),
	}

	gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)

	gl.GenFramebuffers(1, &cube.fbo)
	gl.GenRenderbuffers(1, &cube.rbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, cube.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, cube.width, cube.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, cube.rbo)

	generateCubeMap(&cube.envCubeMap, cube.width, cube.height, true, false)

	//gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)
	//gl.BindRenderbuffer(gl.RENDERBUFFER, cube.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, 32, 32)

	generateCubeMap(&cube.irradianceMap, 32, 32, false, false)

	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, 128, 128)
	generateCubeMap(&cube.prefilterMap, 128, 128, true, true)

	// pre-allocate enough memory for the LUT texture.
	gl.GenTextures(1, &cube.brdfLUTTexture)
	gl.BindTexture(gl.TEXTURE_2D, cube.brdfLUTTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RG16F, width, height, 0, gl.RG, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("FRAMEBUFFER_COMPLETE error, s: 0x%x\n", s))
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return cube
}

func (cube *CubeMap) Update(texture *Texture) {
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

	// convert HDR equirectangular environment map to cubemap equivalent

	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, cube.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, cube.width, cube.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, cube.rbo)

	gl.UseProgram(cube.equirectangularToCubemapShader.Program)

	GLBindTexture(0, cube.equirectangularToCubemapShader.LocEquirectangularMap, texture.ID)

	gl.UniformMatrix4fv(cube.equirectangularToCubemapShader.LocProjection, 1, false, &captureProjection[0])
	gl.Viewport(0, 0, cube.width, cube.height)
	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)

	for i := 0; i < 6; i++ {
		gl.UniformMatrix4fv(cube.equirectangularToCubemapShader.LocView, 1, false, &captureViews[i][0])
		gl.FramebufferTexture2D(
			gl.FRAMEBUFFER,
			gl.COLOR_ATTACHMENT0,
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			cube.envCubeMap,
			0,
		)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		renderCube()
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	// then let OpenGL generate mipmaps from first mip face (combatting visible dots artifact in prefilter map)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cube.envCubeMap)
	gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)

	// irradiance cube map
	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, cube.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, 32, 32)
	//
	gl.UseProgram(cube.irradianceShader.Program)
	gl.UniformMatrix4fv(cube.irradianceShader.LocProjection, 1, false, &captureProjection[0])
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(cube.irradianceShader.LocEnvironmentMap, 0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cube.envCubeMap)

	gl.Viewport(0, 0, 32, 32)
	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)
	for i := 0; i < 6; i++ {
		gl.UniformMatrix4fv(cube.irradianceShader.LocView, 1, false, &captureViews[i][0])
		gl.FramebufferTexture2D(
			gl.FRAMEBUFFER,
			gl.COLOR_ATTACHMENT0,
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			cube.irradianceMap,
			0,
		)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		renderCube()
	}

	// prefilter
	gl.UseProgram(cube.prefilterShader.Program)
	gl.UniformMatrix4fv(cube.prefilterShader.LocProjection, 1, false, &captureProjection[0])
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(cube.prefilterShader.LocEnvironmentMap, 0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cube.envCubeMap)
	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)

	maxMipLevels := 5
	for mip := 0; mip < maxMipLevels; mip++ {
		// resize framebuffer according to mip-level size.
		mipWidth := int32(128 * math.Pow(0.5, float64(mip)))
		mipHeight := int32(128 * math.Pow(0.5, float64(mip)))
		gl.BindRenderbuffer(gl.RENDERBUFFER, cube.rbo)
		gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, mipWidth, mipHeight)
		gl.Viewport(0, 0, mipWidth, mipHeight)

		roughness := float32(mip) / float32(maxMipLevels-1)
		gl.Uniform1f(cube.prefilterShader.LocRoughness, roughness)

		for i := 0; i < 6; i++ {
			gl.UniformMatrix4fv(cube.prefilterShader.LocView, 1, false, &captureViews[i][0])
			gl.FramebufferTexture2D(
				gl.FRAMEBUFFER,
				gl.COLOR_ATTACHMENT0,
				gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
				cube.prefilterMap,
				int32(mip),
			)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			renderCube()
		}
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, cube.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, cube.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, 512, 512)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, cube.brdfLUTTexture, 0)

	gl.Viewport(0, 0, 512, 512)
	gl.UseProgram(cube.brdfShader.Program)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	renderQuad()

	gl.Viewport(0, 0, windowWidth, windowHeight)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func generateCubeMap(texture *uint32, width, height int32, mipMap bool, genMipMap bool) {
	gl.GenTextures(1, texture)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, *texture)
	for i := 0; i < 6; i++ {
		// note that we store each face with 16 bit floating point values
		gl.TexImage2D(
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			0,
			gl.RGB16F,
			width,
			height,
			0,
			gl.RGB,
			gl.FLOAT,
			nil,
		)
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	if mipMap {
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}
	if genMipMap {
		gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
}
