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

	gShader, err := NewShader("gbuffer", "gbuffer")
	if err != nil {
		panic(err)
	}

	shaderLighting, err := NewShader("lighting", "lighting")
	if err != nil {
		panic(err)
	}

	shaderLightBox, err := NewShader("lightbox", "lightbox")
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
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < numLights; i++ {
		// Calculate slightly random offsets
		att := ligthAtt[7]
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

	shaderLighting *Shader

	shaderLightBox *Shader
	pointLights    []*PointLight
	lightMesh      *Mesh
}

func (s *Scene) Render() {
	s.updateTimers()
	view := s.camera.View(s.elapsed)

	// 1. render the gBuffer
	s.gbuffer.BindForWriting()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.gBufferShader.UsePV(s.projection, view)
	s.graph.Render(s.gBufferShader)

	// 2. deferred pass
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.ClearColor(0.1, 0.1, 0.1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.shaderLighting.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(uniformLocation(s.shaderLighting, "gPosition"), 0)
	gl.BindTexture(gl.TEXTURE_2D, s.gbuffer.gPosition)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.Uniform1i(uniformLocation(s.shaderLighting, "gNormal"), 1)
	gl.BindTexture(gl.TEXTURE_2D, s.gbuffer.gNormal)

	gl.ActiveTexture(gl.TEXTURE2)
	gl.Uniform1i(uniformLocation(s.shaderLighting, "gAlbedoSpec"), 2)
	gl.BindTexture(gl.TEXTURE_2D, s.gbuffer.gAlbedoSpec)

	setLights(s.shaderLighting, s.pointLights)
	gl.Uniform3fv(uniformLocation(s.shaderLighting, "viewPos"), 1, &s.camera.position[0])
	renderQuad()

	// 2.5. Copy content of geometry's depth buffer to default framebuffer's depth buffer
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, s.gbuffer.fbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0) // Write to default framebuffer
	//// blit to default framebuffer. Note that this may or may not work as the internal formats of both the FBO and default framebuffer have to match.
	//// the internal formats are implementation defined. This works on all of my systems, but if it doesn't on yours you'll likely have to write to the
	//// depth buffer in another shader stage (or somehow see to match the default framebuffer's internal format with the FBO's internal format).
	gl.BlitFramebuffer(0, 0, s.width, s.height, 0, 0, s.width, s.height, gl.DEPTH_BUFFER_BIT, gl.NEAREST)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// 3. Render lights on top of scene, by blitting
	s.shaderLightBox.UsePV(s.projection, view)

	for _, l := range s.pointLights {
		model := mgl32.Translate3D(l.Position[0], l.Position[1], l.Position[2])
		model = model.Mul4(mgl32.Scale3D(0.02, 0.02, 0.02))
		//rad := l.Radius()
		//model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
		setUniformMatrix4fv(s.shaderLightBox, "model", model)
		gl.Uniform3f(uniformLocation(s.shaderLightBox, "emissive"), l.Color[0], l.Color[1], l.Color[2])
		gl.BindVertexArray(s.lightMesh.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
		gl.BindVertexArray(0)
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
