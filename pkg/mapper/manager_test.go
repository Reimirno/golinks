package mapper

import (
	"testing"

	"github.com/reimirno/golinks/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMapper is a mock implementation of the Mapper interface for testing purposes.
// It uses a in-memory map to store PathUrlPair objects.
// Supports all operations, but does persist anything.
type MockMapper struct {
	mock.Mock
	pairs    types.PathUrlPairMap
	readonly bool
	name     string
}

func (m *MockMapper) GetType() string {
	return "mock"
}

func (m *MockMapper) GetName() string {
	return m.name
}

func (m *MockMapper) Readonly() bool {
	return m.readonly
}

func (m *MockMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	if pair, ok := m.pairs[path]; ok {
		return pair, nil
	}
	return nil, nil
}

func (m *MockMapper) ListUrls() (types.PathUrlPairList, error) {
	return m.pairs.ToList(), nil
}

func (m *MockMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	if m.readonly {
		return nil, ErrOperationNotSupported("put")
	}
	Sanitize(m, pair)
	m.pairs[pair.Path] = pair
	return pair, nil
}

func (m *MockMapper) DeleteUrl(path string) error {
	if m.readonly {
		return ErrOperationNotSupported("delete")
	}
	delete(m.pairs, path)
	return nil
}

func (m *MockMapper) Teardown() error {
	return nil
}

var _ Mapper = (*MockMapper)(nil)

// MockMapperConfigurer is a mock implementation of the MapperConfigurer interface for testing purposes.
// It just returns a new MockMapper instance.
type MockMapperConfigurer struct {
	mock.Mock
	name         string
	singleton    bool
	readonly     bool
	starterPairs types.PathUrlPairMap
}

func (m *MockMapperConfigurer) GetType() string {
	return "mock"
}

func (m *MockMapperConfigurer) GetName() string {
	return m.name
}

func (m *MockMapperConfigurer) GetMapper() (Mapper, error) {
	mapper := new(MockMapper)
	mapper.pairs = m.starterPairs
	for _, pair := range mapper.pairs {
		pair.Mapper = m.name
	}
	mapper.readonly = m.readonly
	mapper.name = m.name
	if mapper.pairs == nil {
		mapper.pairs = make(types.PathUrlPairMap)
	}
	return mapper, nil
}

func (m *MockMapperConfigurer) Singleton() bool {
	return m.singleton
}

var _ MapperConfigurer = (*MockMapperConfigurer)(nil)

var (
	fakePair = &types.PathUrlPair{
		Path: "fk",
		Url:  "https://fake.com",
	}
	fakePairAlt = &types.PathUrlPair{
		Path: "fk",
		Url:  "https://fakealt.com",
	}
	fakePair2 = &types.PathUrlPair{
		Path: "fk2",
		Url:  "https://fake2.com",
	}
	fakePair3 = &types.PathUrlPair{
		Path: "fk3",
		Url:  "https://fake3.com",
	}

	mockConfigurer = &MockMapperConfigurer{
		name:      "mock",
		singleton: false,
		readonly:  false,
		starterPairs: types.PathUrlPairMap{
			"fk":  fakePair,
			"fk2": fakePair2,
		},
	}
	mockConfigurerAlt = &MockMapperConfigurer{
		name:      "mockAlt",
		singleton: false,
		readonly:  false,
		starterPairs: types.PathUrlPairMap{
			"fk": fakePairAlt,
		},
	}
	mockConfigurer2 = &MockMapperConfigurer{
		name:      "mock2",
		singleton: false,
		readonly:  false,
		starterPairs: types.PathUrlPairMap{
			"fk3": fakePair3,
		},
	}
	mockConfigurerReadonly = &MockMapperConfigurer{
		name:      "mockReadonly",
		singleton: false,
		readonly:  true,
	}
	mockConfigurerSingleton = &MockMapperConfigurer{
		name:      "mockSingleton",
		singleton: true,
		readonly:  false,
	}
)

func TestNewMapperManager(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []MapperConfigurer
		persistorName string
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.name,
			wantErr:       false,
		},
		{
			name:          "invalid persistor should fail",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: "invalid",
			wantErr:       true,
		},
		{
			name:          "nil persistor should be fine",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: "",
			wantErr:       false,
		},
		{
			name:          "duplicate non-singleton should pass",
			configurers:   []MapperConfigurer{mockConfigurer, mockConfigurerReadonly},
			persistorName: mockConfigurer.name,
			wantErr:       false,
		},
		{
			name:          "duplicate name should fail",
			configurers:   []MapperConfigurer{mockConfigurer, mockConfigurer},
			persistorName: mockConfigurer.name,
			wantErr:       true,
		},
		{
			name:          "duplicate singleton should fail",
			configurers:   []MapperConfigurer{mockConfigurer, mockConfigurerSingleton},
			persistorName: mockConfigurer.name,
			wantErr:       true,
		},
		{
			name:          "no mappers should fail",
			configurers:   []MapperConfigurer{},
			persistorName: "",
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.persistorName, test.configurers)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, mm)
			} else {
				assert.NoError(t, err)
				if test.persistorName != "" {
					assert.Equal(t, test.persistorName, mm.getPersistor().GetName())
				} else {
					assert.Nil(t, mm.getPersistor())
				}
				assert.Equal(t, len(test.configurers), len(mm.mappers))
			}
		})
	}
}

