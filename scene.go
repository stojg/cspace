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
		xPos := rand.Float32()*30 - 15
		yPos := rand.Float32() + float32(1.2)
		//zPos := rand.Float32()*5.0 - 2.5
		zPos := rand.Float32()*30 - 15
		s.lightPositions = append(s.lightPositions, [3]float32{xPos, yPos, zPos})
		// Also calculate random color
		rColor := rand.Float32()/2 + 0.5 // Between 0.5 and 1.0
		gColor := rand.Float32()/2 + 0.5
		bColor := rand.Float32()/2 + 0.5
		s.lightColors = append(s.lightColors, [3]float32{rColor, gColor, bColor})
	}

	s.lightMesh = newLightMesh()

	return s
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

type Scene struct {
	width, height  int32
	previousTime   float64
	elapsed        float32
	projection     mgl32.Mat4
	camera         *Camera
	graph          *Node
	gbuffer        *Gbuffer
	gBufferShader  *Shader
	shaderLighting *Shader
	shaderLightBox *Shader
	lightPositions [][3]float32
	lightColors    [][3]float32
	lightMesh      *Mesh
}

func (s *Scene) Render() {
	s.updateTimers()
	view := s.camera.View(s.elapsed)

	// render the gBuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.gbuffer.fbo)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	s.gBufferShader.UsePV(s.projection, view)
	s.graph.Render(s.gBufferShader)

	// deferred pass
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

	setLights(s.shaderLighting, s.lightPositions, s.lightColors)
	loc := gl.GetUniformLocation(s.shaderLighting.Program, gl.Str("viewPos\x00"))
	if loc < 0 {
		panic("oh noes")
	}
	gl.Uniform3fv(loc, 1, &s.camera.position[0])
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
	for i := 0; i < numLights; i++ {
		model := mgl32.Translate3D(s.lightPositions[i][0], s.lightPositions[i][1], s.lightPositions[i][2])
		model = model.Mul4(mgl32.Scale3D(0.05, 0.05, 0.05))
		setUniformMatrix4fv(s.shaderLightBox, "model", model)
		gl.Uniform3f(uniformLocation(s.shaderLightBox, "emissive"), s.lightColors[i][0], s.lightColors[i][1], s.lightColors[i][2])
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
