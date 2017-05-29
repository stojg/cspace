package shaders

type IBLIrradiance struct {
	Program           uint32
	LocEnvironmentMap int32
	LocProjection     int32
	LocView           int32
}

func NewIBLIrradiance() *IBLIrradiance {
	c := buildShader("ibl_cubemap", "ibl_irradiance_convolution")

	return &IBLIrradiance{
		Program:           c,
		LocEnvironmentMap: loc(c, "environmentMap"),
		LocProjection:     loc(c, "projection"),
		LocView:           loc(c, "view"),
	}
}
