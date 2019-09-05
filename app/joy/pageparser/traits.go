package pageparser

import (
	"joytotwi/app/joy/parsers"
)

// GetParserTraits of page parser
func GetParserTraits() parsers.Traits {
	return parsers.Traits{
		ID:                 "page",
		SupportsForward:    true,
		SupportsReverse:    true,
		MaxOffsetFromStart: parsers.TraitAnyOffset,
		MaxOffsetFromEnd:   parsers.TraitAnyOffset,
	}
}
