package utils

import "github.com/reimirno/golinks/pkg/types"

func Paginate[T any](list []T, pagination types.Pagination) []T {
	if pagination.Offset >= len(list) {
		return nil
	}
	rightBound := min(pagination.Offset+pagination.Limit, len(list))
	return list[pagination.Offset:rightBound]
}

var DefaultPagination = types.Pagination{Offset: 0, Limit: 100}
