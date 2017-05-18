package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func PBRLevel(graph SceneNode) {

	{
		meshes := LoadModel("models/plane", MaterialMesh)
		for _, mesh := range meshes {
			mesh.Albedo = [3]float32{1, 1, 1}
			mesh.Metallic = 0.0
			mesh.Roughness = 0.8
		}
		t := mgl32.Translate3D(0, 0, 0)
		t = t.Mul4(mgl32.Scale3D(30, 1, 10))

		graph.Add(meshes, t)

		t = mgl32.Translate3D(0, 10, -10)
		t = t.Mul4(mgl32.Scale3D(30, 10, 1))
		t = t.Mul4(mgl32.HomogRotate3D(math.Pi/2, mgl32.Vec3{1, 0, 0}))
		graph.Add(meshes, t)
	}

	for x := float32(0); x < 1.1; x += 0.1 {
		meshes := LoadModel("models/sphere", MaterialMesh)
		for _, mesh := range meshes {
			mesh.Albedo = [3]float32{1, 1, 1}
			mesh.Metallic = x
			mesh.Roughness = 0.3
		}
		t := mgl32.Translate3D(x*25, 2, 0)
		graph.Add(meshes, t)
	}

	for x := float32(0); x < 1.1; x += 0.1 {
		meshes := LoadModel("models/ico", MaterialMesh)
		for _, mesh := range meshes {
			mesh.Albedo = [3]float32{1, 0, 0}
			mesh.Metallic = 0.01
			mesh.Roughness = x
		}
		t := mgl32.Translate3D(x*25, 5, -2)
		graph.Add(meshes, t)
	}

	{
		meshes := LoadModel("models/Wood_Log_qdtdP_256_3d_ms-1", MaterialMesh)
		for _, mesh := range meshes {
			mesh.Albedo = [3]float32{0.486, 0.580, 0.392}
			mesh.Metallic = 0
			mesh.Roughness = 0.5
		}
		t := mgl32.Translate3D(-10, 0, 0)
		t = t.Mul4(mgl32.Scale3D(0.1, 0.1, 0.1))
		graph.Add(meshes, t)
	}
}
