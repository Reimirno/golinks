package file_mapper

import (
	"time"

	"github.com/reimirno/golinks/pkg/mapper"
)

const (
	FileMapperConfigType = "FILE"
)

var _ mapper.MapperConfigurer = (*FileMapperConfig)(nil)

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

func (f *FileMapperConfig) GetMapper() (mapper.Mapper, error) {
	pairs, err := parseFile(f.Path)
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

	mm := &FileMapper{
		name:  f.Name,
		pairs: pairs.ToMap(),
		stop:  stop,
	}
	for _, pair := range mm.pairs {
		mapper.Sanitize(mm, pair)
	}
	return mm, nil
}

func (f *FileMapperConfig) Singleton() bool {
	return false
}
