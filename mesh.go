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

	name := "rock"

	diffuseTexture, err := newTexture(Diffuse, fmt.Sprintf("textures/%s/diffuse.jpg", name))
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, diffuseTexture)

	specularTexture, err := newTexture(Specular, fmt.Sprintf("textures/%s/specular.jpg", name))
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, specularTexture)

	normalTexture, err := newTexture(Normal, fmt.Sprintf("textures/%s/normal.jpg", name))
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, normalTexture)

	depthTexture, err := newTexture(Depth, fmt.Sprintf("textures/%s/depth.jpg", name))
	if err != nil {
		log.Fatalln(err)
	}
	textures = append(textures, depthTexture)

	return NewMesh(vertices, indices, textures)
}

func newLightMesh() *Mesh {
	vertices := getVertices(cubeData)
	var textures []*Texture
	var indices []uint32
	return NewMesh(vertices, indices, textures)
}

func newPlaneMesh() *Mesh {
	var data = []float32{
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0, // top-right
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0, // bottom-left
	}
	vertices := getVertices(data)
	var textures []*Texture
	var indices []uint32
	return NewMesh(vertices, indices, textures)
}

type Vertex struct {
	Position  [3]float32
	Normal    [3]float32
	TexCoords [2]float32
	Tangent   [3]float32
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
	depthNr := 0
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
		case Depth:
			number = depthNr
			depthNr++
		default:
			panic("unknown texture type ")
		}

		uniformName := fmt.Sprintf("mat.%s%d", texture.textureType, number)
		gl.Uniform1i(uniformLocation(shader, uniformName), int32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture.ID)
	}

	location := gl.GetUniformLocation(shader.Program, gl.Str("mat.shininess\x00"))
	if location > 0 {
		gl.Uniform1f(location, 64.0)
	}

	gl.BindVertexArray(s.vao)
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
		v0 := vertices[i]
		v1 := vertices[i+1]
		v2 := vertices[i+2]

		edge1 := edge(v1, v0)
		edge2 := edge(v2, v0)

		deltaU1 := v1.TexCoords[0] - v0.TexCoords[0]
		deltaV1 := v1.TexCoords[1] - v0.TexCoords[1]
		deltaU2 := v2.TexCoords[0] - v0.TexCoords[0]
		deltaV2 := v2.TexCoords[1] - v0.TexCoords[1]

		f := 1.0 / (deltaU1*deltaV2 - deltaU2*deltaV1)

		var tangent [3]float32
		tangent[0] = f * (deltaV2*edge1[0] - deltaV1*edge2[0])
		tangent[1] = f * (deltaV2*edge1[1] - deltaV1*edge2[1])
		tangent[2] = f * (deltaV2*edge1[2] - deltaV1*edge2[2])

		tangent = normalise(tangent)

		copy(vertices[i].Tangent[:], tangent[:])
		copy(vertices[i+1].Tangent[:], tangent[:])
		copy(vertices[i+2].Tangent[:], tangent[:])
	}

	//fmt.Println(vertices)
	//os.Exit(1)
	return vertices
}

var cubeData = []float32{
	// Positions      Normals         Texture Coords
	// Back face
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0, // Bottom-left
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0, // top-right
	0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0, // bottom-right
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0, // top-right
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0, // bottom-left
	-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0, // top-left
	// Front face
	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
	0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
	-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
	// Left face
	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0, // top-right
	-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0, // top-left
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0, // bottom-left
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0, // bottom-left
	-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0, // bottom-right
	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0, // top-right
	// Right face
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // top-left
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0, // bottom-right
	0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0, // top-right
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0, // bottom-right
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // top-left
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0, // bottom-left
	// Bottom face
	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0, // top-right
	0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0, // top-left
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0, // bottom-left
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0, // bottom-left
	-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0, // bottom-right
	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0, // top-right
	// Top face
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0, // top-right
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0, // bottom-left
}
