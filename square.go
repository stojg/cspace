package main

import "github.com/go-gl/gl/v4.1-core/gl"

type pos [3]float32

func newSquare(index uint32, a, b pos) *square {
	s := &square{
		vbo: &index,
		vao: &index,
		vertices: []float32{
			a[0], a[1], 0.0,
			a[0], b[1], 0.0,
			b[0], a[1], 0.0,

			b[0], a[1], 0.0,
			b[0], b[1], 0.0,
			a[0], b[1], 0.0,
		},
	}
	s.setVertexBufferObject()
	s.setVertexArrayObject()
	return s
}

type square struct {
	vbo, vao *uint32
	vertices []float32
}

// vertex buffer object
// Configure the vertex data
// Now an unusual step. Most meshes will use a collection of one or more vertex buffer objects to hold vertex
// points, texture-coordinates, vertex normals, etc. In older GL implementations we would have to bind each one,
// and define their memory layout, every time that we draw the mesh. To simplify that, we have new thing called
// the vertex array object (VAO), which remembers all of the vertex buffers that you want to use, and the memory
// layout of each one. We set up the vertex array object once per mesh. When we want to draw, all we do then is
// bind the VAO and draw.
func (s *square) setVertexBufferObject() {
	gl.GenBuffers(1, s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(s.vertices)*4, gl.Ptr(s.vertices), gl.STATIC_DRAW)
}

func (s *square) setVertexArrayObject() {
	// Here we tell GL to generate a new VAO for us. It sets an unsigned integer to identify it with later.
	gl.GenVertexArrays(1, s.vao)
	// We bind it, to bring it in to focus in the state machine.
	gl.BindVertexArray(*s.vao)
	// This lets us enable the first attribute; 0. We are only using a single vertex buffer, so we know that it will
	// be attribute location 0
	gl.EnableVertexAttribArray(0)
	// The glVertexAttribPointer function defines the layout of our first vertex buffer; "0" means define the layout
	// for attribute number 0. "3" means that the variables are vec3 made from every 3 floats (GL_FLOAT) in the buffer.
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
}

func (s square) Count() int32 {
	return int32(len(s.vertices))
}

func (s square) Draw() {
	gl.BindVertexArray(*s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, s.Count())
}
