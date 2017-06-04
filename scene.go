package main

import (
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/cspace/lib/shaders"
)

const near float32 = 0.5
const far float32 = 200

const maxPointLights = 64

const sizeUboScalar = 4
const sizeUboMat4 = 16 * sizeUboScalar
const sizeUboVec3 = 4 * sizeUboScalar

var windowWidth int32 = 1280
var windowHeight int32 = 720
var viewPortWidth int32
var viewPortHeight int32

var bloomOn = true
var ssaoOn = true
var fxaaOn = true
var dirLightOn = true
var skyBoxOn = true
var showDebug = false

var currentNumLights = 0

var directionLight = &DirectionalLight{
	Direction: normalise([3]float32{1, 0.7, 0}),
	Color:     [3]float32{5, 5, 6},
}

func NewScene() *Scene {

	s := &Scene{
		gBuffer:        NewGBufferPipeline(),
		shadow:         NewShadow(directionLight),
		camera:         NewCamera(),
		projection:     mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), near, far),
		graph:          NewBaseNode(),
		lightBoxShader: shaders.NewEmissive(),
	}

	att := ligthAtt[1]
	for i := 0; i < maxPointLights; i++ {
		s.pointLights = append(s.pointLights, &PointLight{
			Position: [3]float32{rand.Float32()*60 - 30, rand.Float32()*5 + 1, rand.Float32()*60 - 30},
			Color:    [3]float32{rand.Float32() * 20, rand.Float32() * 20, rand.Float32() * 20},
			Constant: att.Constant,
			Linear:   att.Linear,
			Exp:      att.Exp,
			rand:     rand.Float32() * 2,
		})
	}
	chkError("end_of_new_scene")
	return s
}

type Scene struct {
	projection mgl32.Mat4
	camera     *Camera
	graph      SceneNode

	gBuffer *GBufferPipeline
	bloom   *BloomEffect
	shadow  *ShadowFBO
	ssao    *SsaoFBO
	hdr     *HDRFBO
	skybox  *Skybox
	ibl     *IBL
	fxaa    *Fxaa
	tonemap *ToneMap

	pointLightShader *shaders.PointLight
	dirLightShader   *shaders.DirectionalLight
	lightBoxShader   *shaders.Emissive
	passShader       *shaders.Passthrough

	pointLights []*PointLight

	uboMatrices uint32
}

func (s *Scene) Init() {

	s.dirLightShader = shaders.NewDirectionalLight()
	s.pointLightShader = shaders.NewPointLightShader(maxPointLights)
	s.bloom = NewBloomEffect(windowWidth/2, windowHeight/2)
	s.ssao = NewSSAO(windowWidth, windowHeight)
	s.hdr = NewHDRFBO()
	s.ibl = NewCubeMap(512, 512)
	s.skybox = NewSkymap(s.ibl.envCubeMap)

	s.fxaa = NewFxaa(windowWidth, windowHeight)
	s.tonemap = NewToneMap(windowWidth, windowHeight)
	s.passShader = shaders.NewPassthrough()

	s.ibl.Update(GetHDRTexture("sky0016.hdr"))

	gl.GenBuffers(1, &s.uboMatrices)
	gl.BindBuffer(gl.UNIFORM_BUFFER, s.uboMatrices)
	gl.BufferData(gl.UNIFORM_BUFFER, 4*sizeUboMat4+sizeUboVec3, gl.Ptr(nil), gl.STATIC_DRAW)
	// link a specific range of the buffer which in this case is the entire buffer, to binding point 0.
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 0, s.uboMatrices, 0, 4*sizeUboMat4+sizeUboVec3)

	proj := s.projection
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, sizeUboMat4, gl.Ptr(&proj[0]))

	invP := s.projection.Inv()
	gl.BufferSubData(gl.UNIFORM_BUFFER, 2*sizeUboMat4, sizeUboMat4, gl.Ptr(&invP[0]))
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)

	chkError("scene.init")
}

