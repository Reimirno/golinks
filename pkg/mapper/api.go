package mapper

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"reimirno.com/golinks/pkg/logging"
	"reimirno.com/golinks/pkg/types"
)

type MapperIdentityProvider interface {
	GetType() string
	GetName() string
}

type MapperConfigurer interface {
	MapperIdentityProvider
	GetMapper() (Mapper, error)
	Singleton() bool
}

type Mapper interface {
	MapperIdentityProvider
	Readonly() bool

	GetUrl(path string) (*types.PathUrlPair, error)
	ListUrls() (types.PathUrlPairList, error)
	PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error)
	DeleteUrl(path string) error

	Teardown() error
}

func Sanitize(m Mapper, pair *types.PathUrlPair) {
	pair.Mapper = m.GetName()
	pair.Path = strings.Trim(pair.Path, "/")
	if m.Readonly() {
		pair.UseCount = 0
	}
}

type MapperManager struct {
	mappers   []Mapper
	persistor Mapper
	logger    *zap.SugaredLogger
}

func NewMapperManager(persistorName string, mapConfigs []MapperConfigurer) (*MapperManager, error) {
	l := logging.NewLogger("mapper")

	m, err := validateAndGetMappers(mapConfigs)
	if err != nil {
		l.Errorf("failed to validate mapper configurators: %v", err)
		return nil, err
	}
	p, err := getPersistor(persistorName, m)
	if err != nil {
		l.Errorf("failed to get persistor: %v", err)
		return nil, err
	}

	return &MapperManager{
		mappers:   m,
		persistor: p,
		logger:    l,
	}, nil
}

func (m *MapperManager) Teardown() error {
	for _, mapper := range m.mappers {
		err := mapper.Teardown()
		if err != nil {
			m.logger.Errorf("failed to teardown mapper %s: %v", mapper.GetName(), err)
			return err
		}
	}
	m.logger.Debugf("All mappers are down")
	return nil
}

func (m *MapperManager) GetUrl(path string, incrementCounter bool) (*types.PathUrlPair, error) {
	for _, mapper := range m.mappers {
		m.logger.Debugf("Trying mapper %s for path %s", mapper.GetName(), path)
		pair, err := mapper.GetUrl(path)
		if err != nil {
			m.logger.Errorf("Failed to get url at mapper %s: %v", mapper.GetName(), err)
			return nil, err
		}
		if pair != nil {
			m.logger.Debugf("Mapper %s used", mapper.GetName())
			if incrementCounter {
				m.logger.Debugf("Try to increment counter at mapper %s: %d -> %d", mapper.GetName(), pair.UseCount, pair.UseCount+1)
				pair.UseCount = pair.UseCount + 1
				_, err = mapper.PutUrl(pair)
				if err != nil {
					m.logger.Errorf("Failed to increment counter at mapper %s: %v", mapper.GetName(), err)
				}
			}
			return pair, nil
		}
	}
	m.logger.Debugf("No mapper is available for path: %s", path)
	return nil, nil
}

func (m *MapperManager) ListUrls() (types.PathUrlPairList, error) {
	var allUrls []*types.PathUrlPair
	for _, mapper := range m.mappers {
		urls, err := mapper.ListUrls()
		if err != nil {
			return nil, err
		}
		allUrls = append(allUrls, urls...)
	}
	m.logger.Debugf("found %d urls", len(allUrls))
	return allUrls, nil
}

func (m *MapperManager) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	m.logger.Debugf("Setting url: %s -> %s", pair.Path, pair.Url)
	if m.persistor == nil {
		return nil, ErrOperationNotSupported("set")
	}
	old, err := m.GetUrl(pair.Path, false)
	if err != nil {
		return nil, err
	}
	if old == nil {
		pair.UseCount = 0
		return m.persistor.PutUrl(pair)
	}
	mapper := findMapper(m.mappers, old.Mapper)
	if mapper == nil {
		return nil, ErrInvalidMapper(old.Mapper)
	}
	return mapper.PutUrl(pair)
}

func (m *MapperManager) DeleteUrl(path string) error {
	m.logger.Debugf("Deleting url: %s", path)
	if m.persistor == nil {
		return ErrOperationNotSupported("delete")
	}
	old, err := m.GetUrl(path, false)
	if err != nil {
		return err
	}
	if old == nil {
		return nil
	}
	mapper := findMapper(m.mappers, old.Mapper)
	if mapper == nil {
		return ErrInvalidMapper(old.Mapper)
	}
	return mapper.DeleteUrl(path)
}

func validateAndGetMappers(mapConfigs []MapperConfigurer) ([]Mapper, error) {
	mappers := make([]Mapper, 0, len(mapConfigs))
	typesAppeared := make(map[string]bool)
	for _, cfg := range mapConfigs {
		if cfg == nil {
			return nil, ErrMapConfigSetup("mapper configurator is nil")
		}
		if typesAppeared[cfg.GetType()] && cfg.Singleton() {
			return nil, ErrMapConfigSetup(fmt.Sprintf("duplicate singleton mapper type: %s", cfg.GetType()))
		}
		typesAppeared[cfg.GetType()] = true
		mapper, err := cfg.GetMapper()
		if err != nil {
			return nil, ErrMapConfigSetup(fmt.Sprintf("failed to get mapper for config %s: %v", cfg.GetName(), err.Error()))
		}
		mappers = append(mappers, mapper)
	}
	return mappers, nil
}

func getPersistor(persistorName string, mappers []Mapper) (Mapper, error) {
	if persistorName == "" {
		return nil, nil
	}
	persistor := findMapper(mappers, persistorName)
	if persistor == nil {
		return nil, ErrMapConfigSetup(fmt.Sprintf("persistor not found: %s", persistorName))
	}
	if persistor.Readonly() {
		return nil, ErrMapConfigSetup(fmt.Sprintf("persistor is readonly: %s", persistorName))
	}
	return persistor, nil
}

func findMapper(mappers []Mapper, name string) Mapper {
	for _, mapper := range mappers {
		if mapper.GetName() == name {
			return mapper
		}
	}
	return nil
}
