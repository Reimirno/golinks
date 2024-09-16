package file_mapper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"reimirno.com/golinks/pkg/types"
)

func parseFileInput(file string) (types.PathUrlPairList, error) {
	content, err := read(file)
	if err != nil {
		return nil, err
	}
	pointers, err := parse(content)
	if err != nil {
		return nil, err
	}
	return pointers, nil
}

type fileContent struct {
	fileType string
	content  []byte
}

func read(file string) (*fileContent, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	return &fileContent{
		fileType: strings.TrimPrefix(filepath.Ext(info.Name()), "."),
		content:  bytes,
	}, nil
}

func parse(cnt *fileContent) ([]*types.PathUrlPair, error) {
	var parsed []types.PathUrlPair
	var pointers []*types.PathUrlPair
	var err error
	switch cnt.fileType {
	case "yaml":
		err = yaml.Unmarshal(cnt.content, &parsed)
	case "yml":
		err = yaml.Unmarshal(cnt.content, &parsed)
	case "json":
		err = json.Unmarshal(cnt.content, &parsed)
	default:
		err = fmt.Errorf("unsupported file type: %s", cnt.fileType)
	}
	if err != nil {
		return nil, err
	}
	for _, pair := range parsed {
		pair.Mapper = FileMapperConfigType
		pointers = append(pointers, &pair)
	}
	return pointers, nil
}
