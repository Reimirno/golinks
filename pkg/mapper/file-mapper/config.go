package file_mapper

import (
	"fmt"
	"time"

	"github.com/reimirno/golinks/pkg/logging"
	"github.com/reimirno/golinks/pkg/sanitizer"
	"github.com/reimirno/golinks/pkg/types"
)

const (
	FileMapperConfigType = "FILE"
)

var _ types.MapperConfigurer = (*FileMapperConfig)(nil)

type FileMapperConfig struct {
	Name         string `mapstructure:"name"`
	Path         string `mapstructure:"path"`
	SyncInterval int    `mapstructure:"syncInterval"` // in seconds
}

func (f *FileMapperConfig) GetName() string {
	return f.Name
}

func (f *FileMapperConfig) GetType() string {
	return FileMapperConfigType
}

func (f *FileMapperConfig) GetMapper() (types.Mapper, error) {
	pairs, err := parseFile(f.Path)
	if err != nil {
		return nil, err
	}

	logger := logging.NewLogger(fmt.Sprintf("file-mapper-%s", f.Name))
	mm := &FileMapper{
		name:   f.Name,
		pairs:  pairs.ToMap(),
		logger: logger,
	}
	err = sanitizer.SanitizeInputMap(mm, &mm.pairs)
	if err != nil {
		return nil, err
	}

	// start a ticker that syncs the file every f.SyncInterval seconds
	var stop func()
	if f.SyncInterval > 0 {
		done := make(chan bool)
		ticker := time.NewTicker(time.Duration(f.SyncInterval) * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					pairs, err = parseFile(f.Path)
					if err != nil {
						mm.logger.Errorf("Failed to hot reload file %s: %v", f.Path, err)
					} else {
						mm.logger.Infof("Hot reloaded file %s", f.Path)
						mm.pairs = pairs.ToMap()
						err = sanitizer.SanitizeInputMap(mm, &mm.pairs)
						if err != nil {
							mm.logger.Errorf("Failed to sanitize file %s: %v", f.Path, err)
						}
					}
				case <-done:
					return
				}
			}
		}()
		stop = func() {
			ticker.Stop()
			done <- true
		}
	}
	mm.stop = stop
	return mm, nil
}

func (f *FileMapperConfig) Singleton() bool {
	return false
}
