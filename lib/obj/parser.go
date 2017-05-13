package obj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/jonnenauha/obj-simplify/objectfile"
)

var (
	ObjectsParsed int
	GroupsParsed  int
)

func NewMaterial() *Material {
	return &Material{
		Name:         "DefaultMaterial",
		Ambient:      [3]float32{0.5, 0.5, 0.5},
		Diffuse:      [3]float32{0.5, 0.5, 0.5},
		Specular:     [3]float32{0.5, 0.5, 0.5},
		SpecularExp:  32,
		Transparency: 1,
	}
}

type Material struct {
	Name         string
	Ambient      [3]float32
	Diffuse      [3]float32
	Specular     [3]float32
	Transparency float32
	SpecularExp  float32
}

type Object struct {
	Name string
	Data []float32
	Mtr  *Material
}

func LoadObject(filename string) []*Object {
	obj, num, err := ParseFile(filename)
	if err != nil {
		fmt.Println("error at object line: ", num)
		panic(err)
	}

	materials := make(map[string]*Material)

	for _, mtrlib := range obj.MaterialLibraries {
		ms, err := ParseMtr(filepath.Join(filepath.Dir(filename), mtrlib))
		if err != nil {
			fmt.Printf("warning: material parse: %s\n", err)
		}
		for k, m := range ms {
			materials[k] = m
		}
	}

	var objects []*Object

	for _, object := range obj.Objects {
		var data []float32
		// convert the face data into actual data ready for openGL loading
		for _, vert := range object.VertexData {
			data = add(data, vert.Declarations[0])
			data = add(data, vert.Declarations[1])
			data = add(data, vert.Declarations[2])
			for i := 3; i < len(vert.Declarations); i++ {
				data = add(data, vert.Declarations[i-3])
				data = add(data, vert.Declarations[i-1])
				data = add(data, vert.Declarations[i])
			}
		}

		mat := materials[object.Material]
		if mat == nil {
			mat = NewMaterial()
		}

		objects = append(objects, &Object{
			Name: object.Name,
			Data: data,
			Mtr:  mat,
		})
	}
	return objects
}

func add(data []float32, in *objectfile.Declaration) []float32 {
	data = appendValues(data, in.RefVertex, 3)
	data = appendValues(data, in.RefNormal, 3)
	if in.RefUV != nil {
		data = appendValues(data, in.RefUV, 2)
	} else {
		data = append(data, 0, 0)
	}
	return data
}

func appendValues(data []float32, in *objectfile.GeometryValue, count int) []float32 {
	return append(data, toFloat32(in)[:count]...)
}

func ParseFile(path string) (*objectfile.OBJ, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, -1, err
	}
	defer f.Close()
	return parse(f)
}

