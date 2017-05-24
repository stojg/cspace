package shaders

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewPointLightShader(lights int) *PointLight {
	c := buildShader("lighting_point_pbr", "lighting_point_pbr")
	s := &PointLight{
		Program: c,
		lights:  lights,

		LocGNormal:           loc(c, "gNormal"),
		LocGAlbedo:           loc(c, "gAlbedoSpec"),
		LocGDepth:            loc(c, "gDepth"),
		LocGAmbientOcclusion: loc(c, "gAmbientOcclusion"),

		LocNumLights: loc(c, "numLights"),

		LocScreenSize: loc(c, "gScreenSize"),
	}

	blockIndex := gl.GetUniformBlockIndex(c, gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(c, blockIndex, 0)

	for i := 0; i < lights; i++ {
		s.LocLightPos = append(s.LocLightPos, loc(c, fmt.Sprintf("pointLight[%d].Position", i)))
		s.LocLightColor = append(s.LocLightColor, loc(c, fmt.Sprintf("pointLight[%d].Color", i)))
		s.LocLightLinear = append(s.LocLightLinear, loc(c, fmt.Sprintf("pointLight[%d].Linear", i)))
		s.LocLightQuadratic = append(s.LocLightQuadratic, loc(c, fmt.Sprintf("pointLight[%d].Quadratic", i)))
	}

	return s
}

type PointLight struct {
	Program uint32
	lights  int

	LocGNormal           int32
	LocGAlbedo           int32
	LocGDepth            int32
	LocGAmbientOcclusion int32

	LocNumLights      int32
	LocLightPos       []int32
	LocLightColor     []int32
	LocLightLinear    []int32
	LocLightQuadratic []int32
	LocScreenSize     int32
}
