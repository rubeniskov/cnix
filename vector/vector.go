package vector

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Vector struct {
	X float64 `json:"x" regroup:"x"`
	Y float64 `json:"y" regroup:"y"`
	Z float64 `json:"z" regroup:"z"`
}

func (v Vector) Dot(o Vector) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector) Cross(o Vector) Vector {
	return Vector{
		X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X,
	}
}

func (v Vector) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Sum(o Vector) Vector {
	return Vector{
		X: v.X + o.X,
		Y: v.Y + o.Y,
		Z: v.Z + o.Z,
	}
}

func (v Vector) Diff(o Vector) Vector {
	return Vector{
		X: v.X - o.X,
		Y: v.Y - o.Y,
		Z: v.Z - o.Z,
	}
}

func (v Vector) Divide(d float64) Vector {
	return Vector{
		X: v.X / d,
		Y: v.Y / d,
		Z: v.Z / d,
	}
}

func (v Vector) String() string {
	return fmt.Sprintf("Vector{X: %f, Y: %f, Z: %f}", v.X, v.Y, v.Z)
}

func (v Vector) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%f, %f, %f", v.X, v.Y, v.Z)), nil
}

func (v *Vector) UnmarshalText(data []byte) error {
	var err error
	str := string(data)
	chnks := strings.Split(str, ",")
	
	if v.X, err = strconv.ParseFloat(chnks[0], 64); err != nil {
		return err
	}

	if v.Y, err = strconv.ParseFloat(chnks[0], 64); err != nil {
		return err
	}

	if v.Z, err = strconv.ParseFloat(chnks[0], 64); err != nil {
		return err
	} 
	
	return nil
}

func (v Vector) Marshal() ([]byte, error) {
	var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(v)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func (v *Vector) Unmarshal(data []byte) (error) {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(v)

	if err != nil {
		return err
	}
	return nil
}