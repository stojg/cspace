package main

import (
	"path/filepath"

	"github.com/stojg/cspace/lib/obj"
)

type Transform struct {
}

type Model struct {
	Meshes     []Mesh
	Transforms []Transform
}

func LoadModel(directory string) *Mesh {

	filePath := filepath.Join(directory, "model.obj")

	data := obj.LoadObject(filePath)

	glLogf("--- Loaded %s ----\n", directory)
	glLogf("size %d bytes\n", len(data))

	vertices := getVertices(data)
	var textures []*Texture
	var indices []uint32

	diffuseTexture, err := newTexture(Diffuse, filepath.Join(directory, "d.png"), false)
	if err == nil {
		textures = append(textures, diffuseTexture)
	}

	specularTexture, err := newTexture(Specular, filepath.Join(directory, "s.png"), false)
	if err == nil {
		textures = append(textures, specularTexture)
	}

	normalTexture, err := newTexture(Normal, filepath.Join(directory, "n.png"), false)
	if err == nil {
		textures = append(textures, normalTexture)
	}

	glLogf("textures %d \n", len(textures))
	glLogln("------------------------")

	return NewMesh(vertices, indices, textures)
}
