package bolt_mapper

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
)

var _ mapper.Mapper = (*BoltMapper)(nil)

type BoltMapper struct {
	name string
	db   *bolt.DB
}

func (b *BoltMapper) Teardown() error {
	return b.db.Close() // removes file lock
}

func (b *BoltMapper) GetName() string {
	return b.name
}

func (b *BoltMapper) GetType() string {
	return BoltMapperConfigType
}

func (b *BoltMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	bytes, err := b.get(urlMapBucketName, path)
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return nil, nil
	}
	var pair types.PathUrlPair
	err = json.Unmarshal(bytes, &pair)
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (b *BoltMapper) ListUrls() (types.PathUrlPairList, error) {
	var pairs types.PathUrlPairList
	err := b.foreach(urlMapBucketName, func(key string, value []byte) error {
		var pair types.PathUrlPair
		err := json.Unmarshal(value, &pair)
		if err != nil {
			return err
		}
		pairs = append(pairs, &pair)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

func (b *BoltMapper) DeleteUrl(path string) error {
	return b.delete(urlMapBucketName, path)
}

func (b *BoltMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	mapper.Sanitize(b, pair) // always sanitize before saving, so no need to sanitize in get
	bytes, err := json.Marshal(pair)
	if err != nil {
		return nil, err
	}
	return pair, b.put(urlMapBucketName, pair.Path, bytes)
}

func (b *BoltMapper) Readonly() bool {
	return false
}
