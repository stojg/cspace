package main

import (
	"math"

	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/cspace/lib/obj"
)

//func NewGrassMesh() *Mesh {
//	var data = []float32{
//		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
//		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
//		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0, // top-right
//		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
//		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
//		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0, // bottom-left
//	}
//	vertices := getVertices(data)
//	var textures []*Texture
//	//diffuseTexture, err := newTexture(Diffuse, "textures/grass/grass01.jpg", true)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//textures = append(textures, diffuseTexture)
//	//
//	//specularTexture, err := newTexture(Specular, "textures/grass/grass01_s.jpg", true)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//textures = append(textures, specularTexture)
//	//
//	//normalTexture, err := newTexture(Normal, "textures/grass/grass01_n.jpg", true)
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//textures = append(textures, normalTexture)
//	mat := obj.NewMaterial()
//	mat.Ambient = [3]float32{1, 1, 1}
//	mat.Diffuse = [3]float32{0.242558, 0.079845, 0.040770}
//	mat.Name = "grass"
//	mat.SpecularExp = 0
//	mat.Specular = [3]float32{0, 0.3, 0}
//	return NewMesh("grass", vertices, textures, mat)
//}

type Vertex struct {
	Position  [3]float32
	Normal    [3]float32
	TexCoords [2]float32
	Tangent   [3]float32
}

func NewMesh(name string, vertices []Vertex, textures []*Texture, mat *obj.Material, shaderType ShaderType) *Mesh {
	q := &Mesh{
		Name:     name,
		Vertices: vertices,
		Textures: textures,
		Material: mat,
	}
	q.MeshType = shaderType
	q.init()
	return q
}

type Mesh struct {
	Name     string
	Vertices []Vertex
	Indices  []uint32
	Textures []*Texture
	Material *obj.Material
	MeshType ShaderType

	vbo, vao, ebo uint32
}

func (s *Mesh) Render() {

	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.Vertices)))
	// reset
	gl.BindVertexArray(0)

	// reset textures
	//for i := range s.Textures {
	//	gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
	//	gl.BindTexture(gl.TEXTURE_2D, 0)
	//}
}

func (s *Mesh) setTextures(tShader TextureShader) {
	diffuseNr := 0
	specularNr := 0
	normalNr := 0
	for i, texture := range s.Textures {
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
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.Uniform1i(tShader.TextureUniform(texture.textureType, number), int32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture.ID)
	}
}

func (s *Mesh) setMaterial(mShader MaterialShader) {
	gl.Uniform3f(mShader.DiffuseUniform(), s.Material.Diffuse[0], s.Material.Diffuse[1], s.Material.Diffuse[2])
	gl.Uniform1f(mShader.SpecularExpUniform(), 0.001)
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
