package shaders

func NewDirectionalLight() *DirectionalLight {
	c := buildShader("lighting_dir_pbr", "lighting_dir_pbr")
	s := &DirectionalLight{
		Program: c,

		LocGDepth:            loc(c, "gDepth"),
		LocGNormal:           loc(c, "gNormal"),
		LocGAlbedo:           loc(c, "gAlbedoSpec"),
		LocGAmbientOcclusion: loc(c, "gAmbientOcclusion"),

		LocLightDirection: loc(c, "dirLight.Direction"),
		LocLightColor:     loc(c, "dirLight.Color"),

		LocProjMatrixInv:   loc(c, "projMatrixInv"),
		LocViewMatrixInv:   loc(c, "viewMatrixInv"),
		LocShadowMap:       loc(c, "shadowMap"),
		LocLightProjection: loc(c, "lightProjection"),
		LocLightView:       loc(c, "lightView"),
		LocScreenSize:      loc(c, "gScreenSize"),
		LocViewPos:         loc(c, "viewPos"),
	}
	return s
}

type DirectionalLight struct {
	Program uint32

	LocGDepth            int32
	LocGNormal           int32
	LocGAlbedo           int32
	LocGAmbientOcclusion int32

	LocLightDirection  int32
	LocLightColor      int32
	LocProjMatrixInv   int32
	LocViewMatrixInv   int32
	LocShadowMap       int32
	LocLightProjection int32
	LocLightView       int32
	LocScreenSize      int32
	LocViewPos         int32
}

func (s *DirectionalLight) UniformNormalLoc() int32 {
	return s.LocGNormal
}

func (s *DirectionalLight) UniformAlbedoSpecLoc() int32 {
	return s.LocGAlbedo
}

func (s *DirectionalLight) UniformDepthLoc() int32 {
	return s.LocGDepth
}
