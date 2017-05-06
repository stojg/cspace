package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func NewScene(WindowWidth, WindowHeight int32) *Scene {

	s, err := NewShader("material", "material")
	if err != nil {
		panic(err)
	}

	origin := mgl32.Translate3D(0, 0, 0)

	return &Scene{
		previousTime: glfw.GetTime(),
		camera:       NewCamera(),
		gbuffer:      NewGbuffer(WindowWidth, WindowHeight),
		projection:   mgl32.Perspective(mgl32.DegToRad(67.0), float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0),
		graph: &Node{
			transform: &origin,
		},
		simple: s,
	}
}

type Scene struct {
	previousTime float64
	elapsed      float32
	projection   mgl32.Mat4
	camera       *Camera
	graph        *Node
	gbuffer      *Gbuffer
	simple       *Shader
}

func (s *Scene) Render() {
	s.updateTimers()
	s.camera.View(s.elapsed)

	view := s.camera.View(s.elapsed)
	s.simple.UsePV(s.projection, view)

	pos := [][]float32{{1, 0, 1}, {5, 2, 1}}
	colors := [][]float32{{1.000, 0.749, 0.000}, {1, 1, 1}}

	setLights(s.simple, pos, colors)
	SimpleRender(s)

	//DSGeometryPass(s)
	//DSLightPass()
}

func (s *Scene) updateTimers() {
	now := glfw.GetTime()
	s.elapsed = float32(now - s.previousTime)
	s.previousTime = now
}

func SimpleRender(s *Scene) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.graph.Render(s.simple)
}

func DSGeometryPass(s *Scene) {

	//m_DSGeomPassTech.Enable();

	s.gbuffer.BindForWriting()

	//glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

	//p := &Pipeline{}
	//p.Scale(0.1, 0.1, 0.1)
	//p.Rotate(0, 0, 0)
	//p.WorldPos(-0.8, -1, 12)
	//p.SetCamera(s.camera.position, s.camera.front, s.camera.up)
	//p.SetPerspectiveProj(s.projection)

	//m_DSGeomPassTech.SetWVP(p.GetWVPTrans());
	//m_DSGeomPassTech.SetWorldMatrix(p.GetWorldTrans());
	//m_mesh.Render();

}

type Pipeline struct {
	scale      [3]float32
	rotate     [3]float32
	worldPos   [3]float32
	camera     [3][3]float32
	projection mgl32.Mat4
}

func (p *Pipeline) Scale(x, y, z float32) {
	p.scale[0] = x
	p.scale[1] = y
	p.scale[2] = z
}

func (p *Pipeline) Rotate(x, y, z float32) {
	p.rotate[0] = x
	p.rotate[1] = y
	p.rotate[2] = z
}

func (p *Pipeline) WorldPos(x, y, z float32) {
	p.worldPos[0] = x
	p.worldPos[1] = y
	p.worldPos[2] = z
}

func (p *Pipeline) SetCamera(pos, front, up [3]float32) {
	p.camera[0] = pos
	p.camera[1] = front
	p.camera[2] = up
}

func (p *Pipeline) SetPerspectiveProj(proj mgl32.Mat4) {
	p.projection = proj

}
