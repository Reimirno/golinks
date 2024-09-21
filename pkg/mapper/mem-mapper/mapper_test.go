package mem_mapper

import (
	"testing"

	"github.com/reimirno/golinks/pkg/types"
	"github.com/stretchr/testify/assert"
)

var (
	fakeMapper = &MemMapper{
		name: "fake",
		pairs: types.PathUrlPairMap{
			"fk": fakePair,
		},
	}
	fakeMapperEmptyName = &MemMapper{
		name: "",
		pairs: types.PathUrlPairMap{
			"fk": fakePair,
		},
	}
	fakeMapperEmptyPairs = &MemMapper{
		name:  "fakeEmpty",
		pairs: types.PathUrlPairMap{},
	}
)

func TestMemMapper_GetName(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		want   string
	}{
		{name: "happy path", mapper: fakeMapper, want: "fake"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mapper.GetName()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemMapper_GetType(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		want   string
	}{
		{name: "happy path", mapper: fakeMapper, want: MemMapperConfigType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mapper.GetType()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemMapper_Teardown(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		want   error
	}{
		{name: "happy path", mapper: fakeMapper, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mapper.Teardown()
			assert.NoError(t, err)
		})
	}
}

func TestMemMapper_GetUrl(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		path   string
		want   *types.PathUrlPair
	}{
		{name: "happy path", mapper: fakeMapper, path: "fk", want: fakePair},
		{name: "path not found", mapper: fakeMapper, path: "none", want: nil},
		{name: "empty name still happy", mapper: fakeMapperEmptyName, path: "fk", want: fakePair},
		{name: "empty pairs", mapper: fakeMapperEmptyPairs, path: "fk", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mapper.GetUrl(tt.path)
			assert.NoError(t, err)
			assert.True(t, tt.want.Equals(got), "Expected %v, got %v", tt.want, got)
		})
	}
}

func TestMemMapper_ListUrls(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		want   *types.PathUrlPairList
	}{
		{name: "happy path", mapper: fakeMapper, want: &types.PathUrlPairList{fakePair}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.mapper.ListUrls()
			assert.NoError(t, err)
			assert.True(t, tt.want.Equals(&got), "Expected %v, got %v", tt.want, got)
		})
	}
}

func TestMemMapper_PutUrl(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		pair   *types.PathUrlPair
	}{
		{name: "happy path", mapper: fakeMapper, pair: fakePair},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.mapper.PutUrl(tt.pair)
			assert.Error(t, err)
		})
	}
}

func TestMemMapper_DeleteUrl(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		path   string
	}{
		{name: "happy path", mapper: fakeMapper, path: "fk"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mapper.DeleteUrl(tt.path)
			assert.Error(t, err)
		})
	}
}

func TestMemMapper_Readonly(t *testing.T) {
	tests := []struct {
		name   string
		mapper *MemMapper
		want   bool
	}{
		{name: "happy path", mapper: fakeMapper, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mapper.Readonly()
			assert.Equal(t, tt.want, got)
		})
	}
}
