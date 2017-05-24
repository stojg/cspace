package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/cspace/lib/shaders"
)

const near float32 = 0.5
const far float32 = 200
const windowWidth = 1280
const windowHeight = 720
const maxPointLights = 32

var viewPortWidth int32 = windowWidth
var viewPortHeight int32 = windowHeight
var bloomOn = true
var ssaoOn = true
var dirLightOn = true
var fxaaOn = false
var showDebug = false

var u_lumaThreshold float32 = 1 / 16.0 // (1/3), (1/4), (1/8), (1/16)
var u_mulReduce float32 = 1 / 8.0      //
var u_minReduce float32 = 1 / 128.0
var u_maxSpan float32 = 8.0

var currentNumLights = 1

var directionLight = &DirectionalLight{
	Direction: normalise([3]float32{-80, 60, -100}),
	Color:     [3]float32{1, 1, 1},
}

var fxaaShader *DefaultShader
var fxaaTextureloc int32

var fxaaLocU_showEdges int32
var fxaaLocU_lumaThreshold int32
var fxaaLocU_mulReduce int32
var fxaaLocU_minReduce int32
var fxaaLocU_maxSpan int32
var fxaaLoc_enabled int32

func NewScene() *Scene {

	fxaaShader = NewDefaultShader("fx", "fx_fxaa")
	fxaaTextureloc = uniformLocation(fxaaShader, "screenTexture")
	fxaaLocU_showEdges = pUniformLocation(fxaaShader.program, "u_showEdges")
	fxaaLocU_lumaThreshold = pUniformLocation(fxaaShader.program, "u_lumaThreshold")
	fxaaLocU_mulReduce = pUniformLocation(fxaaShader.program, "u_mulReduce")
	fxaaLocU_minReduce = pUniformLocation(fxaaShader.program, "u_minReduce")
	fxaaLocU_maxSpan = pUniformLocation(fxaaShader.program, "u_maxSpan")
	fxaaLoc_enabled = pUniformLocation(fxaaShader.program, "u_enabled")

	s := &Scene{
		gBuffer: NewGBufferPipeline(),

		shadow:         NewShadow(),
		bloom:          NewBloomEffect(),
		previousTime:   glfw.GetTime(),
		camera:         NewCamera(),
		projection:     mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), near, far),
		graph:          NewBaseNode(),
		lightBoxShader: NewDefaultShader("simple", "emissive"),
		icoMesh:        LoadModel("models/ico", MaterialMesh)[0],
		cubeMesh:       LoadModel("models/cube", MaterialMesh)[0],
	}
	s.invProjection = s.projection.Inv()

	s.pointLights = append(s.pointLights, &PointLight{
		Position: [3]float32{-3, 4, -2},
		Color:    [3]float32{10, 9, 8},
		Constant: ligthAtt[1].Constant,
		Linear:   ligthAtt[1].Linear,
		Exp:      ligthAtt[1].Exp,
		rand:     0,
	})

	att := ligthAtt[1]
	for i := 1; i < maxPointLights; i++ {
		s.pointLights = append(s.pointLights, &PointLight{
			Position: [3]float32{rand.Float32()*60 - 30, rand.Float32()*5 + 1, rand.Float32()*60 - 10},
			Color:    [3]float32{rand.Float32()*10 + 0.5, rand.Float32()*10 + 0.5, rand.Float32()*10 + 0.5},
			Constant: att.Constant,
			Linear:   att.Linear,
			Exp:      att.Exp,
			rand:     rand.Float32() * 2,
		})
	}
	s.lightBoxModelLoc = uniformLocation(s.lightBoxShader, "model")
	s.lightBoxEmissiveLoc = uniformLocation(s.lightBoxShader, "emissive")
	chkError("end_of_new_scene")
	return s
}

