package tlparser

import (
	"github.com/gieseladev/tracklist/tlparser/celeste"
	"github.com/gieseladev/tracklist/tlparser/common"
	"github.com/gieseladev/tracklist/tlparser/inception"
	"github.com/gieseladev/tracklist/tlparser/mili"
	"github.com/gieseladev/tracklist/tlparser/sawano"
)

type Parser interface {
	Parse(text string) (common.List, error)
}

var registeredParsers = []Parser{
	sawano.Parser,
	celeste.Parser,
	mili.Parser,
	inception.Parser,
}

func Register(p ...Parser) {
	registeredParsers = append(registeredParsers, p...)
}

func All() []Parser {
	return registeredParsers
}
