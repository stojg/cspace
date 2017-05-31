package main

import (
	"math"

	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/cspace/lib/obj"
)

type Vertex struct {
	Position  [3]float32
	Normal    [3]float32
	TexCoords [2]float32
	Tangent   [3]float32
}

func NewMesh(name string, vertices []Vertex, textures []*Texture, mat *obj.Material, shaderType ShaderType) *Mesh {
	q := &Mesh{
		Name:        name,
		Vertices:    vertices,
		NumVertices: int32(len(vertices)),
		Textures:    textures,
		Material:    mat,
	}
	q.MeshType = shaderType
	q.init()
	return q
}

type Mesh struct {
	Name        string
	Vertices    []Vertex
	NumVertices int32
	Indices     []uint32
	Textures    []*Texture
	Material    *obj.Material
	MeshType    ShaderType
	// PRB material
	Albedo    [3]float32
	Metallic  float32
	Roughness float32

	vbo, vao uint32
}

func (s *Mesh) Render() {
	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, s.NumVertices)
}

func (s *Mesh) setTextures(tShader *GbufferTShader) {
	for i := range s.Textures {
		GLBindTexture(i, tShader.TextureUniform(s.Textures[i].textureType), s.Textures[i].ID)
	}
}

func (s *Mesh) setMaterial(mShader *GbufferMShader) {
	gl.Uniform3f(mShader.LocAlbedo, s.Albedo[0], s.Albedo[1], s.Albedo[2])
	gl.Uniform1f(mShader.LocMetallic, s.Metallic)
	gl.Uniform1f(mShader.LocRoughness, s.Roughness)
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

	// reset, so no other graph accidentally changes this vao
	gl.BindVertexArray(0)
}

func edge(a, b Vertex) [3]float32 {
	return [3]float32{
		a.Position[0] - b.Position[0],
		a.Position[1] - b.Position[1],
		a.Position[2] - b.Position[2],
	}
}

func normalise(vec [3]float32) [3]float32 {
	l := 1.0 / float32(math.Sqrt(float64(vec[0]*vec[0]+vec[1]*vec[1]+vec[2]*vec[2])))
	return [3]float32{vec[0] * l, vec[1] * l, vec[2] * l}
}

func getVertices(meshdata []float32) []Vertex {
	const stride = 8

	if len(meshdata)%stride != 0 {
		panic("the graph data is not a multiple of 8, want [3]Pos, [3]Normals, [2]TexCoords")
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
	return vertices
}
