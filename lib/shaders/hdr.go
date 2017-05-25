package shaders

type HDR struct {
	Program          uint32
	LocScreenTexture int32
	LocExposure      int32
}

func NewHDR() *HDR {
	c := buildShader("fx", "fx_tone")
	return &HDR{
		Program:          c,
		LocScreenTexture: loc(c, "screenTexture"),
		LocExposure:      loc(c, "exposure"),
	}
}
