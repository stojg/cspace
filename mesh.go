package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type mesh struct {
	Vertices []float32
	Textures []uint32
	vbo, vao uint32
}

func newCube(x, y, z float32) *mesh {
	q := &mesh{
		Vertices: []float32{
			-0.5, -0.5, -0.5, 0.0, 0.0,
			0.5, -0.5, -0.5, 1.0, 0.0,
			0.5, 0.5, -0.5, 1.0, 1.0,
			0.5, 0.5, -0.5, 1.0, 1.0,
			-0.5, 0.5, -0.5, 0.0, 1.0,
			-0.5, -0.5, -0.5, 0.0, 0.0,

			-0.5, -0.5, 0.5, 0.0, 0.0,
			0.5, -0.5, 0.5, 1.0, 0.0,
			0.5, 0.5, 0.5, 1.0, 1.0,
			0.5, 0.5, 0.5, 1.0, 1.0,
			-0.5, 0.5, 0.5, 0.0, 1.0,
			-0.5, -0.5, 0.5, 0.0, 0.0,

			-0.5, 0.5, 0.5, 1.0, 0.0,
			-0.5, 0.5, -0.5, 1.0, 1.0,
			-0.5, -0.5, -0.5, 0.0, 1.0,
			-0.5, -0.5, -0.5, 0.0, 1.0,
			-0.5, -0.5, 0.5, 0.0, 0.0,
			-0.5, 0.5, 0.5, 1.0, 0.0,

			0.5, 0.5, 0.5, 1.0, 0.0,
			0.5, 0.5, -0.5, 1.0, 1.0,
			0.5, -0.5, -0.5, 0.0, 1.0,
			0.5, -0.5, -0.5, 0.0, 1.0,
			0.5, -0.5, 0.5, 0.0, 0.0,
			0.5, 0.5, 0.5, 1.0, 0.0,

			-0.5, -0.5, -0.5, 0.0, 1.0,
			0.5, -0.5, -0.5, 1.0, 1.0,
			0.5, -0.5, 0.5, 1.0, 0.0,
			0.5, -0.5, 0.5, 1.0, 0.0,
			-0.5, -0.5, 0.5, 0.0, 0.0,
			-0.5, -0.5, -0.5, 0.0, 1.0,

			-0.5, 0.5, -0.5, 0.0, 1.0,
			0.5, 0.5, -0.5, 1.0, 1.0,
			0.5, 0.5, 0.5, 1.0, 0.0,
			0.5, 0.5, 0.5, 1.0, 0.0,
			-0.5, 0.5, 0.5, 0.0, 0.0,
			-0.5, 0.5, -0.5, 0.0, 1.0,
		},
	}
	q.init()
	return q
}

func (s *mesh) init() {
	const sizeOfFloat = 4

	// Create buffers/arrays
	gl.GenVertexArrays(1, &s.vao)
	gl.GenBuffers(1, &s.vbo)
	//gl.GenBuffers(1, s.ebo);

	gl.BindVertexArray(s.vao)

	// load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(s.Vertices)*sizeOfFloat, gl.Ptr(s.Vertices), gl.STATIC_DRAW)

	// vertex position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*sizeOfFloat, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// vertex texture coordinates
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*sizeOfFloat, gl.PtrOffset(3*sizeOfFloat))
	gl.EnableVertexAttribArray(1)
	// ensure that no one else by
	gl.BindVertexArray(0)
}

func (s mesh) Draw(shader *Shader) {

	// textures
	for i := range s.Textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, s.Textures[i])
		gl.Uniform1i(uniformLocation(shader, fmt.Sprintf("texture_diffuse%d", i+1)), int32(i))
	}

	// draw mesh
	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.Vertices)))

	// set back defaults, good practice stuff
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	//trans := mgl32.Translate3D(s.Position[0], s.Position[1], s.Position[2])
	////trans = trans.Mul4(mgl32.HomogRotate3D(45.0, mgl32.Vec3{0.0, 0.0, 1.0}))
	//trans = trans.Mul4(mgl32.Scale3D(s.Scale[0], s.Scale[1], s.Scale[2]))
	//transformLoc := uniformLocation(shader, "transform")
	//gl.UniformMatrix4fv(transformLoc, 1, false, &trans[0])

}

type pos [3]float32

func newSquare(index uint32, a, b pos) *square {
	s := &square{
		vbo: &index,
		vao: &index,
		vertices: []float32{
			a[0], b[1], 0.0,
			a[0], a[1], 0.0,
			b[0], a[1], 0.0,

			b[0], a[1], 0.0,
			b[0], b[1], 0.0,
			a[0], b[1], 0.0,
		},
	}
	s.init()
	return s
}

type square struct {
	vbo, vao *uint32
	vertices []float32
}

func (s *square) init() {
	gl.GenBuffers(1, s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(s.vertices)*4, gl.Ptr(s.vertices), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, s.vao)
	gl.BindVertexArray(*s.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, *s.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)
}

func (s square) Draw() {
	gl.BindVertexArray(*s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.vertices)))

}
