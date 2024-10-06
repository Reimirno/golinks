package file_mapper

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/reimirno/golinks/pkg/sanitizer"
)

var fakeConfigEmptyPath = FileMapperConfig{
	Name:         "fake",
	Path:         "",
	SyncInterval: -1,
}

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
		name                        string
		tempFileConfig              *tempFileConfig
		want                        *FileMapper
		syncInterval                int
		tempFileConfigNextWrite     string
		expectedPairCountAfterWrite int
		expectedError               bool
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
		{
			name:                        "test config with sync interval",
			tempFileConfig:              yamlFileConfig,
			syncInterval:                1,
			tempFileConfigNextWrite:     altYamlFileConfig.content,
			expectedPairCountAfterWrite: 1,
			want: &FileMapper{
				name:  yamlFileConfig.name,
				pairs: pairList.ToMap(),
			},
		},
		{
			name:                        "test config with sync interval, invalid file",
			tempFileConfig:              yamlFileConfig,
			syncInterval:                1,
			tempFileConfigNextWrite:     malformedYamlFileConfig.content,
			expectedPairCountAfterWrite: 2, // don't error, still keep original pairs
			want: &FileMapper{
				name:  yamlFileConfig.name,
				pairs: pairList.ToMap(),
			},
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
				Name:         tt.tempFileConfig.name,
				Path:         tmpfile.Name(),
				SyncInterval: tt.syncInterval,
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
			// make a clone and sanitize before comparison
			// clone is needed because sanitizer modifies the map, which is reused across tests
			wantClone := tt.want.pairs.Clone()
			err = sanitizer.SanitizeInputMap(fileMapper, wantClone)
			assert.NoError(t, err)
			assert.True(t, wantClone.Equals(&fileMapper.pairs), "Expected %v, got %v", wantClone, fileMapper.pairs)
			if tt.syncInterval > 0 {
				assert.NotNil(t, fileMapper.stop)
				err = os.WriteFile(tmpfile.Name(), []byte(tt.tempFileConfigNextWrite), 0o644)
				assert.NoError(t, err)
				time.Sleep(time.Duration(tt.syncInterval+1) * time.Second)
				assert.Equal(t, tt.expectedPairCountAfterWrite, len(fileMapper.pairs))
			} else {
				assert.Nil(t, fileMapper.stop)
			}

			assert.NoError(t, fileMapper.Teardown())
		})
	}
}
