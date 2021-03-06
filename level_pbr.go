package main

import "github.com/go-gl/mathgl/mgl32"

func PBRLevel(graph SceneNode) {
	{
		meshes := LoadModel("models/cube", MaterialMesh)
		//meshes := LoadModel("models/cube", TexturedMesh)

		albTexture := GetTexture(Albedo, "rock_floor/Base_Color.png", true)
		metallicTexture := GetTexture(Metallic, "rock_floor/Metallic.png", false)
		normalTexture := GetTexture(Normal, "rock_floor/Normal.png", false)
		roughnessTexture := GetTexture(Roughness, "rock_floor/Roughness.png", false)

		for _, mesh := range meshes {
			mesh.Textures = append(mesh.Textures, albTexture)
			mesh.Textures = append(mesh.Textures, metallicTexture)
			mesh.Textures = append(mesh.Textures, roughnessTexture)
			mesh.Textures = append(mesh.Textures, normalTexture)
			mesh.Albedo = [3]float32{1, 1, 1}
			mesh.Metallic = 0.01
			mesh.Roughness = 0.7
		}
		for x := 0; x < 30; x++ {
			for z := 0; z < 30; z++ {
				t := mgl32.Translate3D(float32(x)*2-30, -0.5, float32(z)*2-30)
				t = t.Mul4(mgl32.Scale3D(1, 0.5, 1))
				//t = t.Mul4(mgl32.HomogRotate3D(float32(x)*3.14/2, mgl32.Vec3{0, 1, 0}))
				graph.Add(meshes, t)
			}
		}
	}
	{
		meshes := LoadModel("models/ico", MaterialMesh)
		for _, mesh := range meshes {
			mesh.Albedo = [3]float32{1, 0, 0}
			mesh.Metallic = 0.01
			mesh.Roughness = 0.8
		}
		t := mgl32.Translate3D(25, 5, -1)
		graph.Add(meshes, t)
	}

	{

		albTexture := GetTexture(Albedo, "streaked-marble/streaked-marble-albedo2.png", true)
		metallicTexture := GetTexture(Metallic, "streaked-marble/streaked-marble-metalness.png", false)
		normalTexture := GetTexture(Normal, "streaked-marble/streaked-marble-normal.png", false)
		roughnessTexture := GetTexture(Roughness, "streaked-marble/streaked-marble-roughness1.png", false)

		meshes := LoadModel("models/winged_victory", TexturedMesh)
		for _, mesh := range meshes {
			mesh.Textures = append(mesh.Textures, albTexture)
			mesh.Textures = append(mesh.Textures, metallicTexture)
			mesh.Textures = append(mesh.Textures, roughnessTexture)
			mesh.Textures = append(mesh.Textures, normalTexture)
		}
		t := mgl32.Translate3D(-5, 0, -4)
		t = t.Mul4(mgl32.HomogRotate3D(-3.14/4, mgl32.Vec3{0, 1, 0}))
		graph.Add(meshes, t)
	}

	plasticMetTex := GetTexture(Metallic, "scuffed-plastic/scuffed-plastic-metal.png", false)
	plasticNormTex := GetTexture(Normal, "scuffed-plastic/scuffed-plastic-normal.png", false)
	plasticRoughTex := GetTexture(Roughness, "scuffed-plastic/scuffed-plastic-rough.png", false)
	{
		meshes := LoadModel("models/sphere", TexturedMesh)
		albTexture := GetTexture(Albedo, "scuffed-plastic/scuffed-plastic4-alb.png", true)
		for _, mesh := range meshes {
			mesh.Textures = append(mesh.Textures, albTexture)
			mesh.Textures = append(mesh.Textures, plasticMetTex)
			mesh.Textures = append(mesh.Textures, plasticRoughTex)
			mesh.Textures = append(mesh.Textures, plasticNormTex)
		}
		t := mgl32.Translate3D(-8, 1, 12)
		t = t.Mul4(mgl32.HomogRotate3D(float32(1)*0.314*4, mgl32.Vec3{0, 1, 0}))
		graph.Add(meshes, t)
	}

	{
		meshes := LoadModel("models/test", TexturedMesh)
		albTexture := GetTexture(Albedo, "scuffed-plastic/scuffed-plastic5-alb.png", true)
		for _, mesh := range meshes {
			mesh.Textures = append(mesh.Textures, albTexture)
			mesh.Textures = append(mesh.Textures, plasticMetTex)
			mesh.Textures = append(mesh.Textures, plasticRoughTex)
			mesh.Textures = append(mesh.Textures, plasticNormTex)
		}
		t := mgl32.Translate3D(1, 4, -10)
		graph.Add(meshes, t)
	}

	// green sphere
	{
		meshes := LoadModel("models/sphere", MaterialMesh)
		for _, mesh := range meshes {
			mesh.Albedo = [3]float32{0, 1, 0}
			mesh.Metallic = 0.01
			mesh.Roughness = 0.1
		}
		t := mgl32.Translate3D(0, 1, 16)
		graph.Add(meshes, t)
	}

	// green cube
	{
		meshes := LoadModel("models/beveled_cube", MaterialMesh)
		for _, mesh := range meshes {
			mesh.Albedo = [3]float32{0, 1, 0}
			mesh.Metallic = 0.01
			mesh.Roughness = 0.1
		}

		t := mgl32.Translate3D(0, 0.0, 10)
		graph.Add(meshes, t)
	}

	{
		meshes := LoadModel("models/sphere", TexturedMesh)
		albTexture := GetTexture(Albedo, "scuffed-plastic/scuffed-plastic6-alb.png", true)
		for _, mesh := range meshes {
			mesh.Textures = append(mesh.Textures, albTexture)
			mesh.Textures = append(mesh.Textures, plasticMetTex)
			mesh.Textures = append(mesh.Textures, plasticRoughTex)
			mesh.Textures = append(mesh.Textures, plasticNormTex)
		}
		t := mgl32.Translate3D(-8, 1, 16)
		t = t.Mul4(mgl32.HomogRotate3D(float32(2)*0.314*4, mgl32.Vec3{0, 1, 0}))
		graph.Add(meshes, t)
	}

	{
		meshes := LoadModel("models/sphere", TexturedMesh)
		albTexture := GetTexture(Albedo, "scuffed-plastic/scuffed-plastic-alb.png", true)
		for _, mesh := range meshes {
			mesh.Textures = append(mesh.Textures, albTexture)
			mesh.Textures = append(mesh.Textures, plasticMetTex)
			mesh.Textures = append(mesh.Textures, plasticRoughTex)
			mesh.Textures = append(mesh.Textures, plasticNormTex)
		}
		t := mgl32.Translate3D(-8, 1, 20)
		t = t.Mul4(mgl32.HomogRotate3D(float32(3)*0.314*4, mgl32.Vec3{0, 1, 0}))
		graph.Add(meshes, t)
	}

	{

		meshes := LoadModel("models/sphere_bot", TexturedMesh)
		iAlb := GetTexture(Albedo, "sphere_bot/Robot_innerbody_Albedo.png", true)
		iMet := GetTexture(Metallic, "sphere_bot/Robot_innerbody_Metallic.png", false)
		iNorm := GetTexture(Normal, "sphere_bot/Robot_innerbody_Normal.png", false)
		iRough := GetTexture(Roughness, "sphere_bot/Robot_innerbody_Roughness.png", false)
		meshes[1].Textures = append(meshes[1].Textures, iAlb)
		meshes[1].Textures = append(meshes[1].Textures, iMet)
		meshes[1].Textures = append(meshes[1].Textures, iNorm)
		meshes[1].Textures = append(meshes[1].Textures, iRough)

		oAlb := GetTexture(Albedo, "sphere_bot/Robot_outerbody_Albedo.png", true)
		oMet := GetTexture(Metallic, "sphere_bot/Robot_outerbody_Metallic.png", false)
		oNorm := GetTexture(Normal, "sphere_bot/Robot_outerbody_Normal.png", false)
		oRough := GetTexture(Roughness, "sphere_bot/Robot_outerbody_Roughness.png", false)
		meshes[0].Textures = append(meshes[0].Textures, oAlb)
		meshes[0].Textures = append(meshes[0].Textures, oMet)
		meshes[0].Textures = append(meshes[0].Textures, oNorm)
		meshes[0].Textures = append(meshes[0].Textures, oRough)

		for i := 1; i < 3; i++ {
			t := mgl32.Translate3D(-24, -0.1, float32(i)*7)
			t = t.Mul4(mgl32.HomogRotate3D(float32(i)*0.314*4, mgl32.Vec3{0, 1, 0}))
			graph.Add(meshes, t)
		}
	}

}
