package file_mapper

import (
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
)

const (
	FileMapperConfigType = "FILE"
)

var _ mapper.MapperConfigurer = (*FileMapperConfig)(nil)

type FileMapperConfig struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

func (f *FileMapperConfig) GetName() string {
	return f.Name
}

func (f *FileMapperConfig) GetType() string {
	return FileMapperConfigType
}

func (f *FileMapperConfig) GetMapper() (mapper.Mapper, error) {
	pairs, err := parseFile(f.Path)
	if err != nil {
		return nil, err
	}
	mm := &FileMapper{
		name:  f.Name,
		pairs: pairs.ToMap(),
	}
	for _, pair := range mm.pairs {
		mapper.Sanitize(mm, pair)
	}
	return mm, nil
}

func (f *FileMapperConfig) Singleton() bool {
	return false
}

var _ mapper.Mapper = (*FileMapper)(nil)

type FileMapper struct {
	name  string
	pairs types.PathUrlPairMap
}

func (f *FileMapper) GetType() string {
	return FileMapperConfigType
}

func (f *FileMapper) GetName() string {
	return f.name
}

func (f *FileMapper) Teardown() error {
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