type Scene struct {
	previousTime  float64
	elapsed       float32
	projection    mgl32.Mat4
	invProjection mgl32.Mat4
	camera        *Camera
	graph         SceneNode

	gBuffer *GBufferPipeline
	bloom   *BloomEffect

	shadow *ShadowFBO
	ssao   *SsaoFBO
	hdr    *HDRFBO

	pointLightShader *shaders.PointLight
	dirLightShader   *shaders.DirectionalLight
	hdrShader        *shaders.HDR
	lightBoxShader   *DefaultShader

	pointLights []*PointLight
	icoMesh     *Mesh
	cubeMesh    *Mesh

	lightBoxModelLoc    int32
	lightBoxEmissiveLoc int32

	stencilShader *shaders.Stencil

	cubeMap uint32
	skybox  *shaders.Skybox
}

func (s *Scene) Init() {
	gl.Disable(gl.FRAMEBUFFER_SRGB)
	s.stencilShader = shaders.NewStencil()
	s.dirLightShader = shaders.NewDirectionalLight()
	s.pointLightShader = shaders.NewPointLightShader(maxPointLights)
	s.hdrShader = shaders.NewHDR()
	s.ssao = NewSSAO()
	s.hdr = NewHDRFBO()

	s.cubeMap = GetCubeMap()
	s.skybox = shaders.NewSkybox()

}

