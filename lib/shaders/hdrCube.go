package shaders

type HDRCube struct {
	Program               uint32
	LocEquirectangularMap int32
	LocProjection         int32
	LocView               int32
}

func NewHDRCube() *HDRCube {
	c := buildShader("equi", "equi")

	return &HDRCube{
		Program:               c,
		LocEquirectangularMap: loc(c, "equirectangularMap"),
		LocProjection:         loc(c, "projection"),
		LocView:               loc(c, "view"),
	}
}
