package main

import "github.com/go-gl/mathgl/mgl32"

func TexturedLevel(graph SceneNode) {
	{
		corvette := LoadModel("models/corvette")
		t := mgl32.Translate3D(0, 0.0, 0)
		t = t.Mul4(mgl32.HomogRotate3D(3.14/2, mgl32.Vec3{0, 1, 0}.Normalize()))
		graph.Add(corvette, TextureMesh, t)
	}
}
