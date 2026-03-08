package pagination

import "net/http"

type SortParams struct {
	Column string
	Order  string
}

func ParseSort(r *http.Request, allowedColumns []string, defaultColumn string) SortParams {
	column := r.URL.Query().Get("column")
	order := r.URL.Query().Get("order")

	// Validate column against whitelist (prevent SQL injection)
	valid := false
	for _, col := range allowedColumns {
		if col == column {
			valid = true
			break
		}
	}
	if !valid {
		column = defaultColumn
	}

	// Only allow ASC and DESC
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return SortParams{
		Column: column,
		Order:  order,
	}
}
