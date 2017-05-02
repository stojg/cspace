package main

import (
	"fmt"
	"log"

	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func newCubeMesh() *Mesh {

	const perRowSize = 8

	if len(cubeData)%perRowSize != 0 {
		panic("the mesh data is not a multiple of 8, want [3]Pos, [3]Normals, [2]TexCoords")
	}

	var vertices []Vertex
	var indices []uint32
	for i := 0; i < len(cubeData); i += perRowSize {
		var vertex Vertex
		copy(vertex.Position[:], cubeData[i:i+3])
		copy(vertex.Normal[:], cubeData[i+3:i+6])
		copy(vertex.TexCoords[:], cubeData[i+6:i+8])
		vertices = append(vertices, vertex)
	}

	var textures []*Texture

	diffuseTexture, err := newTexture("diffuse", "textures/crate0/crate0_diffuse.png")
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, diffuseTexture)
	specularTexture, err := newTexture("specular", "textures/specular.png")
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, specularTexture)

	return NewMesh(vertices, indices, textures)
}

type Vertex struct {
	Position  [3]float32
	Normal    [3]float32
	TexCoords [2]float32
}

type Texture struct {
	ID   uint32
	Name string // type of texture, like diffuse, specular or bump
}

func NewMesh(vertices []Vertex, Indices []uint32, textures []*Texture) *Mesh {
	q := &Mesh{
		Vertices: vertices,
		Indices:  Indices,
		Textures: textures,
	}
	q.init()
	return q
}

type Mesh struct {
	Vertices      []Vertex
	Indices       []uint32
	Textures      []*Texture
	vbo, vao, ebo uint32
}

func (s *Mesh) Draw(shader *Shader) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.Textures[0].ID)
	gl.Uniform1i(uniformLocation(shader, "materialDiffuse"), 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, s.Textures[0].ID)
	gl.Uniform1i(uniformLocation(shader, "materialSpecular"), 1)
	gl.Uniform1f(uniformLocation(shader, "materialShininess"), 32.0)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.Vertices)))

	for i := range s.Textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
}

func (s *Mesh) init() {
	const sizeOfFloat = 4

	// Create buffers/arrays
	gl.GenVertexArrays(1, &s.vao)
	gl.GenBuffers(1, &s.vbo)

	gl.BindVertexArray(s.vao)

	// load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	fmt.Println(unsafe.Sizeof(&Vertex{}))
	// 32 is the byte size of the Vertex struct
	gl.BufferData(gl.ARRAY_BUFFER, len(s.Vertices)*32, gl.Ptr(s.Vertices), gl.STATIC_DRAW)

	if len(s.Indices) > 0 {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(s.Indices)*3, gl.Ptr(s.Indices), gl.STATIC_DRAW)
	}

	// vertex position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*sizeOfFloat, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// normals
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*sizeOfFloat, gl.PtrOffset(3*sizeOfFloat))
	gl.EnableVertexAttribArray(1)
	// texture coordinates
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*sizeOfFloat, gl.PtrOffset(6*sizeOfFloat))
	gl.EnableVertexAttribArray(2)

	// reset, so no other mesh accidentally changes this vao
	gl.BindVertexArray(0)
}

var cubeData = []float32{
	// Positions      Normals         Texture Coords
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
	-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
}
