package main

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

func MaterialLevel(graph SceneNode) {

	grass := LoadModel("models/grass1", MaterialMesh)
	for x := 0; x < 20; x++ {
		for z := 0; z < 20; z++ {
			t := mgl32.Translate3D(float32(x)*3-30, 0, float32(z)*3-30)
			t = t.Mul4(mgl32.Scale3D(3, 1, 3))
			graph.Add(grass, t)
		}
	}

	{
		tree := LoadModel("models/tree1", MaterialMesh)
		for i := 0; i < 5; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			s := rand.Float32()/4 + 1
			t = t.Mul4(mgl32.Scale3D(s, s, s))
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}
	{
		rock := LoadModel("models/stone1", MaterialMesh)
		for i := 0; i < 30; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(rock, t)
		}
	}
	{
		rock := LoadModel("models/stone2", MaterialMesh)
		for i := 0; i < 30; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(rock, t)
		}
	}
	{
		tree := LoadModel("models/bush1", MaterialMesh)
		for i := 0; i < 15; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}
	{
		tree := LoadModel("models/grass2", MaterialMesh)
		for i := 0; i < 1000; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}
	{
		tree := LoadModel("models/grass3", MaterialMesh)
		for i := 0; i < 50; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}
	{
		tree := LoadModel("models/tulip1", MaterialMesh)
		for i := 0; i < 10; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}
	{
		tree := LoadModel("models/bush2", MaterialMesh)
		for i := 0; i < 15; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}

	{
		tree := LoadModel("models/tree2", MaterialMesh)
		for i := 0; i < 2; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}

	{
		tree := LoadModel("models/tree3", MaterialMesh)
		for i := 0; i < 8; i++ {
			t := mgl32.Translate3D(rand.Float32()*60-30, 0.1, rand.Float32()*60-30)
			t = t.Mul4(mgl32.HomogRotate3D(rand.Float32()*360, mgl32.Vec3{0, 1, 0}.Normalize()))
			graph.Add(tree, t)
		}
	}

}
