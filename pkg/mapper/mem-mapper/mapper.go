package mem_mapper

import (
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
)

var _ mapper.Mapper = (*MemMapper)(nil)

type MemMapper struct {
	name  string
	pairs types.PathUrlPairMap
}

func (m *MemMapper) GetName() string {
	return m.name
}

func (m *MemMapper) GetType() string {
	return MemMapperConfigType
}

func (m *MemMapper) Teardown() error {
	return nil
}

func (f *MemMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	pair, ok := f.pairs[path]
	if !ok {
		return nil, nil
	}
	return pair, nil
}

func (f *MemMapper) ListUrls() (types.PathUrlPairList, error) {
	return f.pairs.ToList(), nil
}

func (f *MemMapper) DeleteUrl(path string) error {
	return mapper.ErrOperationNotSupported("delete")
}

func (f *MemMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	return nil, mapper.ErrOperationNotSupported("put")
}

func (f *MemMapper) Readonly() bool {
	return true
}
