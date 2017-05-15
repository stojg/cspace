package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func ReferenceLevel(graph SceneNode) {
	{
		model := LoadModel("models/plane")
		t := mgl32.Translate3D(0, 0, 0)
		t = t.Mul4(mgl32.Scale3D(10, 1, 10))
		graph.Add(model, MaterialMesh, t)

		t = mgl32.Translate3D(0, 10, -10)
		t = t.Mul4(mgl32.Scale3D(10, 10, 1))
		t = t.Mul4(mgl32.HomogRotate3D(math.Pi/2, mgl32.Vec3{1, 0, 0}))
		graph.Add(model, MaterialMesh, t)
	}
	{
		model := LoadModel("models/cube")
		t := mgl32.Translate3D(-8, 1, 0)
		graph.Add(model, MaterialMesh, t)
	}
	{
		model := LoadModel("models/ico")
		t := mgl32.Translate3D(-4, 1, 0)
		graph.Add(model, MaterialMesh, t)
	}
	{
		model := LoadModel("models/sphere")
		t := mgl32.Translate3D(0, 1, 0)
		graph.Add(model, MaterialMesh, t)
	}
	{
		model := LoadModel("models/monkey")
		t := mgl32.Translate3D(4, 1, 0)
		graph.Add(model, MaterialMesh, t)
	}
	{
		model := LoadModel("models/rock1")
		t := mgl32.Translate3D(8, 1, 0)
		graph.Add(model, TextureMesh, t)
	}

	{
		model := LoadModel("models/corvette")
		t := mgl32.Translate3D(12, 1, 0)
		t = t.Mul4(mgl32.Scale3D(0.001, 0.001, 0.001))
		graph.Add(model, TextureMesh, t)
	}
}
