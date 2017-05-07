package main

import "math"

type DirectionalLight struct {
	Direction [3]float32
	Color     [3]float32
}

type PointLight struct {
	Position [3]float32
	Color    [3]float32
	Constant float32
	Linear   float32
	Exp      float32
	radius   float32
	rand     float32
}

func (l *PointLight) Radius() float32 {
	if l.radius != 0 {
		return l.radius
	}
	maxChannel := float32(math.Max(math.Max(float64(l.Color[0]), float64(l.Color[1])), float64(l.Color[2])))
	inner := l.Linear*l.Linear - 4*l.Exp*(l.Exp-(256/5)*maxChannel)
	ret := -l.Linear + float32(math.Sqrt(float64(inner)))/(2*l.Exp)
	if math.IsNaN(float64(ret)) {
		panic("CalcPointLightBSphere calculated a NaN")
	}
	l.radius = ret
	return ret

}

type LightAttenuation struct {
	Constant float32
	Linear   float32
	Exp      float32
}

var ligthAtt = map[int]LightAttenuation{
	1:    {1, 0.9, 1.8},
	7:    {1, 0.7, 1.8},
	13:   {1.0, 0.35, 0.44},
	20:   {1.0, 0.22, 0.20},
	32:   {1.0, 0.14, 0.07},
	50:   {1.0, 0.09, 0.032},
	65:   {1.0, 0.07, 0.017},
	100:  {1.0, 0.045, 0.0075},
	160:  {1.0, 0.027, 0.0028},
	200:  {1.0, 0.022, 0.0019},
	325:  {1.0, 0.014, 0.0007},
	600:  {1.0, 0.007, 0.0002},
	3250: {1.0, 0.0014, 0.000007},
}
