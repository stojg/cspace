package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const numLights = 255

var bloom = false
var dirLightOn = true
var fxaa = false
var showDebug = false
var currentNumLights = 0

var directionLight = &DirectionalLight{
	Direction: normalise([3]float32{10, 10, 10}),
	Color:     [3]float32{1, 1, 1},
}

var passthroughShader *PassthroughShader
var depthShader *DefaultShader
var hdrShader *DefaultShader
var fxaaShader *DefaultShader
var fxaaTextureloc int32

func NewScene() *Scene {

	passthroughShader = NewPassthroughShader()
	hdrShader = NewDefaultShader("fx", "fx_tone")
	fxaaShader = NewDefaultShader("fx", "fx_fxaa")
	fxaaTextureloc = uniformLocation(fxaaShader, "screenTexture")

	s := &Scene{
		gBufferPipeline: NewGBufferPipeline(),

		fxBuffer:         NewFXbuffer(),
		bloomEffect:      NewBloomEffect(),
		previousTime:     glfw.GetTime(),
		camera:           NewCamera(),
		projection:       mgl32.Perspective(mgl32.DegToRad(67.0), float32(windowWidth)/float32(windowHeight), 0.5, 200.0),
		graph:            NewBaseNode(),
		pointLightShader: NewPointLightShader("lighting", "lighting_point_pbr"),
		dirLightShader:   NewDirLightShader("lighting", "lighting_dir_pbr"),
		nullShader:       NewDefaultShader("null", "null"),
		lightBoxShader:   NewDefaultShader("simple", "emissive"),
		icoMesh:          LoadModel("models/ico", MaterialMesh)[0],
		cubeMesh:         LoadModel("models/cube", MaterialMesh)[0],
	}

	att := ligthAtt[1]
	s.pointLights = append(s.pointLights, &PointLight{
		Position: [3]float32{12, 12, 10},
		Color:    [3]float32{1, 1, 1},
		Constant: att.Constant,
		Linear:   att.Linear,
		Exp:      att.Exp,
		rand:     rand.Float32() * 2,
		enabled:  true,
	})

	for i := 1; i < numLights; i++ {

		s.pointLights = append(s.pointLights, &PointLight{
			Position: [3]float32{rand.Float32() * 25, rand.Float32() * 25, 3},
			Color:    [3]float32{rand.Float32()*3 + 0.5, rand.Float32()*3 + 0.5, rand.Float32()*3 + 0.5},
			Constant: att.Constant,
			Linear:   att.Linear,
			Exp:      att.Exp,
			rand:     rand.Float32() * 2,
			enabled:  true,
		})
	}

	// tShader location caches
	s.nullShaderModelLoc = uniformLocation(s.nullShader, "model")
	s.pointLightShaderScreenSizeLoc = uniformLocation(s.pointLightShader, "gScreenSize")
	s.pointLightViewPosLoc = uniformLocation(s.pointLightShader, "viewPos")
	s.pointLightModelLoc = uniformLocation(s.pointLightShader, "model")

	s.dirLightShaderScreenSizeLoc = uniformLocation(s.dirLightShader, "gScreenSize")
	s.dirLightViewPosLoc = uniformLocation(s.dirLightShader, "viewPos")
	s.dirLightModelLoc = uniformLocation(s.dirLightShader, "model")

	s.lightBoxModelLoc = uniformLocation(s.lightBoxShader, "model")
	s.lightBoxEmissiveLoc = uniformLocation(s.lightBoxShader, "emissive")
	chkError("end_of_new_scene")
	return s
}

type Scene struct {
	previousTime float64
	elapsed      float32
	projection   mgl32.Mat4
	camera       *Camera
	graph        SceneNode

	fxBuffer        *FXFbo
	gBufferPipeline *GBufferPipeline
	bloomEffect     *BloomEffect

	pointLightShader *PointLightShader
	dirLightShader   *DirLightShader
	nullShader       *DefaultShader
	lightBoxShader   *DefaultShader

	pointLights []*PointLight
	icoMesh     *Mesh
	cubeMesh    *Mesh

	// caches - should move to the individual shaders
	nullShaderModelLoc            int32
	pointLightShaderScreenSizeLoc int32
	pointLightViewPosLoc          int32
	pointLightModelLoc            int32
	dirLightShaderScreenSizeLoc   int32
	dirLightViewPosLoc            int32
	dirLightModelLoc              int32
	lightBoxModelLoc              int32
	lightBoxEmissiveLoc           int32
}

