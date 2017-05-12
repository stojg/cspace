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

var currentNumLights = 32

var directionLight = &DirectionalLight{
	Direction: normalise([3]float32{1, 1, 1}),
	Color:     [3]float32{0.1, 0.1, 0.1},
}

var passthroughShader *PassthroughShader
var bloomColShader *DefaultShader
var shaderBlur *DefaultShader
var bloomBlender *DefaultShader

func NewScene(WindowWidth, WindowHeight int32) *Scene {

	passthroughShader = NewPassthroughShader()
	bloomColShader = NewDefaultShader("fx", "fx_brigthness_sep")
	shaderBlur = NewDefaultShader("fx", "fx_guassian_blur")
	bloomBlender = NewDefaultShader("fx", "fx_bloom_blender")

	graphTransform := mgl32.Ident4()
	s := &Scene{
		gBufferPipeline: NewGBufferPipeline(),

		fxBuffer:         NewFXbuffer(),
		bloomEffect:      NewBloomEffect(),
		previousTime:     glfw.GetTime(),
		camera:           NewCamera(),
		projection:       mgl32.Perspective(mgl32.DegToRad(67.0), float32(WindowWidth)/float32(WindowHeight), 0.1, 200.0),
		graph:            &Node{transform: &graphTransform},
		pointLightShader: NewPointLightShader("lighting", "lighting"),
		dirLightShader:   NewDirLightShader("lighting", "dirlighting"),
		nullShader:       NewDefaultShader("null", "null"),
		lightBoxShader:   NewDefaultShader("simple", "emissive"),
		lightMesh:        newLightMesh(),
		icoMesh:          LoadModel("models/ico"),
	}

	for i := 0; i < numLights; i++ {
		att := ligthAtt[7]
		s.pointLights = append(s.pointLights, &PointLight{
			Position: [3]float32{rand.Float32()*30 - 15, 0, rand.Float32()*30 - 15},
			Color:    [3]float32{1 + rand.Float32()*2, 1 + rand.Float32()*2, 1 + rand.Float32()*2 + 0.5},
			Constant: att.Constant,
			Linear:   att.Linear,
			Exp:      att.Exp,
			rand:     rand.Float32() * 2,
			enabled:  true,
		})
	}

	// shader location caches
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
	graph        *Node

	fxBuffer        *FXFbo
	gBufferPipeline *GBufferPipeline
	bloomEffect     *BloomEffect

	pointLightShader *PointLightShader
	dirLightShader   *DirLightShader
	nullShader       *DefaultShader
	lightBoxShader   *DefaultShader

	pointLights []*PointLight
	lightMesh   *Mesh
	icoMesh     *Mesh

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
	view := s.camera.View(s.elapsed)
	sin := float32(math.Sin(glfw.GetTime()))

	handleInputs()

	s.gBufferPipeline.Render(s.projection, view, s.graph)

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
			gl.DrawBuffer(gl.NONE)
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

			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)

		}

		// 3. PointLight Pass
		{
			s.pointLightShader.UsePV(s.projection, view)

			gl.DrawBuffer(gl.COLOR_ATTACHMENT4)

			gl.ActiveTexture(gl.TEXTURE0)
			gl.Uniform1i(s.pointLightShader.UniformPosLoc(), 0)
			gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gPosition)

			gl.ActiveTexture(gl.TEXTURE1)
			gl.Uniform1i(s.pointLightShader.UniformNormalLoc(), 1)
			gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gNormal)

			gl.ActiveTexture(gl.TEXTURE2)
			gl.Uniform1i(s.pointLightShader.UniformAlbedoSpecLoc(), 2)
			gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gAlbedoSpec)

			gl.StencilFunc(gl.NOTEQUAL, 0, 0xFF)

			gl.Disable(gl.DEPTH_TEST)

			gl.Enable(gl.BLEND)
			gl.BlendEquation(gl.FUNC_ADD)
			gl.BlendFunc(gl.ONE, gl.ONE)

			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.FRONT)

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

			model := mgl32.Translate3D(s.pointLights[i].Position[0], cp.Position[1], s.pointLights[i].Position[2])
			rad := s.pointLights[i].Radius()
			model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
			gl.UniformMatrix4fv(s.pointLightModelLoc, 1, false, &model[0])

			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)

			//gl.BindVertexArray(s.icoMesh.vao)
			//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))
			//gl.BindVertexArray(0)

			gl.CullFace(gl.BACK)
			gl.Disable(gl.BLEND)
		}
	}

	gl.Disable(gl.STENCIL_TEST)

	{ // Render the directional term / ambient
		ident := mgl32.Ident4()
		s.dirLightShader.UsePV(ident, ident)

		gl.DrawBuffer(gl.COLOR_ATTACHMENT4)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(s.dirLightShader.UniformPosLoc(), 0)
		gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gPosition)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.Uniform1i(s.dirLightShader.UniformNormalLoc(), 1)
		gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gNormal)

		gl.ActiveTexture(gl.TEXTURE2)
		gl.Uniform1i(s.dirLightShader.UniformAlbedoSpecLoc(), 2)
		gl.BindTexture(gl.TEXTURE_2D, s.gBufferPipeline.buffer.gAlbedoSpec)

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

	{ // render emissive objects

		// light cubes
		s.lightBoxShader.UsePV(s.projection, view)
		gl.Enable(gl.DEPTH_TEST)
		for i := range s.pointLights[:currentNumLights] {
			if !s.pointLights[i].enabled {
				continue
			}
			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1]+sin+s.pointLights[i].rand, s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(0.05, 0.05, 0.05))

			gl.UniformMatrix4fv(s.lightBoxModelLoc, 1, false, &model[0])
			gl.Uniform3fv(s.lightBoxEmissiveLoc, 1, &s.pointLights[i].Color[0])

			gl.BindVertexArray(s.icoMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))
			gl.BindVertexArray(0)
		}

		// moon
		model := mgl32.Translate3D(100, 100, 100)
		model = model.Mul4(mgl32.Scale3D(8, 8, 8))
		setUniformMatrix4fv(s.lightBoxShader, "model", model)
		gl.Uniform3f(uniformLocation(s.lightBoxShader, "emissive"), 1.4, 1.4, 1.8)
		gl.BindVertexArray(s.icoMesh.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.icoMesh.Vertices)))
		gl.BindVertexArray(0)
	}

	// from here on, there are only texture manipulations
	gl.Disable(gl.DEPTH_TEST)

	out := s.gBufferPipeline.buffer.finalTexture
	if bloom {
		out = s.bloomEffect.Render(s.gBufferPipeline.buffer.finalTexture)
	}

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	passthroughShader.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(passthroughShader.uniformScreenTextureLoc, 0)
	gl.BindTexture(gl.TEXTURE_2D, out)
	renderQuad()
	chkError("end_of_frame")
}
func handleInputs() {
	if keys[glfw.Key1] {
		currentNumLights = 0
	} else if keys[glfw.Key2] {
		currentNumLights = 4
	} else if keys[glfw.Key3] {
		currentNumLights = 8
	} else if keys[glfw.Key4] {
		currentNumLights = 16
	} else if keys[glfw.Key5] {
		currentNumLights = 32
	} else if keys[glfw.Key6] {
		currentNumLights = 64
	} else if keys[glfw.Key7] {
		currentNumLights = 128
	} else if keys[glfw.Key8] {
		currentNumLights = 255
	} else if keys[glfw.KeyEnter] {
		bloom = true
	} else if keys[glfw.KeyEscape] {
		bloom = false
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
