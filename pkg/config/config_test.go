package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tempFileConfig struct {
	name    string
	content string
}

var (
	yamlConfigFileContent = &tempFileConfig{
		name: "test.yaml",
		content: `
server:
  port:
    redirector: 8080
    crud: 8081
  debug: true

mapper:
  persistor: ""
  mappers:
    - type: mem
      name: memory
      pairs:
        - path: ggl
          url: https://google.com
`}
	invalidConfigFileContent = &tempFileConfig{
		name:    "test.yaml",
		content: "invalid content",
	}
)

func createTempFile(tempFileConfig tempFileConfig) (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "*"+tempFileConfig.name)
	if err != nil {
		return nil, err
	}
	_, err = tmpfile.Write([]byte(tempFileConfig.content))
	if err != nil {
		return nil, err
	}
	tmpfile.Close()
	return tmpfile, nil
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		tempFileConfig *tempFileConfig
		wantError      bool
		redirectorPort string
		crudPort       string
		debug          bool
		numMappers     int
	}{
		{
			name:           "happy path",
			tempFileConfig: yamlConfigFileContent,
			redirectorPort: "8080",
			crudPort:       "8081",
			debug:          true,
			numMappers:     1,
		},
		{
			name:           "invalid config file",
			tempFileConfig: invalidConfigFileContent,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := createTempFile(*tt.tempFileConfig)
			assert.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			cfg, err := NewConfig(tmpfile.Name())
			if tt.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.redirectorPort, cfg.Server.Port.Redirector)
			assert.Equal(t, tt.crudPort, cfg.Server.Port.Crud)
			assert.Equal(t, tt.debug, cfg.Server.Debug)
			assert.Equal(t, tt.numMappers, len(cfg.Mapper.Mappers))
		})
	}
}
