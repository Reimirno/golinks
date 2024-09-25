package file_mapper

import (
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
)

var _ mapper.Mapper = (*FileMapper)(nil)

type FileMapper struct {
	name  string
	pairs types.PathUrlPairMap
	stop  func()
}

func (f *FileMapper) GetType() string {
	return FileMapperConfigType
}

func (f *FileMapper) GetName() string {
	return f.name
}

func (f *FileMapper) Teardown() error {
	if f.stop != nil {
		f.stop()
	}
	return nil
}

func (f *FileMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	pair, ok := f.pairs[path]
	if !ok {
		return nil, nil
	}
	return pair, nil
}

func (f *FileMapper) ListUrls() (types.PathUrlPairList, error) {
	return f.pairs.ToList(), nil
}

func (f *FileMapper) DeleteUrl(path string) error {
	return mapper.ErrOperationNotSupported("delete")
}

func (f *FileMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	return nil, mapper.ErrOperationNotSupported("put")
}

func (f *FileMapper) Readonly() bool {
	return true
}