func (s *Scene) Render(elapsed float64) {

	handleInputs()

	view := s.camera.View(elapsed)
	s.updateMatrices(view)

	shadowMap := s.shadow.Render(s.graph)

	s.gBuffer.Render(s.graph)

	aoTexture := s.ssao.Render(s.gBuffer.buffer.gDepth, s.gBuffer.buffer.gNormalRoughness)

	// start drawing light calculations into the finalTexture of the gbuffer
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, s.gBuffer.buffer.fbo)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT3)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// start the light blending
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.ONE, gl.ONE)
	{ // point light pass
		gl.UseProgram(s.pointLightShader.Program)
		gl.Uniform1i(s.pointLightShader.LocNumLights, int32(currentNumLights))

		// for now I have settled with passing all lights into one shader pass. Stencil pass turned out to be very hard
		// to get correct and not very fast. This could possibly be more performant by using tiled deferred rendering.
		for i := range s.pointLights {
			gl.Uniform3f(s.pointLightShader.LocLightPos[i], s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			gl.Uniform3fv(s.pointLightShader.LocLightColor[i], 1, &s.pointLights[i].Color[0])
			gl.Uniform1f(s.pointLightShader.LocLightLinear[i], s.pointLights[i].Linear)
			gl.Uniform1f(s.pointLightShader.LocLightQuadratic[i], s.pointLights[i].Exp)
		}

		GLBindTexture(0, s.pointLightShader.LocGDepth, s.gBuffer.buffer.gDepth)
		GLBindTexture(1, s.pointLightShader.LocGNormal, s.gBuffer.buffer.gNormalRoughness)
		GLBindTexture(2, s.pointLightShader.LocGAlbedo, s.gBuffer.buffer.gAlbedoMetallic)
		GLBindTexture(3, s.pointLightShader.LocGAmbientOcclusion, aoTexture)
		gl.Uniform2f(s.pointLightShader.LocScreenSize, float32(windowWidth), float32(windowHeight))
		renderQuad()
	}
	{
		gl.UseProgram(s.dirLightShader.Program)
		GLBindTexture(0, s.dirLightShader.LocGDepth, s.gBuffer.buffer.gDepth)
		GLBindTexture(1, s.dirLightShader.LocGNormal, s.gBuffer.buffer.gNormalRoughness)
		GLBindTexture(2, s.dirLightShader.LocGAlbedo, s.gBuffer.buffer.gAlbedoMetallic)
		GLBindTexture(3, s.dirLightShader.LocShadowMap, shadowMap)
		GLBindTexture(4, s.dirLightShader.LocGAmbientOcclusion, aoTexture)
		if skyBoxOn {
			gl.Uniform1i(s.dirLightShader.LocIBLEnabled, 1)
			GLBindCubeMap(5, s.dirLightShader.LocIrradianceMap, s.ibl.irradianceMap)
			GLBindCubeMap(6, s.dirLightShader.LocPrefilterMap, s.ibl.prefilterMap)
			GLBindTexture(7, s.dirLightShader.LocPbrdfLUT, s.ibl.brdfLUTTexture)
		} else {
			gl.Uniform1i(s.dirLightShader.LocIBLEnabled, 0)
		}
		gl.UniformMatrix4fv(s.dirLightShader.LocLightProjection, 1, false, &s.shadow.Projection[0])
		gl.UniformMatrix4fv(s.dirLightShader.LocLightView, 1, false, &s.shadow.View[0])
		gl.Uniform2f(s.dirLightShader.LocScreenSize, float32(windowWidth), float32(windowHeight))
		if dirLightOn {
			gl.Uniform1i(s.dirLightShader.LocLightEnabled, 1)
			gl.Uniform3fv(s.dirLightShader.LocLightDirection, 1, &directionLight.Direction[0])
			gl.Uniform3fv(s.dirLightShader.LocLightColor, 1, &directionLight.Color[0])
		} else {
			gl.Uniform1i(s.dirLightShader.LocLightEnabled, 0)
		}
		renderQuad()
	}
	gl.Disable(gl.BLEND)

	{ // render emissive objects
		gl.Enable(gl.DEPTH_TEST)
		gl.UseProgram(s.lightBoxShader.Program)
		for i := range s.pointLights[:currentNumLights] {
			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(0.1, 0.1, 0.1))
			model = model.Mul4(mgl32.HomogRotate3D(float32(math.Cos(glfw.GetTime())), mgl32.Vec3{1, 1, 1}.Normalize()))
			gl.UniformMatrix4fv(s.lightBoxShader.LocModel, 1, false, &model[0])
			gl.Uniform3fv(s.lightBoxShader.LocColor, 1, &s.pointLights[i].Color[0])
			renderCube()
		}
	}

	if skyBoxOn {
		s.skybox.Render(view, s.ibl.envCubeMap)
	}
	out := s.gBuffer.buffer.finalTexture
	if bloomOn {
		out = s.bloom.Render(out)
	}
	out = s.tonemap.Render(out, 1.4)
	if fxaaOn {
		out = s.fxaa.Render(out)
	}

	s.renderToViewport(out)

	// and if debug is on, quad print them on top of everything
	if showDebug {
		gl.Disable(gl.DEPTH_TEST)
		DisplayNormalBufferTexture(s.gBuffer.buffer.gNormalRoughness)
		DisplayAlbedoTexBuffer(s.gBuffer.buffer.gAlbedoMetallic)
		DisplayDepthbufferTexture(s.gBuffer.buffer.gDepth)
		DisplaySsaoTexture(aoTexture)
		DisplayShadowTexture(s.shadow.depthMap)
		DisplayRoughnessTexture(s.gBuffer.buffer.gNormalRoughness)
		DisplayMetallicTexture(s.gBuffer.buffer.gAlbedoMetallic)
		DisplayBloomTexture(s.bloom.pingBuffers[1].textures[0])
	}

	// check if there was an opengl error in this frame, in that case panic
	chkError("end_of_frame")
}

