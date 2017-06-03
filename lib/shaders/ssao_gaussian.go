package shaders

type SSAOGaussian struct {
	Program          uint32
	LocScreenTexture int32
	LocHorizontal    int32
}

func NewSSAOGaussian() *SSAOGaussian {
	c := buildShader("fx", "ssao_blur")

	return &SSAOGaussian{
		Program:          c,
		LocScreenTexture: loc(c, "screenTexture"),
		LocHorizontal:    loc(c, "horizontal"),
	}
}
