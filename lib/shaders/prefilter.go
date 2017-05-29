package shaders

type Prefilter struct {
	Program           uint32
	LocEnvironmentMap int32
	LocProjection     int32
	LocView           int32
	LocRoughness      int32
}

func NewPrefilter() *Prefilter {
	c := buildShader("equi", "equi_prefilter")

	return &Prefilter{
		Program:           c,
		LocEnvironmentMap: loc(c, "environmentMap"),
		LocProjection:     loc(c, "projection"),
		LocView:           loc(c, "view"),
		LocRoughness:      loc(c, "roughness"),
	}
}
