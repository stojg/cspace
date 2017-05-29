package shaders

type IBLPreFilter struct {
	Program           uint32
	LocEnvironmentMap int32
	LocProjection     int32
	LocView           int32
	LocRoughness      int32
}

func NewIBLPrefilter() *IBLPreFilter {
	c := buildShader("ibl_cubemap", "ibl_prefilter")

	return &IBLPreFilter{
		Program:           c,
		LocEnvironmentMap: loc(c, "environmentMap"),
		LocProjection:     loc(c, "projection"),
		LocView:           loc(c, "view"),
		LocRoughness:      loc(c, "roughness"),
	}
}
