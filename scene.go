package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const numLights = 24

var screenShader *DefaultShader

func NewScene(WindowWidth, WindowHeight int32) *Scene {

	screenShader = NewDefaultShader("fx", "fx_pass")

	graphTransform := mgl32.Ident4()
	s := &Scene{
		gbuffer:          NewGbuffer(WindowWidth, WindowHeight),
		fxFBO:            NewFXFbo(),
		width:            WindowWidth,
		height:           WindowHeight,
		previousTime:     glfw.GetTime(),
		camera:           NewCamera(),
		projection:       mgl32.Perspective(mgl32.DegToRad(67.0), float32(WindowWidth)/float32(WindowHeight), 0.1, 200.0),
		graph:            &Node{transform: &graphTransform},
		gBufferShader:    NewDefaultShader("gbuffer", "gbuffer"),
		pointLightShader: NewPointLightShader("lighting", "lighting"),
		dirLightShader:   NewDirLightShader("lighting", "dirlighting"),
		nullShader:       NewDefaultShader("null", "null"),
		lightBoxShader:   NewDefaultShader("simple", "emissive"),
		lightMesh:        newLightMesh(),
	}

	for i := 0; i < numLights; i++ {
		att := ligthAtt[7]
		s.pointLights = append(s.pointLights, &PointLight{
			Position: [3]float32{rand.Float32()*30 - 15, 1, rand.Float32()*30 - 15},
			Color:    [3]float32{rand.Float32()/2 + 0.5, rand.Float32()/2 + 0.5, rand.Float32()/2 + 0.5},
			Constant: att.Constant,
			Linear:   att.Linear,
			Exp:      att.Exp,
			rand:     rand.Float32(),
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
	return s
}

type Scene struct {
	width, height int32
	previousTime  float64
	elapsed       float32
	projection    mgl32.Mat4
	camera        *Camera
	graph         *Node

	fxFBO   *FXFbo
	gbuffer *Gbuffer

	gBufferShader    *DefaultShader
	pointLightShader *PointLightShader
	dirLightShader   *DirLightShader
	nullShader       *DefaultShader
	lightBoxShader   *DefaultShader

	pointLights []*PointLight
	lightMesh   *Mesh

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

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, s.gbuffer.fbo)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT4)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// 1. render into the gBuffer
	{
		s.gBufferShader.UsePV(s.projection, view)

		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, s.gbuffer.fbo)
		var attachments = [3]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2}
		gl.DrawBuffers(3, &attachments[0])

		// Only the geometry pass updates the gDepth buffer
		gl.DepthMask(true)
		gl.ClearColor(0.0, 0.0, 0.0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)

		s.graph.Render(s.gBufferShader)
		// When we get here the gDepth buffer is already populated and the stencil pass depends on it, but it does not write to it.
		gl.DepthMask(false)
	}

	// We need stencil to be enabled in the stencil pass to get the stencil buffer updated and we also need it in the
	// light pass because we render the light only if the stencil passes.
	gl.Enable(gl.STENCIL_TEST)

	var rad float32 = 14

	for i := range s.pointLights {
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
			model = model.Mul4(mgl32.Scale3D(rad, rad, rad))

			gl.UniformMatrix4fv(s.nullShaderModelLoc, 1, false, &model[0])

			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)

		}

		// 3. PointLight Pass
		{
			s.pointLightShader.UsePV(s.projection, view)
			s.gbuffer.BindForLightPass(s.pointLightShader)

			gl.StencilFunc(gl.NOTEQUAL, 0, 0xFF)

			gl.Disable(gl.DEPTH_TEST)

			gl.Enable(gl.BLEND)
			gl.BlendEquation(gl.FUNC_ADD)
			gl.BlendFunc(gl.ONE, gl.ONE)

			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.FRONT)

			gl.Uniform2f(s.pointLightShaderScreenSizeLoc, float32(s.width), float32(s.height))

			cp := PointLight{
				Position: s.pointLights[i].Position,
				Color:    s.pointLights[i].Color,
				Constant: s.pointLights[i].Constant,
				Linear:   s.pointLights[i].Linear,
				Exp:      s.pointLights[i].Exp,
				radius:   s.pointLights[i].radius,
			}
			cp.Position[1] += sin * s.pointLights[i].rand
			s.pointLightShader.SetLight(&cp)

			gl.Uniform3fv(s.pointLightViewPosLoc, 1, &s.camera.position[0])

			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1]+sin*s.pointLights[i].rand, s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
			gl.UniformMatrix4fv(s.pointLightModelLoc, 1, false, &model[0])

			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)

			gl.CullFace(gl.BACK)
			gl.Disable(gl.BLEND)
		}
	}
	// we don't want to use the stencil testing any more
	gl.Disable(gl.STENCIL_TEST)

	{ // Render the directional term / ambient
		directionLight := &DirectionalLight{
			Direction: normalise([3]float32{1, 1, 1}),
			Color:     [3]float32{0.2, 0.2, 0.4},
		}

		ident := mgl32.Ident4()
		s.dirLightShader.UsePV(ident, ident)
		s.gbuffer.BindForLightPass(s.dirLightShader)

		gl.Disable(gl.DEPTH_TEST)
		gl.Enable(gl.BLEND)
		gl.BlendEquation(gl.FUNC_ADD)
		gl.BlendFunc(gl.ONE, gl.ONE)

		s.dirLightShader.SetLight(directionLight)
		gl.Uniform2f(s.dirLightShaderScreenSizeLoc, float32(s.width), float32(s.height))
		gl.Uniform3fv(s.dirLightViewPosLoc, 1, &s.camera.position[0])
		gl.UniformMatrix4fv(s.dirLightModelLoc, 1, false, &ident[0])
		renderQuad()

		gl.Disable(gl.BLEND)
	}

	{ // render emissive light cubes
		s.lightBoxShader.UsePV(s.projection, view)
		gl.Enable(gl.DEPTH_TEST)

		for i := range s.pointLights {
			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1]+sin*s.pointLights[i].rand, s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(0.05, 0.05, 0.05))

			gl.UniformMatrix4fv(s.lightBoxModelLoc, 1, false, &model[0])
			gl.Uniform3fv(s.lightBoxEmissiveLoc, 1, &s.pointLights[i].Color[0])

			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)
		}
	}

	{ // moon
		model := mgl32.Translate3D(100, 100, 100)
		model = model.Mul4(mgl32.Scale3D(8, 8, 8))
		setUniformMatrix4fv(s.lightBoxShader, "model", model)
		gl.Uniform3f(uniformLocation(s.lightBoxShader, "emissive"), 0.7, 0.7, 0.9)

		gl.BindVertexArray(s.lightMesh.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
		gl.BindVertexArray(0)
	}

	// 4. final pass
	{
		s.gbuffer.BindForFinalPass(s.fxFBO.id)
		gl.BlitFramebuffer(0, 0, s.width, s.height, 0, 0, s.width, s.height, gl.COLOR_BUFFER_BIT, gl.LINEAR)
	}

	{ // post effect pass
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
		screenShader.Use()

		gl.Disable(gl.DEPTH_TEST)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(uniformLocation(screenShader, "screenTexture"), 0)
		gl.BindTexture(gl.TEXTURE_2D, s.fxFBO.textures[0])
	}

	renderQuad()

	chkError()
}

func (s *Scene) updateTimers() {
	now := glfw.GetTime()
	s.elapsed = float32(now - s.previousTime)
	s.previousTime = now
}

func chkError() {
	err := gl.GetError()
	if err == 0 {
		return
	}
	switch err {
	case gl.INVALID_OPERATION:
		fmt.Printf("GL Error: INVALID_OPERATION 0x0%x\n", err)
	case gl.INVALID_ENUM:
		fmt.Printf("GL Error: INVALID_ENUM 0x0%x\n", err)
	default:
		fmt.Printf("GL Error: 0x0%x\n", err)
	}
	panic("nope")
}
