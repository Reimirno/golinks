package file_mapper

import (
	"go.uber.org/zap"

	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
	"github.com/reimirno/golinks/pkg/utils"
)

var _ types.Mapper = (*FileMapper)(nil)

type FileMapper struct {
	logger *zap.SugaredLogger
	name   string
	pairs  types.PathUrlPairMap
	stop   func()
}

func (f *FileMapper) GetType() string {
	return FileMapperConfigType
}

func (f *FileMapper) GetName() string {
	return f.name
}

func (f *FileMapper) Teardown() error {
	if f.stop != nil {
		f.stop()
	}
	return nil
}

func (f *FileMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	pair, ok := f.pairs[path]
	if !ok {
		return nil, nil
	}
	return pair, nil
}

func (f *FileMapper) ListUrls(pagination types.Pagination) (types.PathUrlPairList, error) {
	return utils.Paginate(f.pairs.ToList(), pagination), nil
}

func (f *FileMapper) DeleteUrl(path string) error {
	return mapper.ErrOperationNotSupported("delete")
}

func (f *FileMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	return nil, mapper.ErrOperationNotSupported("put")
}

func (f *FileMapper) Readonly() bool {
	return true
}
