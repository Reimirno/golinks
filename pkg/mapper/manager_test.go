package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reimirno/golinks/pkg/sanitizer"
	"github.com/reimirno/golinks/pkg/types"
	"github.com/reimirno/golinks/pkg/utils"
)

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
		Name:        "mock",
		IsSingleton: false,
		IsReadOnly:  false,
		StarterPairs: types.PathUrlPairMap{
			"fk":  fakePair,
			"fk2": fakePair2,
		},
	}
	mockConfigurerAlt = &MockMapperConfigurer{
		Name:        "mockAlt",
		IsSingleton: false,
		IsReadOnly:  false,
		StarterPairs: types.PathUrlPairMap{
			"fk": fakePairAlt,
		},
	}
	mockConfigurer2 = &MockMapperConfigurer{
		Name:        "mock2",
		IsSingleton: false,
		IsReadOnly:  false,
		StarterPairs: types.PathUrlPairMap{
			"fk3": fakePair3,
		},
	}
	mockConfigurerReadonly = &MockMapperConfigurer{
		Name:        "mockReadonly",
		IsSingleton: false,
		IsReadOnly:  true,
	}
	mockConfigurerSingleton = &MockMapperConfigurer{
		Name:        "mockSingleton",
		IsSingleton: true,
		IsReadOnly:  false,
	}
)

func TestNewMapperManager(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []types.MapperConfigurer
		persistorName string
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []types.MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			wantErr:       false,
		},
		{
			name:          "invalid persistor should fail",
			configurers:   []types.MapperConfigurer{mockConfigurer},
			persistorName: "invalid",
			wantErr:       true,
		},
		{
			name:          "nil persistor should be fine",
			configurers:   []types.MapperConfigurer{mockConfigurer},
			persistorName: "",
			wantErr:       false,
		},
		{
			name:          "duplicate non-singleton should pass",
			configurers:   []types.MapperConfigurer{mockConfigurer, mockConfigurerReadonly},
			persistorName: mockConfigurer.Name,
			wantErr:       false,
		},
		{
			name:          "duplicate name should fail",
			configurers:   []types.MapperConfigurer{mockConfigurer, mockConfigurer},
			persistorName: mockConfigurer.Name,
			wantErr:       true,
		},
		{
			name:          "duplicate singleton should fail",
			configurers:   []types.MapperConfigurer{mockConfigurer, mockConfigurerSingleton},
			persistorName: mockConfigurer.Name,
			wantErr:       true,
		},
		{
			name:          "no mappers should fail",
			configurers:   []types.MapperConfigurer{},
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
		configurers []types.MapperConfigurer
		numUrls     int
	}{
		{
			name:        "happy path",
			configurers: []types.MapperConfigurer{mockConfigurer},
			numUrls:     len(mockConfigurer.StarterPairs),
		},
		{
			name:        "multiple mappers",
			configurers: []types.MapperConfigurer{mockConfigurer, mockConfigurer2},
			numUrls:     len(mockConfigurer.StarterPairs) + len(mockConfigurer2.StarterPairs),
		},
		{
			name:        "multiple mappers with overlap",
			configurers: []types.MapperConfigurer{mockConfigurer, mockConfigurerAlt},
			numUrls:     len(mockConfigurer.StarterPairs),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.configurers[0].GetName(), test.configurers)
			assert.NoError(t, err)
			urls, err := mm.ListUrls(utils.DefaultPagination)
			assert.NoError(t, err)
			assert.Equal(t, test.numUrls, len(urls))
		})
	}
}

func TestMapperManager_GetUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []types.MapperConfigurer
		persistorName string
		path          string
		url           *types.PathUrlPair
		wantErr       bool
	}{
		{
			name:          "happy path",
			configurers:   []types.MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "/fk",
			url:           fakePair,
		},
		{
			name:          "path not found",
			configurers:   []types.MapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "invalid",
			url:           nil,
			wantErr:       false,
		},
		{
			name:          "order of configurers is important",
			configurers:   []types.MapperConfigurer{mockConfigurerAlt, mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "/fk",
			url:           fakePairAlt,
		},
		{
			name:          "readonly mapper works",
			configurers:   []types.MapperConfigurer{mockConfigurerReadonly},
			persistorName: "",
			path:          "/fk",
			url:           nil,
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
					err = sanitizer.SanitizeInput(mm.getPersistor(), test.url)
					assert.NoError(t, err)
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
		configurers   []*MockMapperConfigurer
		persistorName string
		pair          *types.PathUrlPair
		wantErr       bool
		finalPair     *types.PathUrlPair // final pair you can GET from update
	}{
		{
			name:          "happy path update",
			configurers:   []*MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			pair:          fakePairAlt,
			finalPair:     fakePairAlt,
		},
		{
			name:          "happy path insert",
			configurers:   []*MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			pair:          fakePair3,
			finalPair:     fakePair3,
		},
		{
			name:          "readonly mapper should fail",
			configurers:   []*MockMapperConfigurer{mockConfigurerReadonly},
			persistorName: "",
			pair:          fakePair3,
			wantErr:       true,
		},
		{
			name:          "update a key repeated at two mappers only update from first",
			configurers:   []*MockMapperConfigurer{mockConfigurerAlt, mockConfigurer},
			persistorName: mockConfigurer.Name,
			pair:          &types.PathUrlPair{Path: "fk", Url: "https://new.com"},
			finalPair:     &types.PathUrlPair{Path: "fk", Url: "https://new.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.persistorName, CloneConfigurers(test.configurers))
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
				// sanitize test.finalPair before comparing
				err = sanitizer.SanitizeInput(mm.getPersistor(), test.finalPair)
				assert.NoError(t, err)
				assert.True(t, test.finalPair.Equals(pair), "Expected %v, got %v", test.finalPair, pair)
			}
		})
	}
}

func TestMapperManager_DeleteUrl(t *testing.T) {
	tests := []struct {
		name          string
		configurers   []*MockMapperConfigurer
		persistorName string
		path          string
		wantErr       bool
		finalPair     *types.PathUrlPair // final pair you can GET from update
	}{
		{
			name:          "happy path",
			configurers:   []*MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "fk",
			wantErr:       false,
			finalPair:     nil,
		},
		{
			name:          "path not found is fine",
			configurers:   []*MockMapperConfigurer{mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "invalid",
			wantErr:       false,
			finalPair:     nil,
		},
		{
			name:          "readonly mapper should fail",
			configurers:   []*MockMapperConfigurer{mockConfigurerReadonly},
			persistorName: "",
			path:          "fk",
			wantErr:       true,
		},
		{
			name:          "delete a key repeated at two mappers only deletes from first",
			configurers:   []*MockMapperConfigurer{mockConfigurerAlt, mockConfigurer},
			persistorName: mockConfigurer.Name,
			path:          "fk",
			finalPair:     fakePair,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm, err := NewMapperManager(test.persistorName, CloneConfigurers(test.configurers))
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
		configurers []types.MapperConfigurer
	}{
		{
			name:        "happy path",
			configurers: []types.MapperConfigurer{mockConfigurer},
		},
		{
			name:        "multiple mappers",
			configurers: []types.MapperConfigurer{mockConfigurer, mockConfigurer2},
		},
		{
			name:        "multiple mappers with overlap",
			configurers: []types.MapperConfigurer{mockConfigurer, mockConfigurerAlt},
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
