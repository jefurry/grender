package simplify

import (
	"bufio"
	"encoding/binary"
	"os"
	"strings"
)

type STLHeader struct {
	_     [80]uint8
	Count uint32
}

type STLTriangle struct {
	N, V1, V2, V3 [3]float32
	_             uint16
}

func LoadBinarySTL(path string) (*Mesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	header := STLHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	count := int(header.Count)
	triangles := make([]*Triangle, count)
	for i := 0; i < count; i++ {
		d := STLTriangle{}
		if err := binary.Read(file, binary.LittleEndian, &d); err != nil {
			return nil, err
		}
		v1 := Vector{float64(d.V1[0]), float64(d.V1[1]), float64(d.V1[2])}
		v2 := Vector{float64(d.V2[0]), float64(d.V2[1]), float64(d.V2[2])}
		v3 := Vector{float64(d.V3[0]), float64(d.V3[1]), float64(d.V3[2])}
		triangles[i] = NewTriangle(v1, v2, v3)
	}
	return NewMesh(triangles), nil
}

func SaveBinarySTL(path string, mesh *Mesh) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	header := STLHeader{}
	header.Count = uint32(len(mesh.Triangles))
	if err := binary.Write(file, binary.LittleEndian, &header); err != nil {
		return err
	}
	for _, triangle := range mesh.Triangles {
		n := triangle.Normal()
		d := STLTriangle{}
		d.N[0] = float32(n.X)
		d.N[1] = float32(n.Y)
		d.N[2] = float32(n.Z)
		d.V1[0] = float32(triangle.V1.X)
		d.V1[1] = float32(triangle.V1.Y)
		d.V1[2] = float32(triangle.V1.Z)
		d.V2[0] = float32(triangle.V2.X)
		d.V2[1] = float32(triangle.V2.Y)
		d.V2[2] = float32(triangle.V2.Z)
		d.V3[0] = float32(triangle.V3.X)
		d.V3[1] = float32(triangle.V3.Y)
		d.V3[2] = float32(triangle.V3.Z)
		if err := binary.Write(file, binary.LittleEndian, &d); err != nil {
			return err
		}
	}
	return nil
}

func LoadSTL(path string) (*Mesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var vertexes []Vector
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 4 && fields[0] == "vertex" {
			f := parseFloats(fields[1:])
			v := Vector{f[0], f[1], f[2]}
			vertexes = append(vertexes, v)
		}
	}
	var triangles []*Triangle
	for i := 0; i < len(vertexes); i += 3 {
		v1 := vertexes[i+0]
		v2 := vertexes[i+1]
		v3 := vertexes[i+2]
		triangles = append(triangles, NewTriangle(v1, v2, v3))
	}
	return NewMesh(triangles), scanner.Err()
}
