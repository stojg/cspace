package main

import (
	"fmt"
	"math/rand"

	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const numLights = 64

func NewScene(WindowWidth, WindowHeight int32) *Scene {

	gShader, err := NewShader("gbuffer", "gbuffer")
	if err != nil {
		panic(err)
	}

	shaderLighting, err := NewLightShader("lighting", "lighting")
	if err != nil {
		panic(err)
	}

	shaderLightBox, err := NewShader("lightbox", "lightbox")
	if err != nil {
		panic(err)
	}

	nullShader, err := NewShader("null", "null")
	if err != nil {
		panic(err)
	}

	origin := mgl32.Translate3D(0, 0, 0)

	s := &Scene{
		width:        WindowWidth,
		height:       WindowHeight,
		previousTime: glfw.GetTime(),
		camera:       NewCamera(),
		gbuffer:      NewGbuffer(WindowWidth, WindowHeight),
		projection:   mgl32.Perspective(mgl32.DegToRad(67.0), float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0),
		graph: &Node{
			transform: &origin,
		},
		gBufferShader:  gShader,
		shaderLighting: shaderLighting,
		shaderLightBox: shaderLightBox,
		nullShader:     nullShader,
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < numLights; i++ {
		// Calculate slightly random offsets
		att := ligthAtt[13]
		l := &PointLight{
			Constant:         att.Constant,
			Linear:           att.Linear,
			Exp:              att.Exp,
			DiffuseIntensity: 1.0,
		}
		l.Position[0] = rand.Float32()*30 - 15
		l.Position[1] = rand.Float32() + float32(1.2)
		l.Position[2] = rand.Float32()*30 - 15
		l.Color[0] = rand.Float32()/2 + 0.5 // Between 0.5 and 1.0
		l.Color[1] = rand.Float32()/2 + 0.5
		l.Color[2] = rand.Float32()/2 + 0.5
		s.pointLights = append(s.pointLights, l)
	}

	s.lightMesh = newLightMesh()
	return s
}

type Scene struct {
	width, height int32
	previousTime  float64
	elapsed       float32
	projection    mgl32.Mat4
	camera        *Camera
	graph         *Node

	gbuffer       *Gbuffer
	gBufferShader *Shader

	shaderLighting *LightShader

	shaderLightBox *Shader
	pointLights    []*PointLight
	lightMesh      *Mesh

	nullShader *Shader
}

func (s *Scene) Render() {
	s.updateTimers()
	view := s.camera.View(s.elapsed)

	s.gbuffer.StartFrame()

	// 1. render the gBuffer
	{
		s.gBufferShader.UsePV(s.projection, view)

		s.gbuffer.BindForGeomPass()
		// Only the geometry pass updates the depth buffer
		gl.DepthMask(true)
		gl.ClearColor(0.1, 0.1, 0.1, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)

		s.graph.Render(s.gBufferShader)
		// When we get here the depth buffer is already populated and the stencil pass depends on it, but it does not
		// write to it.
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

			// We need the stencil test to be enabled but we want it to succeed always. Only the depth test matters.
			gl.StencilFunc(gl.ALWAYS, 0, 0)
			gl.StencilOpSeparate(gl.BACK, gl.KEEP, gl.INCR_WRAP, gl.KEEP)
			gl.StencilOpSeparate(gl.FRONT, gl.KEEP, gl.DECR_WRAP, gl.KEEP)

			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(4, 4, 4))
			//rad := l.Radius()
			//model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
			setUniformMatrix4fv(s.nullShader, "model", model)
			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)

		}

		// 3. PointLight Pass
		{

			s.shaderLighting.UsePV(s.projection, view)
			s.gbuffer.BindForLightPass(s.shaderLighting)

			gl.StencilFunc(gl.NOTEQUAL, 0, 0xFF)

			gl.Disable(gl.DEPTH_TEST)

			gl.Enable(gl.BLEND)
			gl.BlendEquation(gl.FUNC_ADD)
			gl.BlendFunc(gl.ONE, gl.ONE)

			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.FRONT)

			//gl.ClearColor(0.1, 0.1, 0.1, 0)
			//gl.Clear(gl.COLOR_BUFFER_BIT)

			gl.Uniform2f(uniformLocation(s.shaderLighting, "gScreenSize"), float32(s.width), float32(s.height))
			gl.Uniform3f(uniformLocation(s.shaderLighting, "pointLight.Position"), s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			gl.Uniform3f(uniformLocation(s.shaderLighting, "pointLight.Color"), s.pointLights[i].Color[0], s.pointLights[i].Color[1], s.pointLights[i].Color[2])
			gl.Uniform1f(uniformLocation(s.shaderLighting, "pointLight.Radius"), s.pointLights[i].Radius())
			gl.Uniform1f(uniformLocation(s.shaderLighting, "pointLight.Linear"), s.pointLights[i].Linear)
			gl.Uniform1f(uniformLocation(s.shaderLighting, "pointLight.Quadratic"), s.pointLights[i].Exp)
			gl.Uniform1f(uniformLocation(s.shaderLighting, "pointLight.DiffuseIntensity"), s.pointLights[i].DiffuseIntensity)
			gl.Uniform3fv(uniformLocation(s.shaderLighting, "viewPos"), 1, &s.camera.position[0])

			model := mgl32.Translate3D(s.pointLights[i].Position[0], s.pointLights[i].Position[1], s.pointLights[i].Position[2])
			model = model.Mul4(mgl32.Scale3D(4, 4, 4))
			//rad := l.Radius()
			//model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
			setUniformMatrix4fv(s.shaderLighting, "model", model)

			gl.BindVertexArray(s.lightMesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
			gl.BindVertexArray(0)

			gl.CullFace(gl.BACK)
			gl.Disable(gl.BLEND)
		}
	}

	// The directional light does not need a stencil test because its volume
	// is unlimited and the final pass simply copies the texture.
	gl.Disable(gl.STENCIL_TEST)
	// 4. DSDirectionalLightPass();

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
