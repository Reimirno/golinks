package types

import (
	"github.com/orsinium-labs/enum"
)

type MapperIdentityProvider interface {
	GetType() string
	GetName() string
}

type MapperBasicOperator interface {
	GetUrl(path string) (*PathUrlPair, error)
	ListUrls(pagination Pagination) (PathUrlPairList, error)
	PutUrl(pair *PathUrlPair) (*PathUrlPair, error)
	DeleteUrl(path string) error
}

type MapperExtendedOperator interface {
	SearchUrls(query string, mode SearchMode, pagination Pagination) (PathUrlPairList, error)
}

type MapperConfigurer interface {
	MapperIdentityProvider
	GetMapper() (Mapper, error)
	Singleton() bool
}

type Mapper interface {
	MapperIdentityProvider
	MapperBasicOperator
	// MapperExtendedOperator

	Readonly() bool
	Teardown() error
}

type Pagination struct {
	Offset int
	Limit  int
}

type SearchMode enum.Member[string]

var (
	SearchMode_Include = SearchMode{"include"}
	SearchMode_Fuzzy   = SearchMode{"fuzzy"}
)
