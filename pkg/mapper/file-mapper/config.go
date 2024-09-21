package file_mapper

import (
	"github.com/reimirno/golinks/pkg/mapper"
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
