package mem_mapper

import (
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
)

var _ types.Mapper = (*MemMapper)(nil)

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

func (m *MemMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	if pair, ok := m.pairs[path]; ok {
		return pair, nil
	}
	return nil, nil
}

func (m *MemMapper) ListUrls() (types.PathUrlPairList, error) {
	return m.pairs.ToList(), nil
}

func (m *MemMapper) DeleteUrl(path string) error {
	return mapper.ErrOperationNotSupported("delete")
}

func (m *MemMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	return nil, mapper.ErrOperationNotSupported("put")
}

func (m *MemMapper) Readonly() bool {
	return true
}
