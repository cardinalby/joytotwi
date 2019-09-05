package rssparser

import "joytotwi/app/joy/parsers"

// GetParserTraits of page parser
func GetParserTraits() parsers.Traits {
	return parsers.Traits{
		ID:                 "rss",
		SupportsForward:    true,
		SupportsReverse:    false,
		MaxOffsetFromStart: 10,
		MaxOffsetFromEnd:   10,
	}
}
