package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func ReferenceLevel(graph SceneNode) {
	{
		model := LoadModel("models/plane", MaterialMesh)
		t := mgl32.Translate3D(0, 0, 0)
		t = t.Mul4(mgl32.Scale3D(10, 1, 10))
		graph.Add(model, t)

		t = mgl32.Translate3D(0, 10, -10)
		t = t.Mul4(mgl32.Scale3D(10, 10, 1))
		t = t.Mul4(mgl32.HomogRotate3D(math.Pi/2, mgl32.Vec3{1, 0, 0}))
		graph.Add(model, t)
	}
	{
		model := LoadModel("models/cube", MaterialMesh)
		t := mgl32.Translate3D(-8, 1, 0)
		graph.Add(model, t)
	}
	{
		model := LoadModel("models/ico", MaterialMesh)
		t := mgl32.Translate3D(-4, 1, 0)
		graph.Add(model, t)
	}
	{
		model := LoadModel("models/sphere", MaterialMesh)
		t := mgl32.Translate3D(0, 1, 0)
		graph.Add(model, t)
	}
	{
		model := LoadModel("models/monkey", MaterialMesh)
		t := mgl32.Translate3D(4, 1, 0)
		graph.Add(model, t)
	}
	{
		model := LoadModel("models/rock1", TextureMesh)
		t := mgl32.Translate3D(8, 1, 0)
		graph.Add(model, t)
	}
}
