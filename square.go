package main

import "github.com/go-gl/gl/v4.1-core/gl"

type mesh struct {
	Vertices []float32
	Texture  uint32
	vbo, vao *uint32
}

func newCube(index uint32) *mesh {
	q := &mesh{
		vbo: &index,
		vao: &index,
	}
	q.Vertices = []float32{
		// Bottom
		//  X, Y, Z       R, G, B        U, V
		-1.0, -1.0, -1.0, 0.5, 0.8, 0.9, 0.0, 0.0,
		1.0, -1.0, -1.0, 0.5, 0.8, 0.9, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.5, 0.8, 0.9, 0.0, 1.0,
		1.0, -1.0, -1.0, 0.5, 0.8, 0.9, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.5, 0.8, 0.9, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.5, 0.8, 0.9, 0.0, 1.0,

		// Top
		-1.0, 1.0, -1.0, 0.5, 0.8, 0.9, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.5, 0.8, 0.9, 0.0, 1.0,
		1.0, 1.0, -1.0, 0.5, 0.8, 0.9, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.5, 0.8, 0.9, 1.0, 0.0,
		-1.0, 1.0, 1.0, 0.5, 0.8, 0.9, 0.0, 1.0,
		1.0, 1.0, 1.0, 0.5, 0.8, 0.9, 1.0, 1.0,

		// Front
		-1.0, -1.0, 1.0, 0.5, 0.8, 0.9, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.5, 0.8, 0.9, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.5, 0.8, 0.9, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.5, 0.8, 0.9, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.5, 0.8, 0.9, 0.0, 1.0,
		-1.0, 1.0, 1.0, 0.5, 0.8, 0.9, 1.0, 1.0,

		// Back
		-1.0, -1.0, -1.0, 0.5, 0.8, 0.9, 0.0, 0.0,
		-1.0, 1.0, -1.0, 0.5, 0.8, 0.9, 0.0, 1.0,
		1.0, -1.0, -1.0, 0.5, 0.8, 0.9, 1.0, 0.0,
		1.0, -1.0, -1.0, 0.5, 0.8, 0.9, 1.0, 0.0,
		-1.0, 1.0, -1.0, 0.5, 0.8, 0.9, 0.0, 1.0,
		1.0, 1.0, -1.0, 0.5, 0.8, 0.9, 1.0, 1.0,

		// Left
		-1.0, -1.0, 1.0, 0.5, 0.9, 0.8, 0.0, 1.0,
		-1.0, 1.0, -1.0, 0.5, 0.9, 0.8, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.5, 0.9, 0.8, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.5, 0.9, 0.8, 0.0, 1.0,
		-1.0, 1.0, 1.0, 0.5, 0.9, 0.8, 1.0, 1.0,
		-1.0, 1.0, -1.0, 0.5, 0.9, 0.8, 1.0, 0.0,

		// Right
		1.0, -1.0, 1.0, 0.5, 0.9, 0.8, 1.0, 1.0,
		1.0, -1.0, -1.0, 0.5, 0.9, 0.8, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.5, 0.9, 0.8, 0.0, 0.0,
		1.0, -1.0, 1.0, 0.5, 0.9, 0.8, 1.0, 1.0,
		1.0, 1.0, -1.0, 0.5, 0.9, 0.8, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.5, 0.9, 0.8, 0.0, 1.0,
	}
	q.init()
	return q
}

func (s *mesh) init() {
	const sizeOfFloat = 4

	// Create buffers/arrays
	gl.GenVertexArrays(1, s.vao)
	gl.GenBuffers(1, s.vbo)
	//gl.GenBuffers(1, s.ebo);

	gl.BindVertexArray(*s.vao)

	// load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, *s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(s.Vertices)*sizeOfFloat, gl.Ptr(s.Vertices), gl.STATIC_DRAW)

	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.ebo)
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, len(s.Indices)*sizeOfFloat, gl.Ptr(s.Indices), gl.STATIC_DRAW)

	// vertex position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*sizeOfFloat, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// colour (to be replaced with normals)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*sizeOfFloat, gl.PtrOffset(3*sizeOfFloat))
	gl.EnableVertexAttribArray(1)
	// vertex Texture coords
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*sizeOfFloat, gl.PtrOffset(6*sizeOfFloat))
	gl.EnableVertexAttribArray(2)
	// ensure that no one else by
	gl.BindVertexArray(0)
}

func (s mesh) Draw(shader *Shader) {

	// textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.Texture)

	// draw mesh
	gl.BindVertexArray(*s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.Vertices)))

	// set back defaults, good practice stuff
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

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
