package fuzz

import "encoding/json"

func Fuzz(b []byte) int {
	var a struct {
		A int
		B string
	}
	if err := json.Unmarshal(b, &a); err != nil {
		return 0
	}
	if a.A < 10 {
		panic(a.A)
	}
	if len(a.B) > 10 {
		panic(a.B)
	}
	return 1
}
