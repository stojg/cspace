package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func PBRLevel(graph SceneNode) {

	concreteAlbedo, err := newTexture(Diffuse, "textures/patchy_cement1/base.png", false)
	if err != nil {
		panic(err)
	}
	concreteNormal, err := newTexture(Normal, "textures/patchy_cement1/normal.png", false)
	if err != nil {
		panic(err)
	}
	concreteRough, err := newTexture(Specular, "textures/patchy_cement1/roughness.png", false)
	if err != nil {
		panic(err)
	}

	bambooAlbedo, err := newTexture(Diffuse, "textures/bamboo/albedo.png", false)
	if err != nil {
		panic(err)
	}
	bambooNormal, err := newTexture(Normal, "textures/bamboo/normal.png", false)
	if err != nil {
		panic(err)
	}
	bambooRough, err := newTexture(Specular, "textures/bamboo/roughness.png", false)
	if err != nil {
		panic(err)
	}

	plasticAlbedo1, err := newTexture(Diffuse, "textures/scuffed-plastic/scuffed-plastic-alb.png", false)
	if err != nil {
		panic(err)
	}
	plasticAlbedo2, err := newTexture(Diffuse, "textures/scuffed-plastic/scuffed-plastic6-alb.png", false)
	if err != nil {
		panic(err)
	}
	plasticNormal, err := newTexture(Normal, "textures/scuffed-plastic/scuffed-plastic-normal.png", false)
	if err != nil {
		panic(err)
	}
	plasticRough, err := newTexture(Specular, "textures/scuffed-plastic/scuffed-plastic-rough.png", false)
	if err != nil {
		panic(err)
	}

	{
		model := LoadModel("models/plane", TextureMesh)
		for _, m := range model {
			m.Textures = append(m.Textures, concreteAlbedo)
			m.Textures = append(m.Textures, concreteNormal)
			m.Textures = append(m.Textures, concreteRough)
		}
		t := mgl32.Translate3D(0, 0, 0)
		t = t.Mul4(mgl32.Scale3D(10, 1, 10))
		graph.Add(model, t)

		t = mgl32.Translate3D(0, 10, -10)
		t = t.Mul4(mgl32.Scale3D(10, 10, 1))
		t = t.Mul4(mgl32.HomogRotate3D(math.Pi/2, mgl32.Vec3{1, 0, 0}))
		graph.Add(model, t)
	}
	//
	{
		model := LoadModel("models/cube", TextureMesh)
		for _, m := range model {
			m.Textures = append(m.Textures, bambooAlbedo)
			m.Textures = append(m.Textures, bambooNormal)
			m.Textures = append(m.Textures, bambooRough)
		}
		t := mgl32.Translate3D(-8, 1, 0)
		graph.Add(model, t)

	}
	{
		model := LoadModel("models/ico", TextureMesh)
		for _, m := range model {
			m.Textures = append(m.Textures, plasticAlbedo2)
			m.Textures = append(m.Textures, plasticNormal)
			m.Textures = append(m.Textures, plasticRough)
		}
		t := mgl32.Translate3D(-4, 1, 0)
		graph.Add(model, t)
	}
	{

		model := LoadModel("models/sphere", TextureMesh)
		for _, m := range model {
			m.Textures = append(m.Textures, plasticAlbedo1)
			m.Textures = append(m.Textures, plasticNormal)
			m.Textures = append(m.Textures, plasticRough)
		}
		t := mgl32.Translate3D(0, 1, 0)
		graph.Add(model, t)
	}
	{
		model := LoadModel("models/monkey", TextureMesh)
		for _, m := range model {
			m.Textures = append(m.Textures, bambooAlbedo)
			m.Textures = append(m.Textures, bambooNormal)
			m.Textures = append(m.Textures, bambooRough)
		}
		t := mgl32.Translate3D(4, 1, 0)
		graph.Add(model, t)
	}
	{
		model := LoadModel("models/rock1", MaterialMesh)
		t := mgl32.Translate3D(8, 1, 0)
		graph.Add(model, t)
	}

}
