package tracklist

import (
	"github.com/gieseladev/tracklist/formats/celeste"
	"github.com/gieseladev/tracklist/formats/inception"
	"github.com/gieseladev/tracklist/formats/mili"
	"github.com/gieseladev/tracklist/formats/nier"
	"github.com/gieseladev/tracklist/formats/sawano"
)

// Parser is an interface for tracklist parsers.
type Parser interface {
	// Parse takes
	Parse(text string) (List, error)
}

var registeredParsers = []Parser{
	sawano.Parser,
	celeste.Parser,
	mili.Parser, nier.Parser,
	inception.Parser,
}

// Register adds the given Parsers to the registered parsers.
func Register(p ...Parser) {
	registeredParsers = append(registeredParsers, p...)
}

func All() []Parser {
	return registeredParsers
}
