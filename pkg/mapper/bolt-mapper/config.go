package bolt_mapper

import (
	"time"

	"github.com/boltdb/bolt"

	"github.com/reimirno/golinks/pkg/types"
)

const (
	BoltMapperConfigType = "BOLT"
	urlMapBucketName     = "urlMap"
)

var _ types.MapperConfigurer = (*BoltMapperConfig)(nil)

type BoltMapperConfig struct {
	Name    string `mapstructure:"name"`
	Path    string `mapstructure:"path"`
	Timeout int    `mapstructure:"timeout"`
}

func (b *BoltMapperConfig) GetType() string {
	return BoltMapperConfigType
}

func (b *BoltMapperConfig) GetName() string {
	return b.Name
}

func (b *BoltMapperConfig) Singleton() bool {
	return true
}

func (b *BoltMapperConfig) GetMapper() (types.Mapper, error) {
	var err error
	var db *bolt.DB
	db, err = bolt.Open(b.Path, 0o600, &bolt.Options{Timeout: time.Duration(b.Timeout) * time.Second})
	if err != nil {
		return nil, err
	}
	mapper := BoltMapper{
		name: b.Name,
		db:   db,
	}
	err = mapper.initializeBucket(urlMapBucketName)
	if err != nil {
		return nil, err
	}
	return &mapper, nil
}
