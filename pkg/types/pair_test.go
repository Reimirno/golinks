package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathUrlPair_String(t *testing.T) {
	tests := []struct {
		name string
		pair PathUrlPair
		want string
	}{
		{
			name: "basic pair",
			pair: PathUrlPair{Path: "/test", Url: "https://example.com"},
			want: "'/test' -> 'https://example.com'",
		},
		{
			name: "empty pair",
			pair: PathUrlPair{},
			want: "'' -> ''",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pair.String())
			assert.Equal(t, fmt.Sprintf("%s!", tt.want), fmt.Sprintf("%s!", tt.pair))
		})
	}
}

func TestPathUrlPair_Clone(t *testing.T) {
	tests := []struct {
		name     string
		original *PathUrlPair
	}{
		{
			name:     "full pair",
			original: &PathUrlPair{Path: "/test", Url: "https://example.com", Mapper: "testMapper", UseCount: 5},
		},
		{
			name:     "empty pair",
			original: &PathUrlPair{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.original.Clone()
			assert.NotSame(t, tt.original, clone, "Clone should return a new object")
			assert.True(t, tt.original.Equals(clone), "Clone should be equal to original")
		})
	}
}

func TestPathUrlPairMap_Clone(t *testing.T) {
	tests := []struct {
		name     string
		original *PathUrlPairMap
	}{
		{
			name: "full map",
			original: &PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				"/test2": &PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
		},
		{
			name:     "empty map",
			original: &PathUrlPairMap{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.original.Clone()
			assert.NotSame(t, tt.original, clone, "Clone should return a new object")
			assert.True(t, tt.original.Equals(clone), "Clone should be equal to original")
			// each pair in the map should be equal but not same pointer too
			for path, pair := range *tt.original {
				assert.Contains(t, *clone, path, "Clone should contain all original paths")
				assert.True(t, pair.Equals((*clone)[path]), "Clone values should match original values")
				assert.NotSame(t, pair, (*clone)[path], "Clone values should not be the same pointer")
			}
		})
	}
}

