package utils

import "github.com/reimirno/golinks/pkg/types"

func Paginate[T any](list []T, pagination types.Pagination) []T {
	if pagination.Offset >= len(list) {
		return []T{}
	}
	if pagination.Offset < 0 {
		pagination.Offset = 0
	}
	rightBound := min(pagination.Offset+pagination.Limit, len(list))
	if rightBound < pagination.Offset {
		return []T{}
	}
	return list[pagination.Offset:rightBound]
}

var DefaultPagination = types.Pagination{Offset: 0, Limit: 100}