func (s *Scene) Render() {
	s.updateTimers()
	handleInputs()
	view := s.camera.View(s.elapsed)
	sin := float32(math.Sin(glfw.GetTime()))
	invProj := s.projection.Inv()
	invView := view.Inv()

	// bind and clear out the gbuffers final texture
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, s.gBufferPipeline.buffer.fbo)

	// 1. render into the gBuffer
	var attachments = [2]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1}
	gl.DrawBuffers(int32(len(attachments)), &attachments[0])

	// Only the geometry pass updates the depth buffer
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// enable depth mask writing, draw, and disable writing to the depth buffer
	s.graph.Render(s.projection, view, s.gBufferPipeline.tShader, s.gBufferPipeline.mShader)
	gl.DepthMask(false)

	// All rendering should now go into the gbuffers final texture
	gl.DrawBuffer(gl.COLOR_ATTACHMENT4)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// We need stencil to be enabled in the stencil pass to get the stencil buffer updated and we also need it in the
	// light pass because we render the light only if the stencil passes.
	gl.Enable(gl.STENCIL_TEST)
	for i := range s.pointLights[:currentNumLights] {
		if !s.pointLights[i].enabled {
			continue
		}

		// 2. stencil pass
		{
			s.nullShader.UsePV(s.projection, view)

			// Disable color/depth write and enable depth testing
			gl.Enable(gl.DEPTH_TEST)

			// otherwise the light will be inside by the light bounding volume
			gl.Disable(gl.CULL_FACE)

			// clear previous runs stencil tests
			gl.Clear(gl.STENCIL_BUFFER_BIT)

			// We need the stencil test to be enabled but we want it to succeed always. Only the depth test matters.
			gl.StencilFunc(gl.ALWAYS, 0, 0)
			// The following stencil setup that we saw guarantees that only the pixels in the stencil buffer covered by
			// objects inside the bounding sphere will have a value greater than zero
			gl.StencilOpSeparate(gl.BACK, gl.KEEP, gl.INCR_WRAP, gl.KEEP)
			gl.StencilOpSeparate(gl.FRONT, gl.KEEP, gl.DECR_WRAP, gl.KEEP)

			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			rad := s.pointLights[i].Radius()
			model = model.Mul4(mgl32.Scale3D(rad, rad, rad))

			gl.UniformMatrix4fv(s.nullShaderModelLoc, 1, false, &model[0])

			gl.BindVertexArray(s.icoMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))

		}

		// 3. PointLight Pass
		{
			s.pointLightShader.UsePV(s.projection, view)

			gl.Disable(gl.DEPTH_TEST)
			gl.StencilFunc(gl.NOTEQUAL, 0, 0xFF)

			gl.Enable(gl.BLEND)
			gl.BlendEquation(gl.FUNC_ADD)
			gl.BlendFunc(gl.ONE, gl.ONE)

			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.FRONT)

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gDepth)
			gl.Uniform1i(s.pointLightShader.UniformDepthLoc(), 0)

			gl.ActiveTexture(gl.TEXTURE1)
			gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gNormal)
			gl.Uniform1i(s.pointLightShader.UniformNormalLoc(), 1)

			gl.ActiveTexture(gl.TEXTURE2)
			gl.Uniform1i(s.pointLightShader.UniformAlbedoSpecLoc(), 2)
			gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gAlbedoSpec)

			gl.UniformMatrix4fv(s.pointLightShader.locProjMatrixInv, 1, false, &invProj[0])
			gl.UniformMatrix4fv(s.pointLightShader.locViewMatrixInv, 1, false, &invView[0])

			gl.Uniform2f(s.pointLightShaderScreenSizeLoc, float32(windowWidth), float32(windowHeight))

			cp := PointLight{
				Position: s.pointLights[i].Position,
				Color:    s.pointLights[i].Color,
				Constant: s.pointLights[i].Constant,
				Linear:   s.pointLights[i].Linear,
				Exp:      s.pointLights[i].Exp,
				radius:   s.pointLights[i].radius,
			}
			cp.Position[1] += sin + s.pointLights[i].rand
			s.pointLightShader.SetLight(&cp)

			gl.Uniform3fv(s.pointLightViewPosLoc, 1, &s.camera.position[0])

			model := mgl32.Translate3D(cp.Position[0], cp.Position[1], cp.Position[2])
			rad := s.pointLights[i].Radius()
			model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
			gl.UniformMatrix4fv(s.pointLightModelLoc, 1, false, &model[0])

			gl.BindVertexArray(s.icoMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))

			gl.CullFace(gl.BACK)
			gl.Disable(gl.BLEND)
		}
	}
	// we are done with the stencil testing
	gl.Disable(gl.STENCIL_TEST)

	if dirLightOn { // Render the directional term / ambient
		ident := mgl32.Ident4()
		s.dirLightShader.UsePV(ident, ident)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(s.dirLightShader.UniformDepthLoc(), 0)
		gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gDepth)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.Uniform1i(s.dirLightShader.UniformNormalLoc(), 1)
		gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gNormal)

		gl.ActiveTexture(gl.TEXTURE2)
		gl.Uniform1i(s.dirLightShader.UniformAlbedoSpecLoc(), 2)
		gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gAlbedoSpec)

		gl.UniformMatrix4fv(s.dirLightShader.locProjMatrixInv, 1, false, &invProj[0])
		gl.UniformMatrix4fv(s.dirLightShader.locViewMatrixInv, 1, false, &invView[0])

		gl.Disable(gl.DEPTH_TEST)
		gl.Enable(gl.BLEND)
		gl.BlendEquation(gl.FUNC_ADD)
		gl.BlendFunc(gl.ONE, gl.ONE)

		s.dirLightShader.SetLight(directionLight)
		gl.Uniform2f(s.dirLightShaderScreenSizeLoc, float32(windowWidth), float32(windowHeight))
		gl.Uniform3fv(s.dirLightViewPosLoc, 1, &s.camera.position[0])
		gl.UniformMatrix4fv(s.dirLightModelLoc, 1, false, &ident[0])
		renderQuad()

		gl.Disable(gl.BLEND)
	}

	gl.Enable(gl.DEPTH_TEST)

	{ // render emissive objects
		// light cubes
		s.lightBoxShader.UsePV(s.projection, view)

		for i := range s.pointLights[:currentNumLights] {
			if !s.pointLights[i].enabled {
				continue
			}
			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1]+sin+s.pointLights[i].rand, s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(0.03, 0.03, 0.03))
			model = model.Mul4(mgl32.HomogRotate3D(float32(math.Cos(glfw.GetTime())), mgl32.Vec3{1, 1, 1}.Normalize()))

			gl.UniformMatrix4fv(s.lightBoxModelLoc, 1, false, &model[0])
			gl.Uniform3fv(s.lightBoxEmissiveLoc, 1, &s.pointLights[i].Color[0])

			gl.BindVertexArray(s.cubeMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))
		}

		// moon
		model := mgl32.Translate3D(100, 100, 100)
		model = model.Mul4(mgl32.Scale3D(8, 8, 8))
		gl.UniformMatrix4fv(s.lightBoxModelLoc, 1, false, &model[0])
		gl.Uniform3f(s.lightBoxEmissiveLoc, 1.4, 1.4, 1.8)
		gl.BindVertexArray(s.icoMesh.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))
	}

	// from here on, there are only texture manipulations
	gl.Disable(gl.DEPTH_TEST)

	out := s.gBufferPipeline.buffer.finalTexture
	if bloom {
		out = s.bloomEffect.Render(s.gBufferPipeline.buffer.finalTexture)
	}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	if fxaa {
		fxaaShader.Use()
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(fxaaTextureloc, 0)
		gl.BindTexture(gl.TEXTURE_2D, out)
	} else {
		hdrShader.Use()
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(passthroughShader.uniformScreenTextureLoc, 0)
		gl.BindTexture(gl.TEXTURE_2D, out)
	}
	renderQuad()

	if showDebug {
		//DisplayAlbedoBufferTexture(s.bloomEffect.bloomFbo.textures[1])
		DisplayAlbedoBufferTexture(s.gBufferPipeline.buffer.gAlbedoSpec)
		DisplayNormalBufferTexture(s.gBufferPipeline.buffer.gNormal)
		DisplayDepthbufferTexture(s.gBufferPipeline.buffer.gDepth)
	}

	chkError("end_of_frame")
}
func handleInputs() {
	if keys[glfw.Key1] {
		currentNumLights = 4
	} else if keys[glfw.Key2] {
		currentNumLights = 8
	} else if keys[glfw.Key3] {
		currentNumLights = 16
	} else if keys[glfw.Key4] {
		currentNumLights = 32
	} else if keys[glfw.Key5] {
		currentNumLights = 64
	} else if keys[glfw.Key6] {
		currentNumLights = 128
	} else if keys[glfw.Key0] {
		dirLightOn = false
	} else if keys[glfw.KeyF] {
		fxaa = true
	} else if keys[glfw.KeyTab] {
		showDebug = true
	} else if keys[glfw.KeyEnter] {
		bloom = true
	} else if keys[glfw.KeyEscape] {
		currentNumLights = 0
		bloom = false
		showDebug = false
		dirLightOn = true
		fxaa = false
	}
}

func (s *Scene) updateTimers() {
	now := glfw.GetTime()
	s.elapsed = float32(now - s.previousTime)
	s.previousTime = now
}

func chkError(name string) {
	err := gl.GetError()
	if err == 0 {
		return
	}
	switch err {
	case gl.INVALID_OPERATION:
		fmt.Printf("GL Error: INVALID_OPERATION 0x0%x\n", err)
	case gl.INVALID_ENUM:
		fmt.Printf("GL Error: INVALID_ENUM 0x0%x\n", err)
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		fmt.Printf("GL Error: INVALID_FRAMEBUFFER_OPERATION 0x0%x\n", err)
	default:
		fmt.Printf("GL Error: 0x0%x\n", err)
	}
	panic(name)
}
