package main

import "github.com/go-gl/gl/v4.1-core/gl"

type DirLightShader struct {
	*DefaultShader
	uniformDepthLoc      int32
	uniformNormalLoc     int32
	uniformAlbedoSpecLoc int32

	uniformLightDirectionLoc        int32
	uniformLightColorLoc            int32
	uniformLightDiffuseIntensityLoc int32

	locProjMatrixInv   int32
	locViewMatrixInv   int32
	locShadowMap       int32
	locLightProjection int32
	locLightView       int32
}

func (s *DirLightShader) UniformNormalLoc() int32 {
	return s.uniformNormalLoc
}

func (s *DirLightShader) UniformAlbedoSpecLoc() int32 {
	return s.uniformAlbedoSpecLoc
}

func (s *DirLightShader) UniformDepthLoc() int32 {
	return s.uniformDepthLoc
}

func (s *DirLightShader) SetLight(light *DirectionalLight) {
	gl.Uniform3f(s.uniformLightDirectionLoc, light.Direction[0], light.Direction[1], light.Direction[2])
	gl.Uniform3f(s.uniformLightColorLoc, light.Color[0], light.Color[1], light.Color[2])
}

func NewDirLightShader() *DirLightShader {
	c := NewDefaultShader("lighting_dir_pbr", "lighting_dir_pbr")
	s := &DirLightShader{
		DefaultShader:            c,
		uniformDepthLoc:          uniformLocation(c, "gDepth"),
		uniformNormalLoc:         uniformLocation(c, "gNormal"),
		uniformAlbedoSpecLoc:     uniformLocation(c, "gAlbedoSpec"),
		uniformLightDirectionLoc: uniformLocation(c, "dirLight.Direction"),
		uniformLightColorLoc:     uniformLocation(c, "dirLight.Color"),

		locProjMatrixInv:   uniformLocation(c, "projMatrixInv"),
		locViewMatrixInv:   uniformLocation(c, "viewMatrixInv"),
		locShadowMap:       uniformLocation(c, "shadowMap"),
		locLightProjection: uniformLocation(c, "lightProjection"),
		locLightView:       uniformLocation(c, "lightView"),
	}
	return s
}
