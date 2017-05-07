package main

import (
	"fmt"
	"math/rand"

	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const numLights = 32

func NewScene(WindowWidth, WindowHeight int32) *Scene {

	graphTransform := mgl32.Ident4()
	s := &Scene{
		width:            WindowWidth,
		height:           WindowHeight,
		previousTime:     glfw.GetTime(),
		camera:           NewCamera(),
		gbuffer:          NewGbuffer(WindowWidth, WindowHeight),
		projection:       mgl32.Perspective(mgl32.DegToRad(67.0), float32(WindowWidth)/float32(WindowHeight), 0.1, 200.0),
		graph:            &Node{transform: &graphTransform},
		gBufferShader:    NewDefaultShader("gbuffer", "gbuffer"),
		pointLightShader: NewPointLightShader("lighting", "lighting"),
		dirLightShader:   NewDirLightShader("lighting", "dirlighting"),
		nullShader:       NewDefaultShader("null", "null"),
		lightBoxShader:   NewDefaultShader("simple", "emissive"),
		lightMesh:        newLightMesh(),
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < numLights; i++ {
		att := ligthAtt[13]
		s.pointLights = append(s.pointLights, &PointLight{
			Position: [3]float32{rand.Float32()*30 - 15, 1.2, rand.Float32()*30 - 15},
			Color:    [3]float32{rand.Float32()/2 + 0.5, rand.Float32()/2 + 0.5, rand.Float32()/2 + 0.5},
			Constant: att.Constant,
			Linear:   att.Linear,
			Exp:      att.Exp,
		})
	}
	return s
}

type Scene struct {
	width, height int32
	previousTime  float64
	elapsed       float32
	projection    mgl32.Mat4
	camera        *Camera
	graph         *Node

	gbuffer *Gbuffer

	gBufferShader    *DefaultShader
	pointLightShader *PointLightShader
	dirLightShader   *DirLightShader
	nullShader       *DefaultShader
	lightBoxShader   *DefaultShader

	pointLights []*PointLight
	lightMesh   *Mesh
}

func (s *Scene) Render() {
	s.updateTimers()
	view := s.camera.View(s.elapsed)

	s.gbuffer.StartFrame()

	// 1. render the gBuffer
	{
		s.gBufferShader.UsePV(s.projection, view)

		s.gbuffer.BindForGeomPass()
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

	for i := range s.pointLights {
		// 2. stencil pass
		{
			s.nullShader.UsePV(s.projection, view)
			s.gbuffer.BindForStencilPass()

			gl.Enable(gl.DEPTH_TEST)
			gl.Disable(gl.CULL_FACE)

			gl.Clear(gl.STENCIL_BUFFER_BIT)

			// We need the stencil test to be enabled but we want it to succeed always. Only the gDepth test matters.
			gl.StencilFunc(gl.ALWAYS, 0, 0)
			gl.StencilOpSeparate(gl.BACK, gl.KEEP, gl.INCR_WRAP, gl.KEEP)
			gl.StencilOpSeparate(gl.FRONT, gl.KEEP, gl.DECR_WRAP, gl.KEEP)

			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			rad := s.pointLights[i].Radius()
			model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
			setUniformMatrix4fv(s.nullShader, "model", model)
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

			gl.Uniform2f(uniformLocation(s.pointLightShader, "gScreenSize"), float32(s.width), float32(s.height))

			s.pointLightShader.SetLight(s.pointLights[i])

			gl.Uniform3fv(uniformLocation(s.pointLightShader, "viewPos"), 1, &s.camera.position[0])

			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			rad := s.pointLights[i].Radius()
			model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
			setUniformMatrix4fv(s.pointLightShader, "model", model)

			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)

			gl.CullFace(gl.BACK)
			gl.Disable(gl.BLEND)
		}
	}
	// we don't want to use the stencil testing any more
	gl.Disable(gl.STENCIL_TEST)

	{
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
		gl.Uniform2f(uniformLocation(s.dirLightShader, "gScreenSize"), float32(s.width), float32(s.height))
		gl.Uniform3fv(uniformLocation(s.dirLightShader, "viewPos"), 1, &s.camera.position[0])
		setUniformMatrix4fv(s.dirLightShader, "model", ident)
		renderQuad()

		gl.Disable(gl.BLEND)
	}

	s.lightBoxShader.UsePV(s.projection, view)
	gl.Enable(gl.DEPTH_TEST)

	for i := range s.pointLights {
		model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
		model = model.Mul4(mgl32.Scale3D(0.1, 0.1, 0.1))
		setUniformMatrix4fv(s.lightBoxShader, "model", model)
		gl.Uniform3fv(uniformLocation(s.lightBoxShader, "emissive"), 1, &s.pointLights[i].Color[0])

		gl.BindVertexArray(s.lightMesh.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
		gl.BindVertexArray(0)
	}

	model := mgl32.Translate3D(100, 100, 100)
	model = model.Mul4(mgl32.Scale3D(8, 8, 8))
	setUniformMatrix4fv(s.lightBoxShader, "model", model)
	gl.Uniform3f(uniformLocation(s.lightBoxShader, "emissive"), 0.7, 0.7, 0.9)

	gl.BindVertexArray(s.lightMesh.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
	gl.BindVertexArray(0)

	// 4. final pass
	{
		s.gbuffer.BindForFinalPass()
		gl.BlitFramebuffer(0, 0, s.width, s.height, 0, 0, s.width, s.height, gl.COLOR_BUFFER_BIT, gl.LINEAR)
	}

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
