package main

import "github.com/go-gl/gl/v4.1-core/gl"

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

	colours := []float32{0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0}

	gl.GenBuffers(1, s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(s.vertices)*4, gl.Ptr(s.vertices), gl.STATIC_DRAW)

	var color_vbo uint32 = 1
	gl.GenBuffers(1, &color_vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, color_vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(colours)*4, gl.Ptr(colours), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, s.vao)
	gl.BindVertexArray(*s.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, *s.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, color_vbo)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(1)
}

func (s square) Draw() {
	gl.BindVertexArray(*s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.vertices)))
}
