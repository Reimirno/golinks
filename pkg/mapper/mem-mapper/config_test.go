package mem_mapper

import (
	"testing"

	"github.com/reimirno/golinks/pkg/types"
	"github.com/stretchr/testify/assert"
)

var fakePair = &types.PathUrlPair{
	Path: "fk",
	Url:  "https://fake.com",
}

var (
	fakeConfig = MemMapperConfig{
		Name: "fake",
		Pairs: []types.PathUrlPair{
			*fakePair,
		},
	}
	fakeConfigEmptyName = MemMapperConfig{
		Name: "",
		Pairs: []types.PathUrlPair{
			*fakePair,
		},
	}
	fakeConfigEmptyPairs = MemMapperConfig{
		Name:  "fake",
		Pairs: []types.PathUrlPair{},
	}
)

func TestMemMapperConfig_GetName(t *testing.T) {
	tests := []struct {
		name         string
		mapperConfig *MemMapperConfig
		want         string
	}{
		{name: "test", mapperConfig: &fakeConfig, want: fakeConfig.Name},
		{name: "empty name", mapperConfig: &fakeConfigEmptyName, want: fakeConfigEmptyName.Name},
		{name: "empty pairs", mapperConfig: &fakeConfigEmptyPairs, want: fakeConfigEmptyPairs.Name},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mapperConfig.GetName())
		})
	}
}

func TestMemMapperConfig_GetType(t *testing.T) {
	tests := []struct {
		name         string
		mapperConfig *MemMapperConfig
		want         string
	}{
		{name: "test", mapperConfig: &fakeConfig, want: MemMapperConfigType},
		{name: "empty name", mapperConfig: &fakeConfigEmptyName, want: MemMapperConfigType},
		{name: "empty pairs", mapperConfig: &fakeConfigEmptyPairs, want: MemMapperConfigType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mapperConfig.GetType())
		})
	}
}

func TestMemMapperConfig_Singleton(t *testing.T) {
	tests := []struct {
		name         string
		mapperConfig *MemMapperConfig
		want         bool
	}{
		{name: "test", mapperConfig: &fakeConfig, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mapperConfig.Singleton())
		})
	}
}

func TestMemMapperConfig_GetMapper(t *testing.T) {
	tests := []struct {
		name         string
		mapperConfig *MemMapperConfig
		want         *MemMapper
	}{
		{
			name:         "test",
			mapperConfig: &fakeConfig,
			want: &MemMapper{
				name: fakeConfig.Name,
				pairs: types.PathUrlPairMap{
					fakePair.Path: fakePair,
				},
			},
		},
		{
			name:         "empty pairs",
			mapperConfig: &fakeConfigEmptyPairs,
			want: &MemMapper{
				name:  fakeConfigEmptyPairs.Name,
				pairs: types.PathUrlPairMap{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mapperConfig.GetMapper()
			assert.NoError(t, err)
			memMapper, ok := got.(*MemMapper)
			assert.True(t, ok, "Expected *MemMapper, got %T", got)
			assert.Equal(t, tt.want.name, memMapper.name)
			assert.True(t, tt.want.pairs.Equals(&memMapper.pairs), "Expected %v, got %v", tt.want.pairs, memMapper.pairs)
			assert.NoError(t, memMapper.Teardown())
		})
	}
}