func (s *Scene) Render() {
	now := glfw.GetTime()
	s.elapsed = float32(now - s.previousTime)
	s.previousTime = now

	handleInputs()

	view := s.camera.View(s.elapsed)
	sin := float32(math.Sin(glfw.GetTime()))
	invView := view.Inv()

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

		s.graph.Render(s.projection, view, s.gBuffer.tShader, s.gBuffer.mShader)
	}

	// we have written all the depth buffer information we wanted.
	gl.DepthMask(false)

	// we wont be needing the depth tests for until we start a forward rendering again
	gl.Disable(gl.DEPTH_TEST)

	{ // screen space ambient occlusion (SSAO)
		gl.BindFramebuffer(gl.FRAMEBUFFER, s.ssao.fbo)
		gl.UseProgram(s.ssao.shader.Program)

		if ssaoOn {
			gl.Uniform1i(s.ssao.shader.LocEnabled, 1)
		} else {
			gl.Uniform1i(s.ssao.shader.LocEnabled, 0)
		}
		// Send kernel (samples)
		for i, sample := range s.ssao.Kernel {
			gl.Uniform3f(s.ssao.shader.LocSamples[i], sample[0], sample[1], sample[2])
		}
		gl.UniformMatrix4fv(s.ssao.shader.LocProjection, 1, false, &s.projection[0])
		gl.UniformMatrix4fv(s.ssao.shader.LocInvProjection, 1, false, &s.invProjection[0])
		gl.Uniform2f(s.ssao.shader.LocScreenSize, float32(windowWidth), float32(windowHeight))

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gDepth)
		gl.Uniform1i(s.ssao.shader.LocGDepth, 0)

		// see @todo in the shader
		//gl.ActiveTexture(gl.TEXTURE1)
		//gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gNormal)
		//gl.Uniform1i(s.ssao.shader.LocGNormal, 1)

		// see @todo in the shader
		//gl.ActiveTexture(gl.TEXTURE2)
		//gl.BindTexture(gl.TEXTURE_2D, s.ssao.noiseTexture)
		//gl.Uniform1i(s.ssao.shader.LocTexNoise, 2)

		renderQuad()
	}

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

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gDepth)
		gl.Uniform1i(s.pointLightShader.LocGDepth, 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gNormal)
		gl.Uniform1i(s.pointLightShader.LocGNormal, 1)

		gl.ActiveTexture(gl.TEXTURE2)
		gl.Uniform1i(s.pointLightShader.LocGAlbedo, 2)
		gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gAlbedoSpec)

		gl.ActiveTexture(gl.TEXTURE3)
		gl.Uniform1i(s.pointLightShader.LocGAmbientOcclusion, 3)
		gl.BindTexture(gl.TEXTURE_2D, s.ssao.texture)

		gl.UniformMatrix4fv(s.pointLightShader.LocProjMatrixInv, 1, false, &s.invProjection[0])
		gl.UniformMatrix4fv(s.pointLightShader.LocViewMatrixInv, 1, false, &invView[0])
		gl.Uniform2f(s.pointLightShader.LocScreenSize, float32(windowWidth), float32(windowHeight))

		gl.Uniform3fv(s.pointLightShader.LocViewPos, 1, &s.camera.position[0])

		renderQuad()
	}

	if dirLightOn { // Render the directional light, for now there is only one
		gl.UseProgram(s.dirLightShader.Program)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(s.dirLightShader.LocGDepth, 0)
		gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gDepth)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.Uniform1i(s.dirLightShader.LocGNormal, 1)
		gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gNormal)

		gl.ActiveTexture(gl.TEXTURE2)
		gl.Uniform1i(s.dirLightShader.LocGAlbedo, 2)
		gl.BindTexture(gl.TEXTURE_2D, s.gBuffer.buffer.gAlbedoSpec)

		gl.ActiveTexture(gl.TEXTURE3)
		gl.Uniform1i(s.dirLightShader.LocShadowMap, 3)
		gl.BindTexture(gl.TEXTURE_2D, s.shadow.depthMap)

		gl.ActiveTexture(gl.TEXTURE4)
		gl.Uniform1i(s.dirLightShader.LocGAmbientOcclusion, 4)
		gl.BindTexture(gl.TEXTURE_2D, s.ssao.texture)

		gl.UniformMatrix4fv(s.dirLightShader.LocProjMatrixInv, 1, false, &s.invProjection[0])
		gl.UniformMatrix4fv(s.dirLightShader.LocViewMatrixInv, 1, false, &invView[0])
		gl.UniformMatrix4fv(s.dirLightShader.LocLightProjection, 1, false, &lightProjection[0])
		gl.UniformMatrix4fv(s.dirLightShader.LocLightView, 1, false, &lightView[0])
		gl.Uniform3fv(s.dirLightShader.LocLightDirection, 1, &directionLight.Direction[0])
		gl.Uniform3fv(s.dirLightShader.LocLightColor, 1, &directionLight.Color[0])
		gl.Uniform3fv(s.dirLightShader.LocViewPos, 1, &s.camera.position[0])
		gl.Uniform2f(s.dirLightShader.LocScreenSize, float32(windowWidth), float32(windowHeight))

		renderQuad()
	}
	gl.Disable(gl.BLEND)

	// start a forward rendering pass from here
	gl.Enable(gl.DEPTH_TEST)

	// draw sky box, this could probably coded in a way that is more performant
	if dirLightOn {
		gl.DepthFunc(gl.LEQUAL)
		gl.UseProgram(s.skybox.Program)
		gl.UniformMatrix4fv(s.skybox.LocProjection, 1, false, &s.projection[0])
		skyboxView := view.Mat3().Mat4()
		gl.UniformMatrix4fv(s.skybox.LocView, 1, false, &skyboxView[0])
		gl.BindVertexArray(s.skybox.SkyboxVAO)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(s.skybox.LocScreenTexture, 0)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, s.cubeMap)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}

	{ // render emissive objects
		s.lightBoxShader.UsePV(s.projection, view)

		for i := range s.pointLights[:currentNumLights] {
			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1]+sin, s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(0.03, 0.03, 0.03))
			model = model.Mul4(mgl32.HomogRotate3D(float32(math.Cos(glfw.GetTime())), mgl32.Vec3{1, 1, 1}.Normalize()))

			gl.UniformMatrix4fv(s.lightBoxModelLoc, 1, false, &model[0])
			gl.Uniform3fv(s.lightBoxEmissiveLoc, 1, &s.pointLights[i].Color[0])

			gl.BindVertexArray(s.cubeMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))
		}
	}

	// Here we start the full screen effects, so we don't want any depth test shenanigans going on
	gl.Disable(gl.DEPTH_TEST)

	out := s.gBuffer.buffer.finalTexture
	if bloomOn {
		out = s.bloom.Render(s.gBuffer.buffer.finalTexture)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, s.hdr.fbo)
	gl.UseProgram(s.hdrShader.Program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(s.hdrShader.LocScreenTexture, 0)
	gl.BindTexture(gl.TEXTURE_2D, out)
	renderQuad()
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// do the final rendering to the backbuffer
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	// taking care of retina having more actual pixels
	gl.Viewport(0, 0, viewPortWidth, viewPortHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	fxaaShader.Use()
	if fxaaOn {
		gl.Uniform1i(fxaaLoc_enabled, 1)
	} else {
		gl.Uniform1i(fxaaLoc_enabled, 0)
	}
	if showDebug {
		gl.Uniform1i(fxaaLocU_showEdges, 1)
	} else {
		gl.Uniform1i(fxaaLocU_showEdges, 0)
	}
	gl.Uniform1f(fxaaLocU_lumaThreshold, u_lumaThreshold)
	gl.Uniform1f(fxaaLocU_minReduce, u_minReduce)
	gl.Uniform1f(fxaaLocU_mulReduce, u_mulReduce)
	gl.Uniform1f(fxaaLocU_maxSpan, u_maxSpan)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(fxaaTextureloc, 0)
	gl.BindTexture(gl.TEXTURE_2D, s.hdr.texture)
	renderQuad()

	// and if debug is on, quad print them on top off everything
	if showDebug {
		DisplayNormalBufferTexture(s.gBuffer.buffer.gNormal)
		DisplayColorTexBuffer(s.gBuffer.buffer.gAlbedoSpec)
		DisplayDepthbufferTexture(s.gBuffer.buffer.gDepth)
		//DisplayDepthbufferTexture(s.ssao.texture)
		//DisplayDepthbufferTexture(s.shadow.depthMap)
	}

	// check if there was an opengl error in this frame, in that case panic
	chkError("end_of_frame")
}

func handleInputs() {
	if keys[glfw.Key1] {
		currentNumLights = 1
	} else if keys[glfw.Key2] {
		currentNumLights = 2
	} else if keys[glfw.Key3] {
		currentNumLights = 4
	} else if keys[glfw.Key4] {
		currentNumLights = 8
	} else if keys[glfw.Key5] {
		currentNumLights = 16
	} else if keys[glfw.Key6] {
		currentNumLights = 32
	} else if keys[glfw.Key0] {
		dirLightOn = false
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
		u_lumaThreshold = 1 / 16.0
		u_mulReduce = 1 / 8.0
		u_minReduce = 1 / 128.0
		u_maxSpan = 8

	} else if keys[glfw.KeyH] {
		u_lumaThreshold -= 0.05
		fmt.Println("threshold", u_lumaThreshold)
	} else if keys[glfw.KeyY] {
		u_lumaThreshold += 0.05
		fmt.Println("threshold", u_lumaThreshold)
	} else if keys[glfw.KeyJ] {
		u_mulReduce *= 2.0
		fmt.Println("mulReduce", u_mulReduce)
	} else if keys[glfw.KeyU] {
		u_mulReduce /= 2.0
		fmt.Println("mulReduce", u_mulReduce)
	} else if keys[glfw.KeyK] {
		u_minReduce *= 2.0
		fmt.Println("minreduce", u_minReduce)
	} else if keys[glfw.KeyI] {
		u_minReduce /= 2.0
		fmt.Println("minreduce", u_minReduce)
	} else if keys[glfw.KeyL] {
		u_maxSpan -= 1.0
		fmt.Println("maxSpan", u_maxSpan)
	} else if keys[glfw.KeyO] {
		u_maxSpan += 1.0
		fmt.Println("maxSpan", u_maxSpan)
	}

}
