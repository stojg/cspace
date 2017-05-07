package main

type PointLight struct {
	Position         [3]float32
	Color            [3]float32
	AmbientIntensity float32
	DiffuseIntensity float32
	Constant         float32
	Linear           float32
	Exp              float32
}
