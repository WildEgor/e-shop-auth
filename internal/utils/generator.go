package utils

import (
	"math/rand"
	"strconv"
)

type CodeGenerator struct {
	Length uint8
}

func NewCodeGenerator(length uint8) *CodeGenerator {
	return &CodeGenerator{
		Length: length,
	}
}

func (g *CodeGenerator) GenShortCode() string {
	num := rand.Intn(9000) + 1000
	return strconv.Itoa(num)
}
