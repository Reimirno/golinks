package sanitizer

import (
	"testing"

	"github.com/reimirno/golinks/pkg/types"
	"github.com/stretchr/testify/assert"
)

type NameOnlyMapper struct {
	Name       string
	IsReadOnly bool
}

func (m *NameOnlyMapper) GetType() string {
	return "nameonly"
}

func (m *NameOnlyMapper) GetName() string {
	return m.Name
}

func (m *NameOnlyMapper) Readonly() bool {
	return m.IsReadOnly
}

func (m *NameOnlyMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	return nil, nil
}

func (m *NameOnlyMapper) ListUrls() (types.PathUrlPairList, error) {
	return nil, nil
}

func (m *NameOnlyMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	return nil, nil
}

func (m *NameOnlyMapper) DeleteUrl(path string) error {
	return nil
}

func (m *NameOnlyMapper) Teardown() error {
	return nil
}

var _ types.Mapper = (*NameOnlyMapper)(nil)

var (
	nameOnlyMapper         = &NameOnlyMapper{Name: "nameonly"}
	nameOnlyMapperReadOnly = &NameOnlyMapper{Name: "nameonlyreadonly", IsReadOnly: true}
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		mapper   types.Mapper
		input    *types.PathUrlPair
		expected *types.PathUrlPair
		wantErr  bool
	}{
		{
			"Noop",
			nameOnlyMapper,
			&types.PathUrlPair{Path: "/noop", Url: "https://noop"},
			&types.PathUrlPair{Path: "/noop", Url: "https://noop"},
			false,
		},
		{
			"Trim slashes",
			nameOnlyMapper,
			&types.PathUrlPair{Path: "/example/path/", Url: "https://example.com"},
			&types.PathUrlPair{Path: "/example/path", Url: "https://example.com"},
			false,
		},
		{
			"Zero out use count",
			nameOnlyMapper,
			&types.PathUrlPair{Path: "/example/path/", Url: "https://example.com", UseCount: 1},
			&types.PathUrlPair{Path: "/example/path", Url: "https://example.com"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeInput(tt.mapper, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.expected.Equals(tt.input), "Expected %v, got %v", tt.expected, tt.input)
				assert.Equal(t, tt.mapper.GetName(), tt.input.Mapper)
				assert.Equal(t, 0, tt.input.UseCount)
			}
		})
	}
}

func TestSanitizeOutput(t *testing.T) {
	tests := []struct {
		name     string
		mapper   types.Mapper
		input    *types.PathUrlPair
		expected *types.PathUrlPair
	}{
		{
			"Noop",
			nameOnlyMapper,
			&types.PathUrlPair{Path: "/noop", Url: "https://noop"},
			&types.PathUrlPair{Path: "/noop", Url: "https://noop"},
		},
		{
			"Zero out use count for readonly mapper",
			nameOnlyMapperReadOnly,
			&types.PathUrlPair{Path: "/noop", Url: "https://noop", UseCount: 1},
			&types.PathUrlPair{Path: "/noop", Url: "https://noop"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SanitizeOutput(tt.mapper, tt.input)
			assert.True(t, tt.expected.Equals(tt.input), "Expected %v, got %v", tt.expected, tt.input)
			assert.Equal(t, tt.expected.UseCount, tt.input.UseCount)
		})
	}
}

func TestCanonicalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"Noop", "/noop", "/noop", false},
		{"Trim slashes", "/example/path/", "/example/path", false},
		{"Remove special characters", "ex_am.ple-path", "/examplepath", false},
		{"Replace multiple slashes", "//example///path//", "/example/path", false},
		{"Add leading slash", "example/path", "/example/path", false},
		{"Reserved path /", "/", "", true},
		{"Reserved path /d", "/d", "", true},
		{"Reserved path /d/", "/d/example", "", true},
		{"Not a reserved path /dd/", "/dd/example", "/dd/example", false},
		{"Invalid characters escaped", "/example/path with spaces", "/example/path%20with%20spaces", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CanonicalizePath(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCanonicalizeUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"Noop", "https://noop", "https://noop", false},
		{"Trim spaces", " https://example.com ", "https://example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CanonicalizeUrl(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
