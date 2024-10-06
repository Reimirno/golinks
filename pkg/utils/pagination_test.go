package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reimirno/golinks/pkg/types"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name       string
		list       []int
		pagination types.Pagination
		want       []int
	}{
		{
			name:       "Normal case",
			list:       []int{1, 2, 3, 4, 5},
			pagination: types.Pagination{Offset: 1, Limit: 2},
			want:       []int{2, 3},
		},
		{
			name:       "Offset exceeds list length",
			list:       []int{1, 2, 3},
			pagination: types.Pagination{Offset: 5, Limit: 2},
			want:       []int{},
		},
		{
			name:       "Limit exceeds remaining items",
			list:       []int{1, 2, 3, 4, 5},
			pagination: types.Pagination{Offset: 3, Limit: 10},
			want:       []int{4, 5},
		},
		{
			name:       "Empty list",
			list:       []int{},
			pagination: types.Pagination{Offset: 0, Limit: 5},
			want:       []int{},
		},
		{
			name:       "Zero limit",
			list:       []int{1, 2, 3, 4, 5},
			pagination: types.Pagination{Offset: 1, Limit: 0},
			want:       []int{},
		},
		{
			name:       "Offset at last element",
			list:       []int{1, 2, 3},
			pagination: types.Pagination{Offset: 2, Limit: 2},
			want:       []int{3},
		},
		{
			name:       "Negative offset",
			list:       []int{1, 2, 3, 4, 5},
			pagination: types.Pagination{Offset: -1, Limit: 2},
			want:       []int{1, 2},
		},
		{
			name:       "Negative limit",
			list:       []int{1, 2, 3, 4, 5},
			pagination: types.Pagination{Offset: 1, Limit: -2},
			want:       []int{},
		},
		{
			name:       "Offset and limit both zero",
			list:       []int{1, 2, 3, 4, 5},
			pagination: types.Pagination{Offset: 0, Limit: 0},
			want:       []int{},
		},
		{
			name:       "Offset and limit exceed list length",
			list:       []int{1, 2, 3},
			pagination: types.Pagination{Offset: 5, Limit: 10},
			want:       []int{},
		},
		{
			name:       "Large offset and limit values",
			list:       []int{1, 2, 3, 4, 5},
			pagination: types.Pagination{Offset: 1000000, Limit: 1000000},
			want:       []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Paginate(tt.list, tt.pagination)
			assert.Equal(t, tt.want, got)
		})
	}
}
