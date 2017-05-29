// Package tags contains code for parsing struct tags.
package tags

import (
	"strconv"
	"strings"
)

// A StructTag aliases a struct tag string to add parsing functionality.
type StructTag string

// A Handler function handles a tag key/value pair.
// Returns (true, nil) if the tag is handled successfully.
type Handler func(k, v string) (bool, error)

// ForEach parses tag and iterates over the key/value pairs, passing each to handler.
// Iteration may be terminated early if handler returns (true, nil).
// Derived from reflect/type.go Get
func (tag StructTag) ForEach(handler Handler) error {
	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, _ := strconv.Unquote(qvalue)

		if handled, err := handler(name, value); err != nil {
			return err
		} else if handled {
			// Tag was handled. We're done.
			return nil
		}
		// Tag was not handled. Continue.
	}
	// No tags were handled.
	return nil
}

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns ("", false).
// Similar to reflect/type.go Get, but distinguishes between empty tag values and missing tag keys.
func (tag StructTag) Get(key string) (string, bool) {
	value, ok := "", false
	tag.ForEach(Handler(func(tagKey, v string) (bool, error) {
		if tagKey == key {
			value, ok = v, true
			return true, nil
		}
		return false, nil
	}))
	return value, ok
}

// TagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
// Copy of unexported type tagOptions from encoding/json/tags.go
type TagOptions string

// ParseTag splits a struct field's json tag into its name and
// comma-separated options.
func ParseTag(tag string) (string, TagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], TagOptions(tag[idx+1:])
	}
	return tag, TagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o TagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
