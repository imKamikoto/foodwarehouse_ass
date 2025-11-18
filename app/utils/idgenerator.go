package utils

import "fmt"

type IDGenerator struct {
	pattern string
	next    int
}

func NewGenerator(pattern string, iniValue int) *IDGenerator {
	return &IDGenerator{
		pattern: pattern,
		next:    iniValue,
	}
}

func (g *IDGenerator) Next() string {
	g.next++
	return fmt.Sprintf("%s - B%04d", g.pattern, g.next)
}
