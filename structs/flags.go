package structs

import "strings"

type CustomFlag struct {
	value   []string
	changed bool
}

// String returns a string representation of the CustomFlag value.
// The elements of the value slice are joined by comma separators. If the value slice is empty,
// an empty string is returned.
func (c *CustomFlag) String() string {
	return strings.Join(c.value, ",")
}

// Set sets the value of the CustomFlag based on the provided string.
// If the string is empty, it sets the value to an empty string slice and marks the flag as changed.
// Otherwise, it splits the string by commas and sets the value to the resulting string slice, marking the flag as changed.
// It always returns nil.
func (c *CustomFlag) Set(s string) error {
	if s == "" {
		c.changed = true
		c.value = []string{}
	} else {
		c.changed = true
		c.value = strings.Split(s, ",")
	}

	return nil
}

// Get returns the value of the CustomFlag as a slice of strings.
func (c *CustomFlag) Get() []string {
	return c.value
}

// IsChanged returns a boolean value indicating whether the CustomFlag has been changed.
func (c *CustomFlag) IsChanged() bool {
	return c.changed
}

// NewCustomFlag creates a new instance of the CustomFlag struct and returns a pointer to it.
func NewCustomFlag() *CustomFlag {
	return &CustomFlag{}
}
