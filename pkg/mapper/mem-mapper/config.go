package mem_mapper

import (
	"github.com/reimirno/golinks/pkg/sanitizer"
	"github.com/reimirno/golinks/pkg/types"
)

const (
	MemMapperConfigType = "MEM"
)

var _ types.MapperConfigurer = (*MemMapperConfig)(nil)

type MemMapperConfig struct {
	Name  string              `mapstructure:"name"`
	Pairs []types.PathUrlPair `mapstructure:"pairs"`
}

func (m *MemMapperConfig) GetName() string {
	return m.Name
}

func (m *MemMapperConfig) GetType() string {
	return MemMapperConfigType
}

func (m *MemMapperConfig) GetMapper() (types.Mapper, error) {
	pairs := make(types.PathUrlPairMap)
	for _, pair := range m.Pairs {
		pairs[pair.Path] = &pair
	}
	mm := &MemMapper{
		name:  m.Name,
		pairs: pairs,
	}
	err := sanitizer.SanitizeInputMap(mm, &mm.pairs)
	if err != nil {
		return nil, err
	}
	return mm, nil
}

func (m *MemMapperConfig) Singleton() bool {
	return true
}
