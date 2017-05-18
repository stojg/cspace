package shaders

func NewPointLightShader() *PointLight {
	c := buildShader("lighting", "lighting_point_pbr")
	s := &PointLight{
		Program: c,

		LocModel:      loc(c, "model"),
		LocView:       loc(c, "view"),
		LocProjection: loc(c, "projection"),

		LocGNormal: loc(c, "gNormal"),
		LocGAlbedo: loc(c, "gAlbedoSpec"),
		LocGDepth:  loc(c, "gDepth"),

		LocLightPos:       loc(c, "pointLight.Position"),
		LocLightColor:     loc(c, "pointLight.Color"),
		LocLightLinear:    loc(c, "pointLight.Linear"),
		LocLightQuadratic: loc(c, "pointLight.Quadratic"),

		LocProjMatrixInv: loc(c, "projMatrixInv"),
		LocViewMatrixInv: loc(c, "viewMatrixInv"),
		LocScreenSize:    loc(c, "gScreenSize"),
		LocViewPos:       loc(c, "viewPos"),
	}
	return s
}

type PointLight struct {
	Program uint32

	LocModel          int32
	LocView           int32
	LocProjection     int32
	LocGNormal        int32
	LocGAlbedo        int32
	LocGDepth         int32
	LocLightPos       int32
	LocLightColor     int32
	LocLightLinear    int32
	LocLightQuadratic int32
	LocScreenSize     int32
	LocViewPos        int32
	LocProjMatrixInv  int32
	LocViewMatrixInv  int32
}

//func (s *PointLight) SetLight(light *PointLight) {
//
//}
