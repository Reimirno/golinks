package file_mapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fakeConfigEmptyPath = FileMapperConfig{
		Name:         "fake",
		Path:         "",
		SyncInterval: -1,
	}
)

func TestFileMapperConfig_GetName(t *testing.T) {
	tests := []struct {
		name         string
		mapperConfig *FileMapperConfig
		want         string
	}{
		{name: "test", mapperConfig: &fakeConfigEmptyPath, want: fakeConfigEmptyPath.Name},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mapperConfig.GetName())
		})
	}
}

func TestFileMapperConfig_GetType(t *testing.T) {
	tests := []struct {
		name         string
		mapperConfig *FileMapperConfig
		want         string
	}{
		{name: "test", mapperConfig: &fakeConfigEmptyPath, want: FileMapperConfigType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mapperConfig.GetType())
		})
	}
}

func TestFileMapperConfig_Singleton(t *testing.T) {
	tests := []struct {
		name         string
		mapperConfig *FileMapperConfig
		want         bool
	}{
		{name: "test", mapperConfig: &fakeConfigEmptyPath, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.mapperConfig.Singleton())
		})
	}
}

func TestFileMapperConfig_GetMapper(t *testing.T) {
	tests := []struct {
		name           string
		tempFileConfig *tempFileConfig
		want           *FileMapper
		expectedError  bool
	}{
		{
			name:           "test config with yaml file",
			tempFileConfig: yamlFileConfig,
			want: &FileMapper{
				name:  yamlFileConfig.name,
				pairs: pairList.ToMap(),
			},
		},
		{
			name:           "test config with json file",
			tempFileConfig: jsonFileConfig,
			want: &FileMapper{
				name:  jsonFileConfig.name,
				pairs: pairList.ToMap(),
			},
		},
		{
			name:           "test config with malformed yaml file",
			tempFileConfig: malformedYamlFileConfig,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := createTempFile(*tt.tempFileConfig)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())
			mapperConfig := FileMapperConfig{
				Name: tt.tempFileConfig.name,
				Path: tmpfile.Name(),
			}

			got, err := mapperConfig.GetMapper()

			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fileMapper, ok := got.(*FileMapper)
			assert.True(t, ok, "Expected *FileMapper, got %T", got)
			assert.Equal(t, tt.want.name, fileMapper.name)
			assert.True(t, tt.want.pairs.Equals(&fileMapper.pairs), "Expected %v, got %v", tt.want.pairs, fileMapper.pairs)
		})
	}
}