func (s *Scene) updateMatrices(view mgl32.Mat4) {
	gl.BindBuffer(gl.UNIFORM_BUFFER, s.uboMatrices)
	gl.BufferSubData(gl.UNIFORM_BUFFER, sizeUboMat4, sizeUboMat4, gl.Ptr(&view[0]))
	invView := view.Inv()
	gl.BufferSubData(gl.UNIFORM_BUFFER, 3*sizeUboMat4, sizeUboMat4, gl.Ptr(&invView[0]))
	camPos := s.camera.position
	gl.BufferSubData(gl.UNIFORM_BUFFER, 4*sizeUboMat4, sizeUboVec3, gl.Ptr(&camPos[0]))
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
}

func (s *Scene) renderToViewport(outTexture uint32) {
	gl.Disable(gl.DEPTH_TEST)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	// taking care of retina screens have different amount of pixel between actual viewport and requested window size
	gl.Viewport(0, 0, viewPortWidth, viewPortHeight)
	gl.UseProgram(s.passShader.Program)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	GLBindTexture(0, s.passShader.LocScreenTexture, outTexture)
	renderQuad()
	GLUnbindTexture(0)
}

func handleInputs() {
	if keys[glfw.Key1] {
		skyBoxOn = true
		dirLightOn = true
	} else if keys[glfw.Key2] {
		currentNumLights = 1
	} else if keys[glfw.Key3] {
		currentNumLights = 2
	} else if keys[glfw.Key4] {
		currentNumLights = 4
	} else if keys[glfw.Key5] {
		currentNumLights = 8
	} else if keys[glfw.Key6] {
		currentNumLights = 16
	} else if keys[glfw.Key6] {
		currentNumLights = 32
	} else if keys[glfw.Key7] {
		currentNumLights = 64
	} else if keys[glfw.Key0] {
		dirLightOn = false
		skyBoxOn = false
		currentNumLights = 0
	} else if keys[glfw.KeyF] {
		fxaaOn = true
	} else if keys[glfw.KeyR] {
		fxaaOn = false
	} else if keys[glfw.KeyTab] {
		showDebug = true
	} else if keys[glfw.KeyEnter] {
		bloomOn = true
	} else if keys[glfw.KeyG] {
		ssaoOn = true
	} else if keys[glfw.KeyT] {
		ssaoOn = false
	} else if keys[glfw.KeyEscape] {
		currentNumLights = 0
		dirLightOn = true
		skyBoxOn = true
		bloomOn = false
		fxaaOn = false
		ssaoOn = false
		showDebug = false
	}

}
