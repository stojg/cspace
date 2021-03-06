package shaders

import "github.com/go-gl/gl/v4.1-core/gl"

func NewDirectionalLight() *DirectionalLight {
	c := buildShader("lighting_dir_pbr", "lighting_dir_pbr")
	s := &DirectionalLight{
		Program: c,

		LocGDepth:            loc(c, "gDepth"),
		LocGNormal:           loc(c, "gNormal"),
		LocGAlbedo:           loc(c, "gAlbedoSpec"),
		LocGAmbientOcclusion: loc(c, "gAmbientOcclusion"),
		LocIBLEnabled:        loc(c, "iblEnabled"),
		LocIrradianceMap:     loc(c, "irradianceMap"),
		LocPrefilterMap:      loc(c, "prefilterMap"),
		LocPbrdfLUT:          loc(c, "brdfLUT"),

		LocLightDirection: loc(c, "dirLight.Direction"),
		LocLightColor:     loc(c, "dirLight.Color"),
		LocLightEnabled:   loc(c, "dirLight.Enabled"),

		LocShadowMap:       loc(c, "shadowMap"),
		LocLightProjection: loc(c, "lightProjection"),
		LocLightView:       loc(c, "lightView"),
		LocScreenSize:      loc(c, "gScreenSize"),
	}

	blockIndex := gl.GetUniformBlockIndex(c, gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(c, blockIndex, 0)

	return s
}

type DirectionalLight struct {
	Program uint32

	LocGDepth            int32
	LocGNormal           int32
	LocGAlbedo           int32
	LocGAmbientOcclusion int32
	LocIBLEnabled        int32
	LocIrradianceMap     int32
	LocPrefilterMap      int32
	LocPbrdfLUT          int32

	LocLightDirection  int32
	LocLightColor      int32
	LocLightEnabled    int32
	LocShadowMap       int32
	LocLightProjection int32
	LocLightView       int32
	LocScreenSize      int32
}
