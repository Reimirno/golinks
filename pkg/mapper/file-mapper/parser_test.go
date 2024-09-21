package file_mapper

import (
	"os"
	"testing"

	"github.com/reimirno/golinks/pkg/types"
	"github.com/stretchr/testify/assert"
)

type tempFileConfig struct {
	name    string
	content string
}

var (
	pairList = &types.PathUrlPairList{
		{
			Path: "example",
			Url:  "https://example.com",
		},
		{
			Path: "test",
			Url:  "https://test.com",
		},
	}

	yamlFileConfig = &tempFileConfig{
		name: "test.yaml",
		content: `
data:
  - path: "example"
    url: "https://example.com"
  - path: "test"
    url: "https://test.com"
`,
	}

	jsonFileConfig = &tempFileConfig{
		name: "test.json",
		content: `
{
  "data": [
    {
      "path": "example",
      "url": "https://example.com"
    },
    {
      "path": "test",
      "url": "https://test.com"
    }
  ]
}
`,
	}

	malformedYamlFileConfig = &tempFileConfig{
		name: "test-malformed.yaml",
		content: `
data:
  - path: "example" asdf
    url: "https://example.com"
  - --
`,
	}

	emptyFileConfig = &tempFileConfig{
		name:    "test.yaml",
		content: "",
	}

	invalidFileConfig = &tempFileConfig{
		name:    "test.yaml",
		content: "invalid content",
	}

	unsupportedFileConfig = &tempFileConfig{
		name:    "test.txt",
		content: "unsupported content",
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

func TestParseFile(t *testing.T) {
	tests := []struct {
		name           string
		tempFileConfig *tempFileConfig
		expectedPairs  *types.PathUrlPairList
		expectedError  bool
	}{
		{
			name:           "yaml file",
			tempFileConfig: yamlFileConfig,
			expectedPairs:  pairList,
		},
		{
			name:           "json file",
			tempFileConfig: jsonFileConfig,
			expectedPairs:  pairList,
		},
		{
			name:           "malformed yaml file",
			tempFileConfig: malformedYamlFileConfig,
			expectedError:  true,
		},
		{
			name:           "empty file",
			tempFileConfig: emptyFileConfig,
			expectedError:  false,
		},
		{
			name:           "invalid file",
			tempFileConfig: invalidFileConfig,
			expectedError:  true,
		},
		{
			name:           "unsupported file",
			tempFileConfig: unsupportedFileConfig,
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpfile, err := createTempFile(*test.tempFileConfig)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			pairs, err := parseFile(tmpfile.Name())
			if test.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, test.expectedPairs.Equals(&pairs))
		})
	}
}
