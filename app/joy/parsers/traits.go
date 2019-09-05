package parsers

// TraitAnyOffset is a value for ParserTrairs.MaxOffsetFromStart/MaxOffsetFromEnd
const TraitAnyOffset = -1

// GetTraits interface of function to get parser traits
type GetTraits func() Traits

// Traits definition of Parser
type Traits struct {
	ID                 string
	SupportsForward    bool
	SupportsReverse    bool
	MaxOffsetFromStart int
	MaxOffsetFromEnd   int
}
