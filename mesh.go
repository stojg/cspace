package main

import "github.com/go-gl/gl/v4.1-core/gl"

type mesh struct {
	Vertices []float32
	Textures []uint32
	vbo, vao uint32
}

func newCube(x, y, z float32) *mesh {
	q := &mesh{
		Vertices: []float32{
			-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
			0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
			0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
			0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
			-0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
			-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,

			-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
			0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
			0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
			0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
			-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
			-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,

			-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,
			-0.5, 0.5, -0.5, -1.0, 0.0, 0.0,
			-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
			-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
			-0.5, -0.5, 0.5, -1.0, 0.0, 0.0,
			-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,

			0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
			0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
			0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
			0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
			0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
			0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

			-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
			0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
			0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
			0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
			-0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
			-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,

			-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
			0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
			0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
			0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
			-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
			-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*sizeOfFloat, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*sizeOfFloat, gl.PtrOffset(3*sizeOfFloat))
	gl.EnableVertexAttribArray(1)

	// vertex texture coordinates
	//gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 6*sizeOfFloat, gl.PtrOffset(3*sizeOfFloat))
	//gl.EnableVertexAttribArray(1)
	// ensure that no one else by
	gl.BindVertexArray(0)
}
