package types

type MapperIdentityProvider interface {
	GetType() string
	GetName() string
}

type MapperBasicOperator interface {
	GetUrl(path string) (*PathUrlPair, error)
	ListUrls() (PathUrlPairList, error)
	PutUrl(pair *PathUrlPair) (*PathUrlPair, error)
	DeleteUrl(path string) error
}

type MapperExtendedOperator interface {
	SearchUrls(query string, limit int) (PathUrlPairList, error)
}

type MapperConfigurer interface {
	MapperIdentityProvider
	GetMapper() (Mapper, error)
	Singleton() bool
}

type Mapper interface {
	MapperIdentityProvider
	MapperBasicOperator

	Readonly() bool
	Teardown() error
}
