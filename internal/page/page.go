package page

import (
	"fmt"
	"strconv"
)

type Page struct {
	number int
	rows   int
}

func Parse(page string, rowsPerPage string) (Page, error) {
	number := 1
	if page != "" {
		var err error
		number, err = strconv.Atoi(page)
		if err != nil {
			return Page{}, fmt.Errorf("page conversion: %w", err)
		}
	}

	rows := 10
	if rowsPerPage != "" {
		var err error
		rows, err = strconv.Atoi(rowsPerPage)
		if err != nil {
			return Page{}, fmt.Errorf("rows conversion: %w", err)
		}
	}

	//validation
	if number <= 0 {
		return Page{}, fmt.Errorf("page number too low, must be greater that 0: %d", number)
	}

	if rows <= 0 {
		return Page{}, fmt.Errorf("rows per page is too low, must be greater that 0: %d", rows)
	}

	if rows > 100 {
		return Page{}, fmt.Errorf("rows per page too high, must be less than 100: %d", rows)
	}

	return Page{number: number, rows: rows}, nil
}
