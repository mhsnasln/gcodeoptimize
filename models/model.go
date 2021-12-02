package models

type Layer struct {
	Items []Point `json:"items"`
}

type Block struct {
	Parts []Part
}

type Part struct {
	Z     float64
	Items []Point
}

type Point struct {
	G    int32   `json:"g"`
	X    float64 `json:"x"`
	XInc float64 `json:"x_inc"`
	Y    float64 `json:"y"`
	Z    float64 `json:"Z"`
	M1   float64 `json:"m1"`
	M2   float64 `json:"m2"`
}
