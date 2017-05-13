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

func LoadModel(directory string) []*Mesh {

	var result []*Mesh

	filePath := filepath.Join(directory, "model.obj")
	objects := obj.LoadObject(filePath)

	for _, object := range objects {
		glLogf("--- Loaded %s ----\n", object.Name)
		glLogf("size %d bytes\n", len(object.Data))

		vertices := getVertices(object.Data)
		var textures []*Texture

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

		result = append(result, NewMesh(object.Name, vertices, textures, object.Mtr))
	}
	return result

}
