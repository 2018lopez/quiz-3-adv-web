//filename : internal/data/filters.go

package data

import (
	"strings"

	"todo.imerlopez.net/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
	SortList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	//Check page and pageSize params
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 1000, "page", "must be a maximum of 1000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.Page <= 100, "page_size", "must be a maximum of 100")

	//check that the sort params matches a values in the acceptable sort list
	v.Check(validator.In(f.Sort, f.SortList...), "sort", "invalid sort value")
}

// The sortColumn() method safety extracted the sort field query parameter
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// the sortOrder() determine by asc or desc
func (f Filters) sortOrder() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}
