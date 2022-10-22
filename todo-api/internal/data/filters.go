//filename : internal/data/filters.go

package data

import (
	"math"
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

// The limit() method determines the LIMIT
func (f Filters) limit() int {
	return f.PageSize
}

// The offset() method calculates the OFFSET
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// The Metadata type contains metadata to help with pagination
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// The calculateMetadata() function computes the values for the Metadata fields
func calculateMetadata(totalRecrods int, page int, pageSize int) Metadata {
	if totalRecrods == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecrods) / float64(pageSize))),
		TotalRecords: totalRecrods,
	}
}
