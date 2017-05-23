package main

import "github.com/go-gl/gl/v4.1-core/gl"

func NewPointLightShader() *PointLightShader {
	c := NewDefaultShader("lighting", "lighting_point_pbr")
	s := &PointLightShader{
		DefaultShader: c,
		locNormal:     uniformLocation(c, "gNormal"),

		LocGAlbedo:      uniformLocation(c, "gAlbedoSpec"),
		uniformDepthLoc: uniformLocation(c, "gDepth"),

		uniformLightPosLoc:       uniformLocation(c, "pointLight.Position"),
		uniformLightColorLoc:     uniformLocation(c, "pointLight.Color"),
		uniformLightLinearLoc:    uniformLocation(c, "pointLight.Linear"),
		uniformLightQuadraticLoc: uniformLocation(c, "pointLight.Quadratic"),

		locProjMatrixInv: uniformLocation(c, "projMatrixInv"),
		locViewMatrixInv: uniformLocation(c, "viewMatrixInv"),
	}
	return s
}

type PointLightShader struct {
	*DefaultShader
	locNormal       int32
	LocGAlbedo      int32
	uniformDepthLoc int32

	uniformLightPosLoc       int32
	uniformLightColorLoc     int32
	uniformLightLinearLoc    int32
	uniformLightQuadraticLoc int32

	locProjMatrixInv int32
	locViewMatrixInv int32
}

func (s *PointLightShader) UniformNormalLoc() int32 {
	return s.locNormal
}

func (s *PointLightShader) UniformAlbedoSpecLoc() int32 {
	return s.LocGAlbedo
}

func (s *PointLightShader) UniformDepthLoc() int32 {
	return s.uniformDepthLoc
}

func (s *PointLightShader) SetLight(light *PointLight) {
	gl.Uniform3f(s.uniformLightPosLoc, light.Position[0], light.Position[1], light.Position[2])
	gl.Uniform3f(s.uniformLightColorLoc, light.Color[0], light.Color[1], light.Color[2])
	gl.Uniform1f(s.uniformLightLinearLoc, light.Linear)
	gl.Uniform1f(s.uniformLightQuadraticLoc, light.Exp)
}
