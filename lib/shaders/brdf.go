package shaders

type Brdf struct {
	Program           uint32
	LocEnvironmentMap int32
	LocProjection     int32
	LocView           int32
	LocRoughness      int32
}

func NewBrdf() *Brdf {
	c := buildShader("brdf", "equi_brdf")

	return &Brdf{
		Program: c,
		//LocEnvironmentMap: loc(c, "environmentMap"),
	}
}
