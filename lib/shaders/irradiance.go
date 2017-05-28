package shaders

type Irradiance struct {
	Program           uint32
	LocEnvironmentMap int32
	LocProjection     int32
	LocView           int32
}

func NewIrradiance() *Irradiance {
	c := buildShader("equi", "equi_irradiance")

	return &Irradiance{
		Program:           c,
		LocEnvironmentMap: loc(c, "environmentMap"),
		LocProjection:     loc(c, "projection"),
		LocView:           loc(c, "view"),
	}
}
