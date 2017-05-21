package ucreg

import "regexp"

type Parser struct {
	pattern string
	reg     *regexp.Regexp
}

func NewParser(pattern string) *Parser {
	p := &Parser{}
	p.reg = regexp.MustCompile(pattern)
	p.pattern = pattern
	return p
}

func (this *Parser) Match(str string) bool {
	if this.reg == nil {
		return false
	}
	return this.reg.MatchString(str)
}

func (this *Parser) Find(str string) []string {
	result := make([]string, 0)
	if this.reg == nil {
		return result
	}

	submatch := this.reg.FindSubmatch([]byte(str))
	for i := 0; i < len(submatch); i++ {
		result = append(result, string(submatch[i]))
	}
	return result
}