func TestMapperManager_ListUrls(t *testing.T) {
	tests := []struct {
		name        string
		configurers []MapperConfigurer
		numUrls     int
	}{
		{
			name:        "happy path",
			configurers: []MapperConfigurer{mockConfigurer},
			numUrls:     len(mockConfigurer.starterPairs),
		},
		{
			name:        "multiple mappers",
			configurers: []MapperConfigurer{mockConfigurer, mockConfigurer2},
			numUrls:     len(mockConfigurer.starterPairs) + len(mockConfigurer2.starterPairs),
		},
		{
			name:        "multiple mappers with overlap",
			configurers: []MapperConfigurer{mockConfigurer, mockConfigurerAlt},
			numUrls:     len(mockConfigurer.starterPairs),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.configurers[0].GetName(), test.configurers)
			assert.NoError(t, err)
			urls, err := mm.ListUrls()
			assert.NoError(t, err)
			assert.Equal(t, test.numUrls, len(urls))
		})
	}
}

func TestMapperManager_GetUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []MapperConfigurer
		persistorName string
		path          string
		url           *types.PathUrlPair
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.name,
			path:          "fk",
			url:           fakePair,
		},
		{
			name:          "path not found",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.name,
			path:          "invalid",
			url:           nil,
			wantErr:       false,
		},
		{
			name:          "order of configurers is important",
			configurers:   []MapperConfigurer{mockConfigurerAlt, mockConfigurer},
			persistorName: mockConfigurer.name,
			path:          "fk",
			url:           fakePairAlt,
		},
		{
			name:          "readonly mapper works",
			configurers:   []MapperConfigurer{mockConfigurerReadonly},
			persistorName: "",
			path:          "fk",
			url:           fakePair,
			wantErr:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.persistorName, test.configurers)
			assert.NoError(t, err)
			pair, err := mm.GetUrl(test.path, false)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pair)
			} else {
				assert.NoError(t, err)
				if test.url != nil {
					assert.True(t, test.url.Equals(pair), "Expected %v, got %v", test.url, pair)
				} else {
					assert.Nil(t, pair)
				}
			}
		})
	}
}

func TestMapperManager_PutUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []MapperConfigurer
		persistorName string
		pair          *types.PathUrlPair
		wantErr       bool
		finalPair     *types.PathUrlPair // final pair you can GET from update
	}{
		{
			name:          "happy path update",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.name,
			pair:          fakePairAlt,
			finalPair:     fakePairAlt,
		},
		{
			name:          "happy path insert",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.name,
			pair:          fakePair3,
			finalPair:     fakePair3,
		},
		{
			name:          "readonly mapper should fail",
			configurers:   []MapperConfigurer{mockConfigurerReadonly},
			persistorName: "",
			pair:          fakePair3,
			wantErr:       true,
		},
		{
			name:          "update still respects precedence",
			configurers:   []MapperConfigurer{mockConfigurerAlt, mockConfigurer},
			persistorName: mockConfigurer.name,
			pair:          &types.PathUrlPair{Path: "fk", Url: "https://new.com"},
			finalPair:     fakePairAlt, // GET still retrieve from mockConfigurerAlt, as it takes precedence
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.persistorName, test.configurers)
			assert.NoError(t, err)
			pair, err := mm.PutUrl(test.pair)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pair)
			} else {
				assert.NoError(t, err)
				assert.True(t, test.pair.Equals(pair), "Expected %v, got %v", test.pair, pair)
				pair, err := mm.GetUrl(test.pair.Path, false)
				assert.NoError(t, err)
				assert.True(t, test.finalPair.Equals(pair), "Expected %v, got %v", test.finalPair, pair)
			}
		})
	}
}

func TestMapperManager_DeleteUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []MapperConfigurer
		persistorName string
		path          string
		wantErr       bool
		finalPair     *types.PathUrlPair // final pair you can GET from update
	}{
		{
			name:          "happy path",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.name,
			path:          "fk",
			wantErr:       false,
			finalPair:     nil,
		},
		{
			name:          "path not found is fine",
			configurers:   []MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.name,
			path:          "invalid",
			wantErr:       false,
			finalPair:     nil,
		},
		{
			name:          "readonly mapper should fail",
			configurers:   []MapperConfigurer{mockConfigurerReadonly},
			persistorName: "",
			path:          "fk",
			wantErr:       true,
		},
		{
			name:          "delete still respects precedence",
			configurers:   []MapperConfigurer{mockConfigurerAlt, mockConfigurer},
			persistorName: mockConfigurer.name,
			path:          "fk",
			finalPair:     fakePairAlt, // GET still retrieve from mockConfigurerAlt, as it takes precedence
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.persistorName, test.configurers)
			assert.NoError(t, err)
			err = mm.DeleteUrl(test.path)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				pair, err := mm.GetUrl(test.path, false)
				assert.NoError(t, err)
				assert.True(t, test.finalPair.Equals(pair), "Expected %v, got %v", test.finalPair, pair)
			}
		})
	}
}

func TestMapperManager_Teardown(t *testing.T) {
	tests := []struct {
		name        string
		configurers []MapperConfigurer
	}{
		{
			name:        "happy path",
			configurers: []MapperConfigurer{mockConfigurer},
		},
		{
			name:        "multiple mappers",
			configurers: []MapperConfigurer{mockConfigurer, mockConfigurer2},
		},
		{
			name:        "multiple mappers with overlap",
			configurers: []MapperConfigurer{mockConfigurer, mockConfigurerAlt},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.configurers[0].GetName(), test.configurers)
			assert.NoError(t, err)
			err = mm.Teardown()
			assert.NoError(t, err)
		})
	}
}
