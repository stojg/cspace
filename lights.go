package main

import "math"

type PointLight struct {
	Position         [3]float32
	Color            [3]float32
	AmbientIntensity float32
	DiffuseIntensity float32
	Constant         float32
	Linear           float32
	Exp              float32
}

func (l *PointLight) Radius() float32 {

	maxChannel := float32(math.Max(math.Max(float64(l.Color[0]), float64(l.Color[1])), float64(l.Color[2])))

	inner := l.Linear*l.Linear - 4*l.Exp*(l.Exp-256*maxChannel*l.DiffuseIntensity)

	ret := -l.Linear + float32(math.Sqrt(float64(inner)))/(2*l.Exp)

	if math.IsNaN(float64(ret)) {
		panic("CalcPointLightBSphere calculated a NaN")
	}
	return ret

}
