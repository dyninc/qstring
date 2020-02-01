package qstring

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	operatorEquals    = "="
	operatorGreater   = ">"
	operatorGreaterEq = ">="
	operatorLesser    = "<"
	operatorLesserEq  = "<="
	operatorLike      = "~"
	operatorDifferent = "!"
)

// parseOperator parses a leading logical operator out of the provided string
func parseOperator(s string) string {
	switch s[0] {
	case operatorLesser[0]: // "<"
		if 1 == len(s) {
			return operatorLesser
		}
		switch s[1] {
		case operatorEquals[0]: // "="
			return operatorLesserEq
		default:
			return operatorLesser
		}
	case operatorGreater[0]: // ">"
		if 1 == len(s) {
			return operatorGreater
		}
		switch s[1] {
		case operatorEquals[0]: // "="
			return operatorGreaterEq
		default:
			return operatorGreater
		}
	case operatorLike[0]: // "~"
		return operatorLike
	case operatorDifferent[0]: // "!"
		return operatorDifferent
	default:
		// no operator found, default to "="
		return operatorEquals
	}
}

type ComparativeString struct {
	Operator string
	Str      string
}

// ComparativeTime is a field that can be used for specifying a query parameter
// which includes a conditional operator and a timestamp
type ComparativeTime struct {
	Operator string
	Time     time.Time
}

// NewComparativeTime returns a new ComparativeTime instance with a default
// operator of "="
func NewComparativeTime() *ComparativeTime {
	return &ComparativeTime{Operator: "="}
}

// Parse is used to parse a query string into a ComparativeTime instance
func (c *ComparativeTime) Parse(query string) error {
	if len(query) <= 2 {
		return errors.New("qstring: Invalid Timestamp Query")
	}

	c.Operator = parseOperator(query)

	// if no operator was provided and we defaulted to an equality operator
	if !strings.HasPrefix(query, c.Operator) {
		query = fmt.Sprintf("=%s", query)
	}

	var err error
	c.Time, err = time.Parse(time.RFC3339, query[len(c.Operator):])
	if err != nil {
		return err
	}

	return nil
}

// String returns this ComparativeTime instance in the form of the query
// parameter that it came in on
func (c ComparativeTime) String() string {
	return fmt.Sprintf("%s%s", c.Operator, c.Time.Format(time.RFC3339))
}

// Parse is used to parse a query string into a ComparativeString instance
func (c *ComparativeString) Parse(query string) error {
	if len(query) <= 2 {
		return errors.New("qstring: Invalid Query")
	}

	c.Operator = parseOperator(query)

	if c.Operator != operatorDifferent && c.Operator != operatorLike && c.Operator != operatorEquals {
		return errors.New(fmt.Sprintf("qstring: Invalid operator for %T", c))
	}
	if c.Operator == operatorEquals {
		c.Operator = ""
	}

	// if no operator was provided and we defaulted to an equality operator
	if !strings.HasPrefix(query, c.Operator) {
		query = fmt.Sprintf("=%s", query)
	}

	var err error
	c.Str = query[len(c.Operator):]
	if err != nil {
		return err
	}

	return nil
}

// String returns this ComparativeString instance in the form of the query
// parameter that it came in on
func (c ComparativeString) String() string {
	return fmt.Sprintf("%s%s", c.Operator, c.Str)
}
