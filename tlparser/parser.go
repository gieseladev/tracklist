package tlparser

import (
	"github.com/gieseladev/tracklist/tlparser/celeste"
	"github.com/gieseladev/tracklist/tlparser/common"
	"github.com/gieseladev/tracklist/tlparser/inception"
	"github.com/gieseladev/tracklist/tlparser/mili"
	"github.com/gieseladev/tracklist/tlparser/nier"
	"github.com/gieseladev/tracklist/tlparser/sawano"
)

// Parser is an interface for tracklist parsers.
type Parser interface {
	// Parse takes
	Parse(text string) (common.List, error)
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