func TestPathUrlPairList_Clone(t *testing.T) {
	tests := []struct {
		name     string
		original *PathUrlPairList
	}{
		{
			name: "full list",
			original: &PathUrlPairList{
				&PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				&PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
		},
		{
			name:     "empty list",
			original: &PathUrlPairList{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.original.Clone()
			assert.NotSame(t, tt.original, clone, "Clone should return a new object")
			assert.True(t, tt.original.Equals(clone), "Clone should be equal to original")
			// each pair in the list should be equal but not same pointer too
			for idx := range *tt.original {
				assert.True(t, (*tt.original)[idx].Equals((*clone)[idx]), "Clone values should match original values")
				assert.NotSame(t, (*tt.original)[idx], (*clone)[idx], "Clone values should not be the same pointer")
			}
		})
	}
}

func TestPathUrlPairMap_ToList(t *testing.T) {
	tests := []struct {
		name string
		m    PathUrlPairMap
	}{
		{
			name: "non-empty map",
			m: PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				"/test2": &PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
		},
		{
			name: "empty map",
			m:    PathUrlPairMap{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := tt.m.ToList()
			assert.Len(t, list, len(tt.m), "List length should match map length")
			for _, pair := range list {
				assert.Contains(t, tt.m, pair.Path, "List should contain all map keys")
				assert.True(t, tt.m[pair.Path].Equals(pair), "List values should match map values")
			}
		})
	}
}

func TestPathUrlPairList_ToMap(t *testing.T) {
	tests := []struct {
		name string
		list PathUrlPairList
	}{
		{
			name: "non-empty list",
			list: PathUrlPairList{
				&PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				&PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
		},
		{
			name: "empty list",
			list: PathUrlPairList{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.list.ToMap()
			assert.Len(t, m, len(tt.list), "Map length should match list length")
			for _, pair := range tt.list {
				assert.Contains(t, m, pair.Path, "Map should contain all list paths")
				assert.True(t, pair.Equals(m[pair.Path]), "Map values should match list values")
			}
		})
	}
}

func TestPathUrlPair_Equals(t *testing.T) {
	tests := []struct {
		name string
		p1   *PathUrlPair
		p2   *PathUrlPair
		want bool
	}{
		{
			name: "equal pairs",
			p1:   &PathUrlPair{Path: "/test", Url: "https://example.com"},
			p2:   &PathUrlPair{Path: "/test", Url: "https://example.com"},
			want: true,
		},
		{
			name: "equal pairs with irrelevant fields",
			p1:   &PathUrlPair{Path: "/test", Url: "https://example.com", Mapper: "testMapper", UseCount: 5},
			p2:   &PathUrlPair{Path: "/test", Url: "https://example.com", Mapper: "haha", UseCount: 9},
			want: true,
		},
		{
			name: "different pairs",
			p1:   &PathUrlPair{Path: "/test", Url: "https://example.com"},
			p2:   &PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			want: false,
		},
		{
			name: "empty pairs",
			p1:   &PathUrlPair{},
			p2:   &PathUrlPair{},
			want: true,
		},
		{
			name: "nil pairs",
			p1:   nil,
			p2:   nil,
			want: true,
		},
		{
			name: "one nil pair",
			p1:   nil,
			p2:   &PathUrlPair{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p1.Equals(tt.p2))
		})
	}
}

func TestPathUrlPairMap_Equals(t *testing.T) {
	tests := []struct {
		name string
		m1   *PathUrlPairMap
		m2   *PathUrlPairMap
		want bool
	}{
		{
			name: "equal maps",
			m1: &PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				"/test2": &PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
			m2: &PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				"/test2": &PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
			want: true,
		},
		{
			name: "equal maps with irrelevant fields",
			m1: &PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com", Mapper: "testMapper", UseCount: 5},
				"/test2": &PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
			m2: &PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				"/test2": &PathUrlPair{Path: "/test2", Url: "https://example2.com", Mapper: "haha", UseCount: 9},
			},
			want: true,
		},
		{
			name: "different maps",
			m1: &PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				"/test2": &PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
			m2: &PathUrlPairMap{
				"/test1": &PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				"/test3": &PathUrlPair{Path: "/test3", Url: "https://example3.com"},
			},
			want: false,
		},
		{
			name: "empty maps",
			m1:   &PathUrlPairMap{},
			m2:   &PathUrlPairMap{},
			want: true,
		},
		{
			name: "nil maps",
			m1:   nil,
			m2:   nil,
			want: true,
		},
		{
			name: "one nil map",
			m1:   nil,
			m2:   &PathUrlPairMap{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m1.Equals(tt.m2))
		})
	}
}

func TestPathUrlPairList_Equals(t *testing.T) {
	tests := []struct {
		name string
		l1   *PathUrlPairList
		l2   *PathUrlPairList
		want bool
	}{
		{
			name: "equal lists",
			l1: &PathUrlPairList{
				&PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				&PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
			l2: &PathUrlPairList{
				&PathUrlPair{Path: "/test2", Url: "https://example2.com"},
				&PathUrlPair{Path: "/test1", Url: "https://example1.com"},
			},
			want: true,
		},
		{
			name: "equal lists with irrelevant fields",
			l1: &PathUrlPairList{
				&PathUrlPair{Path: "/test1", Url: "https://example1.com", Mapper: "testMapper", UseCount: 5},
				&PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
			l2: &PathUrlPairList{
				&PathUrlPair{Path: "/test2", Url: "https://example2.com", Mapper: "haha", UseCount: 9},
			},
		},
		{
			name: "different lists",
			l1: &PathUrlPairList{
				&PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				&PathUrlPair{Path: "/test2", Url: "https://example2.com"},
			},
			l2: &PathUrlPairList{
				&PathUrlPair{Path: "/test1", Url: "https://example1.com"},
				&PathUrlPair{Path: "/test3", Url: "https://example3.com"},
			},
			want: false,
		},
		{
			name: "empty lists",
			l1:   &PathUrlPairList{},
			l2:   &PathUrlPairList{},
			want: true,
		},
		{
			name: "nil lists",
			l1:   nil,
			l2:   nil,
			want: true,
		},
		{
			name: "one nil list",
			l1:   nil,
			l2:   &PathUrlPairList{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l1.Equals(tt.l2))
		})
	}
}
