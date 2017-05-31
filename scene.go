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

var bloomOn = true // 2ms
var ssaoOn = true  // 8ms ?
var fxaaOn = true
var dirLightOn = true
var showDebug = false

var currentNumLights = 0

var directionLight = &DirectionalLight{
	Direction: normalise([3]float32{-80, 60, -100}),
	Color:     [3]float32{20, 20, 20},
}

func NewScene() *Scene {

	s := &Scene{
		gBuffer: NewGBufferPipeline(),

		shadow:         NewShadow(),
		previousTime:   glfw.GetTime(),
		camera:         NewCamera(),
		projection:     mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), near, far),
		graph:          NewBaseNode(),
		lightBoxShader: shaders.NewEmissive(),
		icoMesh:        LoadModel("models/ico", MaterialMesh)[0],
		cubeMesh:       LoadModel("models/cube", MaterialMesh)[0],
	}

	gl.GenBuffers(1, &s.uboMatrices)
	gl.BindBuffer(gl.UNIFORM_BUFFER, s.uboMatrices)
	gl.BufferData(gl.UNIFORM_BUFFER, 4*sizeUboMat4+sizeUboVec3, gl.Ptr(nil), gl.STATIC_DRAW) // allocate 150 bytes of memory
	// link a specific range of the buffer which in this case is the entire buffer, to binding point 0.
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 0, s.uboMatrices, 0, 4*sizeUboMat4+sizeUboVec3)

	proj := s.projection
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, sizeUboMat4, gl.Ptr(&proj[0]))

	invP := s.projection.Inv()
	gl.BufferSubData(gl.UNIFORM_BUFFER, 2*sizeUboMat4, sizeUboMat4, gl.Ptr(&invP[0]))
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)

	att := ligthAtt[1]
	for i := 0; i < maxPointLights; i++ {
		s.pointLights = append(s.pointLights, &PointLight{
			Position: [3]float32{rand.Float32()*60 - 30, rand.Float32()*5 + 1, rand.Float32()*60 - 30},
			Color:    [3]float32{rand.Float32() * 10, rand.Float32() * 10, rand.Float32() * 10},
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
	previousTime float64
	elapsed      float32
	projection   mgl32.Mat4
	camera       *Camera
	graph        SceneNode

	gBuffer *GBufferPipeline
	bloom   *BloomEffect

	shadow   *ShadowFBO
	ssao     *SsaoFBO
	hdr      *HDRFBO
	exposure *AverageExposure

	pointLightShader *shaders.PointLight
	dirLightShader   *shaders.DirectionalLight
	lightBoxShader   *shaders.Emissive
	passShader       *shaders.Passthrough

	pointLights []*PointLight
	icoMesh     *Mesh
	cubeMesh    *Mesh

	stencilShader *shaders.Stencil

	skybox        *shaders.Skybox
	skyBoxTexture *Texture
	cubeMap       *IBL

	fxaa    *Fxaa
	tonemap *ToneMap

	uboMatrices uint32
}

func (s *Scene) Init() {
	gl.Disable(gl.FRAMEBUFFER_SRGB)
	s.stencilShader = shaders.NewStencil()
	s.dirLightShader = shaders.NewDirectionalLight()
	s.pointLightShader = shaders.NewPointLightShader(maxPointLights)
	s.bloom = NewBloomEffect(windowWidth/2, windowHeight/2)
	s.ssao = NewSSAO(windowWidth/2, windowHeight/2)
	s.hdr = NewHDRFBO()
	s.exposure = NewAverageExposure()

	s.skybox = shaders.NewSkybox()
	s.skyBoxTexture = GetHDRTexture("woods_1k.hdr")
	s.cubeMap = NewCubeMap(512, 512)

	s.cubeMap.Update(s.skyBoxTexture)

	s.fxaa = NewFxaa(windowWidth, windowHeight)
	s.tonemap = NewToneMap(windowWidth, windowHeight)
	s.passShader = shaders.NewPassthrough()

	chkError("scene.init")

}

func (s *Scene) Render() {

	now := glfw.GetTime()
	s.elapsed = float32(now - s.previousTime)
	s.previousTime = now

	handleInputs()

	sin := float32(math.Sin(glfw.GetTime()))

	gl.BindBuffer(gl.UNIFORM_BUFFER, s.uboMatrices)
	view := s.camera.View(s.elapsed)
	gl.BufferSubData(gl.UNIFORM_BUFFER, sizeUboMat4, sizeUboMat4, gl.Ptr(&view[0]))
	invView := view.Inv()
	gl.BufferSubData(gl.UNIFORM_BUFFER, 3*sizeUboMat4, sizeUboMat4, gl.Ptr(&invView[0]))
	camPos := s.camera.position
	gl.BufferSubData(gl.UNIFORM_BUFFER, 4*sizeUboMat4, sizeUboVec3, gl.Ptr(&camPos[0]))
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
	chkError("ubo binding")

	// @todo move somewhere and calculate the proper bounding box
	lightProjection := mgl32.Ortho(-44, 40, -25, 25, -45, 40)
	lightView := mgl32.LookAt(directionLight.Direction[0], directionLight.Direction[1], directionLight.Direction[2], 0, 0, 0, 0, 1, 0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	{ // get the directional lights shadow mask and push that into a shadow depth texture
		gl.BindFramebuffer(gl.FRAMEBUFFER, s.shadow.fbo)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		s.shadow.shader.Use()
		lightSpaceMatrix := lightProjection.Mul4(lightView)
		gl.UniformMatrix4fv(s.shadow.locLightSpaceMatrix, 1, false, &lightSpaceMatrix[0])

		gl.Viewport(0, 0, s.shadow.Width, s.shadow.Height)
		s.graph.SimpleRender(s.shadow.shader)
		gl.Viewport(0, 0, windowWidth, windowHeight)
	}

	{ // render into the gBuffer
		gl.BindFramebuffer(gl.FRAMEBUFFER, s.gBuffer.buffer.fbo)
		var attachments = [2]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1}
		gl.DrawBuffers(int32(len(attachments)), &attachments[0])
		gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)
		s.graph.Render(s.gBuffer.tShader, s.gBuffer.mShader)
	}

	// we have written all the depth buffer information we wanted.
	gl.DepthMask(false)

	// we wont be needing the depth tests for until we start a forward rendering again
	gl.Disable(gl.DEPTH_TEST)

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
			gl.Uniform3f(s.pointLightShader.LocLightPos[i], s.pointLights[i].Position[0], s.pointLights[i].Position[1]+sin, s.pointLights[i].Position[2])
			gl.Uniform3fv(s.pointLightShader.LocLightColor[i], 1, &s.pointLights[i].Color[0])
			gl.Uniform1f(s.pointLightShader.LocLightLinear[i], s.pointLights[i].Linear)
			gl.Uniform1f(s.pointLightShader.LocLightQuadratic[i], s.pointLights[i].Exp)
		}

		GLBindTexture(0, s.pointLightShader.LocGDepth, s.gBuffer.buffer.gDepth)
		GLBindTexture(1, s.pointLightShader.LocGNormal, s.gBuffer.buffer.gNormalRoughness)
		GLBindTexture(2, s.pointLightShader.LocGAlbedo, s.gBuffer.buffer.gAlbedoMetallic)
		GLBindTexture(3, s.pointLightShader.LocGAmbientOcclusion, s.ssao.texture)

		gl.Uniform2f(s.pointLightShader.LocScreenSize, float32(windowWidth), float32(windowHeight))

		renderQuad()
	}

	if dirLightOn { // Render the directional light, for now there is only one
		gl.UseProgram(s.dirLightShader.Program)

		GLBindTexture(0, s.dirLightShader.LocGDepth, s.gBuffer.buffer.gDepth)
		GLBindTexture(1, s.dirLightShader.LocGNormal, s.gBuffer.buffer.gNormalRoughness)
		GLBindTexture(2, s.dirLightShader.LocGAlbedo, s.gBuffer.buffer.gAlbedoMetallic)
		GLBindTexture(3, s.dirLightShader.LocShadowMap, s.shadow.depthMap)
		GLBindTexture(4, s.dirLightShader.LocGAmbientOcclusion, aoTexture)
		GLBindCubeMap(5, s.dirLightShader.LocIrradianceMap, s.cubeMap.irradianceMap)
		GLBindCubeMap(6, s.dirLightShader.LocPrefilterMap, s.cubeMap.prefilterMap)
		GLBindTexture(7, s.dirLightShader.LocPbrdfLUT, s.cubeMap.brdfLUTTexture)

		gl.UniformMatrix4fv(s.dirLightShader.LocLightProjection, 1, false, &lightProjection[0])
		gl.UniformMatrix4fv(s.dirLightShader.LocLightView, 1, false, &lightView[0])
		gl.Uniform3fv(s.dirLightShader.LocLightDirection, 1, &directionLight.Direction[0])
		gl.Uniform3fv(s.dirLightShader.LocLightColor, 1, &directionLight.Color[0])
		gl.Uniform2f(s.dirLightShader.LocScreenSize, float32(windowWidth), float32(windowHeight))

		renderQuad()
	}
	gl.Disable(gl.BLEND)

	// start a forward rendering pass from here
	gl.Enable(gl.DEPTH_TEST)

	// render the skybox
	if dirLightOn {
		gl.DepthFunc(gl.LEQUAL)
		gl.UseProgram(s.skybox.Program)
		// remove the rotation
		skyboxView := view.Mat3().Mat4()
		gl.UniformMatrix4fv(s.skybox.LocSkyView, 1, false, &skyboxView[0])
		GLBindCubeMap(0, s.skybox.LocScreenTexture, s.cubeMap.envCubeMap)
		gl.BindVertexArray(s.skybox.SkyboxVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}

	{ // render emissive objects
		gl.UseProgram(s.lightBoxShader.Program)

		for i := range s.pointLights[:currentNumLights] {
			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1]+sin, s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(0.1, 0.1, 0.1))
			model = model.Mul4(mgl32.HomogRotate3D(float32(math.Cos(glfw.GetTime())), mgl32.Vec3{1, 1, 1}.Normalize()))

			gl.UniformMatrix4fv(s.lightBoxShader.LocModel, 1, false, &model[0])
			gl.Uniform3fv(s.lightBoxShader.LocColor, 1, &s.pointLights[i].Color[0])

			renderCube()
		}
	}

	// Here we start the full screen effects, so we don't want any depth test shenanigans going on
	gl.Disable(gl.DEPTH_TEST)

	out := s.gBuffer.buffer.finalTexture
	if bloomOn {
		out = s.bloom.Render(out)
	}

	out = s.tonemap.Render(out)
	if fxaaOn {
		out = s.fxaa.Render(out)
	}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	// taking care of retina screens have different amount of pixel between actual viewport and requested window size
	gl.Viewport(0, 0, windowWidth, windowHeight)
	gl.UseProgram(s.passShader.Program)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	GLBindTexture(0, s.passShader.LocScreenTexture, out)
	renderQuad()
	GLUnbindTexture(0)

	//DisplayAlbedoTexBuffer(out)

	// and if debug is on, quad print them on top of everything
	if showDebug {
		DisplayNormalBufferTexture(s.gBuffer.buffer.gNormalRoughness)
		DisplayAlbedoTexBuffer(s.gBuffer.buffer.gAlbedoMetallic)
		DisplayDepthbufferTexture(s.gBuffer.buffer.gDepth)
		DisplaySsaoTexture(s.ssao.outTexture)
		DisplayShadowTexture(s.shadow.depthMap)
		DisplayRoughnessTexture(s.gBuffer.buffer.gNormalRoughness)
		DisplayMetallicTexture(s.gBuffer.buffer.gAlbedoMetallic)
		DisplayBloomTexture(s.bloom.pingBuffers[1].textures[0])
		//DisplayDepthbufferTexture(s.shadow.depthMap)
	}

	// check if there was an opengl error in this frame, in that case panic
	chkError("end_of_frame")
}

func handleInputs() {
	if keys[glfw.Key1] {
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
		bloomOn = false
		fxaaOn = false
		ssaoOn = false
		showDebug = false
	}

}
