package order

import (
	"fmt"
	"strings"
)

// set of directions for ordering
const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

type By struct {
	Field     string
	Direction string
}

func NewBy(field, direction string) By {
	if _, ok := directions[direction]; !ok {
		return By{Field: field, Direction: ASC} //default direction
	}

	return By{Field: field, Direction: direction}
}

func Parse(fieldMappings map[string]string, orderBy string, defaultOrder By) (By, error) {
	if orderBy == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(orderBy, ",")
	fieldNameRaw := strings.TrimSpace(orderParts[0])
	fieldName, ok := fieldMappings[fieldNameRaw]
	if !ok {
		return By{}, fmt.Errorf("unknows field: %q", fieldNameRaw)
	}

	switch len(orderParts) {
	case 1:
		return NewBy(fieldName, ASC), nil
	case 2:
		directionRaw := strings.TrimSpace(orderParts[1])
		direction, ok := directions[directionRaw]
		if !ok {
			return By{}, fmt.Errorf("unknows direction: %q", directionRaw)
		}
		return NewBy(fieldName, direction), nil
	default:
		return By{}, fmt.Errorf("invalid order by format: %q", orderBy)
	}
}
