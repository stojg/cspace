package main

import "github.com/go-gl/mathgl/mgl32"

func PBRLevel(graph SceneNode) {

	//{
	//	model := LoadModel("models/plane", MaterialMesh)
	//	t := mgl32.Translate3D(0, 0, 0)
	//	t = t.Mul4(mgl32.Scale3D(10, 1, 10))
	//	graph.Add(model, t)
	//
	//	t = mgl32.Translate3D(0, 10, -10)
	//	t = t.Mul4(mgl32.Scale3D(10, 10, 1))
	//	t = t.Mul4(mgl32.HomogRotate3D(math.Pi/2, mgl32.Vec3{1, 0, 0}))
	//	graph.Add(model, t)
	//}

	for y := float32(0); y < 1.1; y += 0.1 {
		for x := float32(0); x < 1.1; x += 0.1 {
			meshes := LoadModel("models/sphere", MaterialMesh)
			for _, mesh := range meshes {
				mesh.Albedo = [3]float32{0.9, 0.5, 0.0}
				mesh.Metallic = y
				mesh.Roughness = x
			}
			t := mgl32.Translate3D(x*25, y*25, 0)
			graph.Add(meshes, t)
		}

	}
}
