/*
Package scanner provides Scanner interface and it's default implementation,
which can be used to scan key-value pairs from property files or struct tags.

This scanner implementation is an alternative to "reflect".StructTag.Get(),
more convenient to parse complex tag values.

It gives:

  - multiline tags,
  - custom separators between key and value parts,
  - custom quote characters,
  - custom escape character,
  - more relaxed syntax (have a look at tests to get an idea),
  - conventional syntax is also supported.

Instead of Get/Lookup interface, scanner returns a map of key-value pairs.

The same scanner can be used to parse short and simple property files, similar
in syntax to Java property files. Though, this is not a main goal of this
package. So, synatx only resambles one of Java property files, there can be
many differencies.

Scanner should work with unicode within files/tags, but this is not tested yet.
*/
package scanner

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Scanner is an interface of key-value pairs scanner.
type Scanner interface {
	Tags(reflect.StructTag) (map[string]string, error)
	Scan(io.Reader) (map[string]string, error)
}

// New creates Scanner intstance with custom field separators, quotes and
// escape character.
func New(fieldSeparators, quotes []rune, escape rune) Scanner {
	return &scanner{
		splitter{
			fieldSeparators: fieldSeparators,
			quotes:          quotes,
			esc:             escape,
			noesc:           utf8.RuneError, // assume shouldn't exist
		},
	}
}

// Default is an instance of Scanner with sensible defaults.
var Default = New([]rune{':', '='}, []rune{'"', '\'', '`'}, '\\')

type scanner struct {
	spl splitter
}

func (s scanner) Tags(tag reflect.StructTag) (map[string]string, error) {
	return s.Scan(strings.NewReader(string(tag)))
}

func (s scanner) Scan(r io.Reader) (map[string]string, error) {
	values := make(map[string]string)
	bs := bufio.NewScanner(r)
	k := ""
	bs.Split(s.spl.split)
	for bs.Scan() {
		t := bs.Text()
		switch s.spl.tokenType {
		case tokenTypeKey:
			k = t
		case tokenTypeValue:
			if k != "" {
				values[k] = t
				k = ""
			}
		}
	}
	if err := bs.Err(); err != nil {
		return values, err
	}
	return values, nil
}

const (
	tokenTypeSpace = iota
	tokenTypeFieldSeparator
	tokenTypeKey
	tokenTypeValue
	tokenTypeText
)

type splitter struct {
	fieldSeparators []rune
	quotes          []rune
	esc             rune
	noesc           rune

	tokenType int
	waitValue bool
}

func (s *splitter) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < 2 {
		return 0, nil, nil
	}
	c := data[0]
	r := rune(c)
	isFieldSeparator := isOneOf(s.fieldSeparators...)
	isQuote := isOneOf(s.quotes...)
	switch {
	case unicode.IsSpace(r): // space
		advance, token, err = consume(data, s.noesc, unicode.IsSpace)
		s.tokenType = tokenTypeSpace
	case isFieldSeparator(r): // field separator
		advance, token, err = consume(data, s.noesc, isFieldSeparator)
		s.tokenType = tokenTypeFieldSeparator
		s.waitValue = true
	case s.waitValue: // value
		if isQuote(r) {
			data = data[1:]
			advance, token, err = consume(data, s.esc, func(rr rune) bool { return rr != r })
			advance += 2
		} else {
			advance, token, err = consume(data, s.esc, not(isOneOf('\r', '\n')))
		}
		s.tokenType = tokenTypeValue
		s.waitValue = false
	case isKeyRune(r): // key
		advance, token, err = consume(data, s.noesc, isKeyRune)
		s.tokenType = tokenTypeKey
	default: // some non-space garbage
		advance, token, err = consume(data, s.noesc, not(unicode.IsSpace))
		s.tokenType = tokenTypeText
	}
	return
}

func consume(data []byte, escape rune, f func(rune) bool) (int, []byte, error) {
	tok := []byte{}
	skip := false
	for i, b := range string(data) {
		if b == escape {
			skip = true
			continue
		}
		if f(b) {
			if skip { // current character shouldn't be escaped, assume '\' literaly
				tok = append(tok, []byte(string(escape))...)
			}
			tok = append(tok, []byte(string(b))...)
			skip = false
			continue
		}
		if skip {
			tok = append(tok, []byte(string(b))...)
			skip = false
			continue
		}
		return i, tok, nil
	}
	return len(tok), tok, bufio.ErrFinalToken
}

func isKeyRune(r rune) bool {
	switch r {
	case '_', '$':
		return true
	}
	return unicode.In(r, unicode.Letter, unicode.Digit, unicode.Dash, unicode.Hyphen)
}

func isOneOf(rns ...rune) func(rune) bool {
	return func(r rune) bool {
		for _, rr := range rns {
			if rr == r {
				return true
			}
		}
		return false
	}
}

func not(f func(rune) bool) func(rune) bool {
	return func(r rune) bool { return !f(r) }
}
