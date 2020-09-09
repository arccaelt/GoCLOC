package types

type ParserError struct {
	message string
}

func (p *ParserError) Error() string {
	return p.message
}