func parse(src io.Reader) (*objectfile.OBJ, int, error) {
	dest := objectfile.NewOBJ()
	geom := dest.Geometry

	scanner := bufio.NewScanner(src)
	linenum := 0

	var (
		currentObject           *objectfile.Object
		currentObjectName       string
		currentObjectChildIndex int
		currentMaterial         string
		currentSmoothGroup      string
	)

	fakeObject := func(material string) *objectfile.Object {
		ot := objectfile.ChildObject
		if currentObject != nil {
			ot = currentObject.Type
		}
		currentObjectChildIndex++
		name := fmt.Sprintf("%s_%d", currentObjectName, currentObjectChildIndex)
		return dest.CreateObject(ot, name, material)
	}

	for scanner.Scan() {
		linenum++

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		t, value := parseLineType(line)

		// Force GC and release mem to OS for >1 million
		// line source files, every million lines.
		//
		// @todo We should also do data structure optimizations to handle
		// multiple gig source files without swapping on low mem machines.
		// A 4.5gb 82 million line test source file starts swapping on my 8gb
		// mem machine (though this app used ~5gb) at about the 40 million line mark.
		//
		// Above should be done when actualy users have a real use case for such
		// large files :)
		if linenum%1000000 == 0 {
			rt := time.Now()
			debug.FreeOSMemory()
			fmt.Println("%s lines parsed - Forced GC took %s", rt)
		}

		switch t {

		// comments
		case objectfile.Comment:
			if currentObject == nil && len(dest.MaterialLibraries) == 0 {
				dest.Comments = append(dest.Comments, value)
			} else if currentObject != nil {
				// skip comments that might refecence vertex, normal, uv, polygon etc.
				// counts as they wont be most likely true after this tool is done.
				if len(value) > 0 && !strContainsAny(value, []string{"vertices", "normals", "uvs", "texture coords", "polygons", "triangles"}, caseInsensitive) {
					currentObject.Comments = append(currentObject.Comments, value)
				}
			}

			// mtl file ref
		case objectfile.MtlLib:
			dest.MaterialLibraries = append(dest.MaterialLibraries, value)

			// geometry
		case objectfile.Vertex, objectfile.Normal, objectfile.UV, objectfile.Param:
			if _, err := geom.ReadValue(t, value, true); err != nil {
				return nil, linenum, wrapErrorLine(err, linenum)
			}

			// object, group
		case objectfile.ChildObject, objectfile.ChildGroup:
			currentObjectName = value
			currentObjectChildIndex = 0
			// inherit currently declared material
			currentObject = dest.CreateObject(t, currentObjectName, currentMaterial)
			if t == objectfile.ChildObject {
				ObjectsParsed++
			} else if t == objectfile.ChildGroup {
				GroupsParsed++
			}

			// object: material
		case objectfile.MtlUse:
			// obj files can define multiple materials inside a single object/group.
			// usually these are small face groups that kill performance on 3D engines
			// as they have to render hundreds or thousands of meshes with the same material,
			// each mesh containing a few faces.
			//
			// this app will convert all these "multi material" objects into
			// separate object, later merging all meshes with the same material into
			// a single draw call geometry.
			//
			// this might be undesirable for certain users, renderers and authoring software,
			// in this case don't use this simplified on your obj files. simple as that.

			// only fake if an object has been declared
			if currentObject != nil {
				// only fake if the current object has declared vertex data (faces etc.)
				// and the material name actually changed (ecountering the same usemtl
				// multiple times in a row would be rare, but check for completeness)
				if len(currentObject.VertexData) > 0 && currentObject.Material != value {
					currentObject = fakeObject(value)
				}
			}

			// store material value for inheriting
			currentMaterial = value

			// set material to current object
			if currentObject != nil {
				currentObject.Material = currentMaterial
			}

			// object: faces
		case objectfile.Face, objectfile.Line, objectfile.Point:
			// most tools support the file not defining a o/g prior to face declarations.
			// I'm not sure if the spec allows not declaring any o/g.
			// Our data structures and parsing however requires objects to put the faces into,
			// create a default object that is named after the input file (without suffix).
			if currentObject == nil {
				currentObject = dest.CreateObject(objectfile.ChildObject, "default", currentMaterial)
			}
			vd, vdErr := currentObject.ReadVertexData(t, value, true)
			if vdErr != nil {
				return nil, linenum, wrapErrorLine(vdErr, linenum)
			}
			// attach current smooth group and reset it
			if len(currentSmoothGroup) > 0 {
				vd.SetMeta(objectfile.SmoothingGroup, currentSmoothGroup)
				currentSmoothGroup = ""
			}

		case objectfile.SmoothingGroup:
			// smooth group can change mid vertex data declaration
			// so it is attched to the vertex data instead of current object directly
			currentSmoothGroup = value

			// unknown
		case objectfile.Unkown:
			return nil, linenum, wrapErrorLine(fmt.Errorf("Unsupported line %q\n\nPlease submit a bug report. If you can, provide this file as an attachement.\n> %s\n", line, "https://github.com/jonnenauha/obj-simplify//issues"), linenum)
		default:
			return nil, linenum, wrapErrorLine(fmt.Errorf("Unsupported line %q\n\nPlease submit a bug report. If you can, provide this file as an attachement.\n> %s\n", line, "https://github.com/jonnenauha/obj-simplify//issues"), linenum)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, linenum, err
	}
	return dest, linenum, nil
}

func wrapErrorLine(err error, linenum int) error {
	return fmt.Errorf("line:%d %s", linenum, err.Error())
}

func parseLineType(str string) (objectfile.Type, string) {
	value := ""
	if i := strings.Index(str, " "); i != -1 {
		value = strings.TrimSpace(str[i+1:])
		str = str[0:i]
	}
	return objectfile.TypeFromString(str), value
}

func toFloat32(val *objectfile.GeometryValue) []float32 {
	return []float32{float32(val.X), float32(val.Y), float32(val.Z), float32(val.Z)}
}

type caseSensitivity int

const (
	caseSensitive   caseSensitivity = 0
	caseInsensitive caseSensitivity = 1
)

func strContains(str, part string, cs caseSensitivity) bool {
	if cs == caseSensitive {
		return strings.Contains(str, part)
	}
	return strings.Contains(strings.ToLower(str), strings.ToLower(part))
}

func strContainsAny(str string, parts []string, cs caseSensitivity) bool {
	for _, part := range parts {
		if strContains(str, part, cs) {
			return true
		}
	}
	return false
}

func ParseMtr(path string) (map[string]*Material, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseMtr(f)
}

func parseMtr(r io.Reader) (map[string]*Material, error) {
	result := make(map[string]*Material)
	current := ""
	var f1, f2, f3 float32
	reader := bufio.NewReader(r)
	line, e1 := reader.ReadString('\n')
	for ; e1 == nil; line, e1 = reader.ReadString('\n') {
		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "Ka": // ambient
			if _, e := fmt.Sscanf(line, "Ka %f %f %f", &f1, &f2, &f3); e != nil {
				return result, fmt.Errorf("could not parse ambient values %s", e)
			}
			result[current].Ambient = [3]float32{f1, f2, f3}
		case "Kd": // diffuse
			if _, e := fmt.Sscanf(line, "Kd %f %f %f", &f1, &f2, &f3); e != nil {
				return result, fmt.Errorf("could not parse diffuse values %s", e)
			}
			result[current].Diffuse = [3]float32{f1, f2, f3}
		case "Ks": // specular
			if _, e := fmt.Sscanf(line, "Ks %f %f %f", &f1, &f2, &f3); e != nil {
				return result, fmt.Errorf("could not parse specular values %s", e)
			}
			result[current].Specular = [3]float32{f1, f2, f3}
		case "d": // transparency
			a, _ := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 32)
			result[current].Transparency = float32(a)
		case "Ns": // specular exponent
			ns, _ := strconv.ParseFloat(strings.TrimSpace(tokens[1]), 32)
			result[current].SpecularExp = float32(ns)
		case "newmtl": // material name
			current = strings.TrimSpace(tokens[1])
			result[current] = &Material{
				Name: tokens[1],
			}
		case "Ni": // optical density - scaler. Ignored for now.
		case "illum": // illumination model - int. Ignored for now.
		}
	}
	return result, nil
}
