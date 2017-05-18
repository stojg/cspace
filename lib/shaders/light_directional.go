package shaders

/*
	s.dirLightShaderScreenSizeLoc = uniformLocation(s.dirLightShader, "gScreenSize")
	s.dirLightViewPosLoc = uniformLocation(s.dirLightShader, "viewPos")
	s.dirLightModelLoc = uniformLocation(s.dirLightShader, "model")

*/
func NewDirectionalLight() *DirectionalLight {
	c := buildShader("lighting_dir_pbr", "lighting_dir_pbr")
	s := &DirectionalLight{
		Program: c,

		LocGDepth:  loc(c, "gDepth"),
		LocGNormal: loc(c, "gNormal"),
		LocGAlbedo: loc(c, "gAlbedoSpec"),

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

	LocGDepth  int32
	LocGNormal int32
	LocGAlbedo int32

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

//func (s *DirectionalLight) SetLight(light *DirectionalLight) {
//	gl.Uniform3f(s.LocLightDirection, light.Direction[0], light.Direction[1], light.Direction[2])
//	gl.Uniform3f(s.LocLightColor, light.Color[0], light.Color[1], light.Color[2])
//}
