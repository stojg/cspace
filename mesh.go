package main

import (
	"fmt"
	"log"

	"math"

	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func newCrateMesh() *Mesh {
	vertices := getVertices(cubeData)
	var textures []*Texture
	var indices []uint32

	diffuseTexture, err := newTexture(Diffuse, "textures/crate0/crate0_diffuse.png")
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, diffuseTexture)

	specularTexture, err := newTexture(Specular, "textures/specular.png")
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, specularTexture)

	normalTexture, err := newTexture(Normal, "textures/crate0/crate0_normal.png")
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, normalTexture)

	return NewMesh(vertices, indices, textures)
}

func newLightMesh() *Mesh {
	vertices := getVertices(cubeData)
	var textures []*Texture
	var indices []uint32
	return NewMesh(vertices, indices, textures)
}

type Vertex struct {
	Position  [3]float32
	Normal    [3]float32
	TexCoords [2]float32
	Tangent   [3]float32
	BiTangent [3]float32
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
	Vertices []Vertex
	Indices  []uint32
	Textures []*Texture

	vbo, vao, ebo uint32
}

func (s *Mesh) Draw(shader *Shader) {
	diffuseNr := 0
	specularNr := 0
	normalNr := 0
	for i, texture := range s.Textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		var number int
		switch texture.textureType {
		case Specular:
			number = specularNr
			specularNr++
		case Diffuse:
			number = diffuseNr
			diffuseNr++
		case Normal:
			number = normalNr
			normalNr++
		default:
			panic("unknown texture type ")
		}

		uniformName := fmt.Sprintf("mat.%s%d", texture.textureType, number)
		gl.Uniform1i(uniformLocation(shader, uniformName), int32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture.ID)
	}

	location := gl.GetUniformLocation(shader.Program, gl.Str("mat.shininess\x00"))
	if location > 0 {
		gl.Uniform1f(location, 32.0)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.Vertices)))

	// reset textures to they don't leak into some other poor mesh
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

	size := int32(unsafe.Sizeof(Vertex{}))
	gl.BufferData(gl.ARRAY_BUFFER, len(s.Vertices)*int(size), gl.Ptr(s.Vertices), gl.STATIC_DRAW)

	if len(s.Indices) > 0 {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(s.Indices)*3, gl.Ptr(s.Indices), gl.STATIC_DRAW)
	}

	// vertex position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, size, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// normals
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, size, gl.PtrOffset(3*sizeOfFloat))
	gl.EnableVertexAttribArray(1)

	// texture coordinates
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, size, gl.PtrOffset(6*sizeOfFloat))
	gl.EnableVertexAttribArray(2)

	// tangents
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, size, gl.PtrOffset(8*sizeOfFloat))
	gl.EnableVertexAttribArray(3)

	// bi-tangents
	gl.VertexAttribPointer(4, 3, gl.FLOAT, false, size, gl.PtrOffset(11*sizeOfFloat))
	gl.EnableVertexAttribArray(4)

	// reset, so no other mesh accidentally changes this vao
	gl.BindVertexArray(0)
}

func edge(a, b Vertex) [3]float32 {
	return [3]float32{
		a.Position[0] - b.Position[0],
		a.Position[1] - b.Position[1],
		a.Position[2] - b.Position[2],
	}
}

func deltaUV(a, b Vertex) [2]float32 {
	return [2]float32{
		a.TexCoords[0] - b.TexCoords[0],
		a.TexCoords[1] - b.TexCoords[1],
	}
}

func normalise(vec [3]float32) [3]float32 {
	l := 1.0 / float32(math.Sqrt(float64(vec[0]*vec[0]+vec[1]*vec[1]+vec[2]*vec[2])))
	return [3]float32{vec[0] * l, vec[1] * l, vec[2] * l}
}

func getVertices(meshdata []float32) []Vertex {
	const stride = 8

	if len(meshdata)%stride != 0 {
		panic("the mesh data is not a multiple of 8, want [3]Pos, [3]Normals, [2]TexCoords")
	}
	var vertices []Vertex

	for i := 0; i < len(meshdata); i += stride {
		var vertex Vertex
		copy(vertex.Position[:], meshdata[i:i+3])
		copy(vertex.Normal[:], meshdata[i+3:i+6])
		copy(vertex.TexCoords[:], meshdata[i+6:i+8])
		vertices = append(vertices, vertex)
	}

	// calculate tangents and bi-tangents
	for i := 0; i < len(vertices); i += 3 {
		deltaPos1 := edge(vertices[i+1], vertices[i])
		deltaPos2 := edge(vertices[i+2], vertices[i])

		deltaUV1 := deltaUV(vertices[i+1], vertices[i])
		deltaUV2 := deltaUV(vertices[i+2], vertices[i])

		f := 1.0 / (deltaUV1[0]*deltaUV2[1] - deltaUV2[0]*deltaUV1[1])

		var tangent [3]float32
		tangent[0] = f * (deltaUV2[1]*deltaPos1[0] - deltaUV1[1]*deltaPos2[0])
		tangent[1] = f * (deltaUV2[1]*deltaPos1[1] - deltaUV1[1]*deltaPos2[1])
		tangent[2] = f * (deltaUV2[1]*deltaPos1[2] - deltaUV1[1]*deltaPos2[2])
		tangent = normalise(tangent)

		var biTangent [3]float32
		biTangent[0] = f * (-deltaUV2[0]*deltaPos1[0] - deltaUV1[1]*deltaPos2[0])
		biTangent[1] = f * (-deltaUV2[0]*deltaPos1[1] - deltaUV1[1]*deltaPos2[1])
		biTangent[2] = f * (-deltaUV2[0]*deltaPos1[2] - deltaUV1[1]*deltaPos2[2])
		biTangent = normalise(biTangent)

		copy(vertices[i].Tangent[:], tangent[:])
		copy(vertices[i].BiTangent[:], biTangent[:])

		copy(vertices[i+1].Tangent[:], tangent[:])
		copy(vertices[i+1].BiTangent[:], biTangent[:])

		copy(vertices[i+2].Tangent[:], tangent[:])
		copy(vertices[i+2].BiTangent[:], biTangent[:])
	}

	return vertices
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
