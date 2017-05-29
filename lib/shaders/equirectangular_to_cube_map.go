package shaders

type EquiRectToCubeMap struct {
	Program               uint32
	LocEquirectangularMap int32
	LocProjection         int32
	LocView               int32
}

func NewEquiRectToCubeMap() *EquiRectToCubeMap {
	c := buildShader("ibl_cubemap", "equirectangular_to_cubemap")

	return &EquiRectToCubeMap{
		Program:               c,
		LocEquirectangularMap: loc(c, "equirectangularMap"),
		LocProjection:         loc(c, "projection"),
		LocView:               loc(c, "view"),
	}
}
