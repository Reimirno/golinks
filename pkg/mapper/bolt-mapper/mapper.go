package bolt_mapper

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/reimirno/golinks/pkg/types"
)

var _ types.Mapper = (*BoltMapper)(nil)

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

func (b *BoltMapper) ListUrls(pagination types.Pagination) (types.PathUrlPairList, error) {
	var pairs types.PathUrlPairList
	curIdx := 0
	err := b.forsome(urlMapBucketName, func(key string, value []byte) error {
		if curIdx < pagination.Offset {
			curIdx++
			return nil
		}
		var pair types.PathUrlPair
		err := json.Unmarshal(value, &pair)
		if err != nil {
			return err
		}
		pairs = append(pairs, &pair)
		curIdx++
		return nil
	}, pagination.Offset+pagination.Limit)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

func (b *BoltMapper) DeleteUrl(path string) error {
	return b.delete(urlMapBucketName, path)
}

func (b *BoltMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	bytes, err := json.Marshal(pair)
	if err != nil {
		return nil, err
	}
	return pair, b.put(urlMapBucketName, pair.Path, bytes)
}

func (b *BoltMapper) Readonly() bool {
	return false
}
