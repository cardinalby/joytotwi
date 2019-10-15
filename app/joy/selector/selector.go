package selector

import (
	"errors"
	"fmt"
	"joytotwi/app/joy"
	"joytotwi/app/joy/pageparser"
	"joytotwi/app/joy/parsers"
	"joytotwi/app/joy/rssparser"
)

type postsBidirectionalReader func(
	userName string,
	reverse bool,
	offset int,
	limit int,
	done chan struct{},
) (chan *joy.Post, chan error)

type readerWithTraits struct {
	reader postsBidirectionalReader
	traits parsers.Traits
}

func getParsersTrairs() []readerWithTraits {
	return []readerWithTraits{
		{pageparser.GetPosts, pageparser.GetParserTraits()},
		{rssparser.GetPosts, rssparser.GetParserTraits()},
	}
}

// GetPostReader returns implementation specified in sourceType which should be one
// of the ParserID consts defined in parsers packages
func GetPostReader(
	sourceType string,
	reverse bool,
	offset int,
	limit int,
) (joy.PostsReader, error) {
	parserWithTrait := getReaderWithTraitByID(sourceType)
	if parserWithTrait == nil {
		return nil, fmt.Errorf("parser with '%s' not supported", sourceType)
	}
	traits := parserWithTrait.traits
	if !reverse && !traits.SupportsForward {
		return nil, errors.New("parser doesn't support forward parsing")
	}
	if reverse && !traits.SupportsReverse {
		return nil, errors.New("parser doesn't support reverse parsing")
	}
	if traits.MaxOffsetFromStart != parsers.TraitAnyOffset && offset > traits.MaxOffsetFromStart {
		return nil, fmt.Errorf("parser doesn't support reading out of %d offset", traits.MaxOffsetFromStart)
	}
	if traits.MaxOffsetFromEnd != parsers.TraitAnyOffset && offset+limit > traits.MaxOffsetFromEnd {
		return nil, fmt.Errorf("parser doesn't support reading out of %d offset", traits.MaxOffsetFromStart)
	}

	return func(userName string, done chan struct{}) (chan *joy.Post, chan error) {
		return parserWithTrait.reader(userName, reverse, offset, limit, done)
	}, nil
}

func getReaderWithTraitByID(id string) *readerWithTraits {
	for _, item := range getParsersTrairs() {
		if item.traits.ID == id {
			return &item
		}
	}
	return nil
}
