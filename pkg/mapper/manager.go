package mapper

import (
	"fmt"

	"github.com/reimirno/golinks/pkg/logging"
	"github.com/reimirno/golinks/pkg/sanitizer"
	"github.com/reimirno/golinks/pkg/types"
	"go.uber.org/zap"
)

// func Sanitize(m Mapper, pair *types.PathUrlPair) {
// 	pair.Mapper = m.GetName()
// 	pair.Path = strings.Trim(pair.Path, "/")
// 	if m.Readonly() {
// 		pair.UseCount = 0
// 	}
// }

type MapperManager struct {
	mappers   []types.Mapper
	persistor types.Mapper
	logger    *zap.SugaredLogger
}

func NewMapperManager(persistorName string, mapConfigs []types.MapperConfigurer) (*MapperManager, error) {
	l := logging.NewLogger("mapper")

	m, err := validateAndGetMappers(mapConfigs)
	if err != nil {
		l.Errorf("Failed to validate mapper configurators: %v", err)
		return nil, err
	}
	pIdx, err := getPersistorIndex(persistorName, m)
	if err != nil {
		l.Errorf("Failed to get persistor: %v", err)
		return nil, err
	}

	var p types.Mapper
	if pIdx >= 0 {
		p = m[pIdx]
		if pIdx > 0 {
			l.Warn("Persistor is not the first mapper; some overrides may not work as expected")
		}
	} else {
		l.Warn("No persistor is configured")
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
			m.logger.Errorf("Failed to teardown mapper %s: %v", mapper.GetName(), err)
			return err
		}
	}
	m.logger.Debugf("All mappers are down")
	return nil
}

func (m *MapperManager) GetUrl(path string, incrementCounter bool) (*types.PathUrlPair, error) {
	m.logger.Debugf("Getting url: %s", path)
	canonicalPath, err := sanitizer.CanonicalizePath(path)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Path canonicalized: %s -> %s", path, canonicalPath)
	// mapper order is important here
	// mappers in the front takes precedence over mappers in the back
	for _, mapper := range m.mappers {
		m.logger.Debugf("Trying mapper %s for path %s", mapper.GetName(), canonicalPath)
		pair, err := mapper.GetUrl(canonicalPath)
		if err != nil {
			m.logger.Errorf("Failed to get url at mapper %s: %v", mapper.GetName(), err)
			return nil, err
		}
		if pair != nil {
			m.logger.Debugf("Mapper %s used", mapper.GetName())
			if incrementCounter && !mapper.Readonly() {
				m.logger.Debugf("Try to increment counter at mapper %s: %d -> %d", mapper.GetName(), pair.UseCount, pair.UseCount+1)
				pair.UseCount = pair.UseCount + 1
				_, err = mapper.PutUrl(pair)
				if err != nil {
					m.logger.Errorf("Failed to increment counter at mapper %s: %v", mapper.GetName(), err)
				}
			}
			sanitizer.SanitizeOutput(mapper, pair)
			return pair, nil
		}
	}
	m.logger.Debugf("No mapper is available for path %s (raw: %s)", canonicalPath, path)
	return nil, nil
}

func (m *MapperManager) ListUrls() (types.PathUrlPairList, error) {
	m.logger.Debugf("Listing urls")
	// mapper order is important here
	// mappers in the front takes precedence over mappers in the back
	var urlMap = make(types.PathUrlPairMap)
	for _, mapper := range m.mappers {
		urls, err := mapper.ListUrls()
		if err != nil {
			return nil, err
		}
		for _, url := range urls {
			if urlMap[url.Path] == nil {
				sanitizer.SanitizeOutput(mapper, url)
				urlMap[url.Path] = url
			}
		}
	}
	m.logger.Debugf("found %d urls", len(urlMap))
	return urlMap.ToList(), nil
}

func (m *MapperManager) getPersistor() types.Mapper {
	return m.persistor
}

func (m *MapperManager) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	m.logger.Debugf("Setting url: %s -> %s", pair.Path, pair.Url)
	if m.getPersistor() == nil {
		return nil, ErrOperationNotSupported("set")
	}
	canonicalPath, err := sanitizer.CanonicalizePath(pair.Path)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Path canonicalized: %s -> %s", pair.Path, canonicalPath)
	old, err := m.GetUrl(canonicalPath, false)
	if err != nil {
		return nil, err
	}
	if old == nil {
		// Create path
		pair.UseCount = 0
		persistor := m.getPersistor()
		err = sanitizer.SanitizeInput(persistor, pair)
		if err != nil {
			return nil, err
		}
		pair, err = persistor.PutUrl(pair)
		if err != nil {
			return nil, err
		}
		sanitizer.SanitizeOutput(persistor, pair)
		return pair, nil
	}
	// Update path
	mapper := findMapper(m.mappers, old.Mapper)
	if mapper == nil {
		return nil, ErrInvalidMapper(old.Mapper)
	}
	err = sanitizer.SanitizeInput(mapper, pair)
	if err != nil {
		return nil, err
	}
	pair, err = mapper.PutUrl(pair)
	if err != nil {
		return nil, err
	}
	sanitizer.SanitizeOutput(mapper, pair)
	return pair, nil
}

func (m *MapperManager) DeleteUrl(path string) error {
	m.logger.Debugf("Deleting url: %s", path)
	if m.getPersistor() == nil {
		return ErrOperationNotSupported("delete")
	}
	canonicalPath, err := sanitizer.CanonicalizePath(path)
	if err != nil {
		return err
	}
	m.logger.Debugf("Path canonicalized: %s -> %s", path, canonicalPath)
	old, err := m.GetUrl(canonicalPath, false)
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
	return mapper.DeleteUrl(canonicalPath)
}

func validateAndGetMappers(mapConfigs []types.MapperConfigurer) ([]types.Mapper, error) {
	if len(mapConfigs) == 0 {
		return nil, ErrMapConfigSetup("no mappers are configured")
	}
	mappers := make([]types.Mapper, 0, len(mapConfigs))
	typesAppeared := make(map[string]bool)
	namesAppeared := make(map[string]bool)
	for _, cfg := range mapConfigs {
		if cfg == nil {
			return nil, ErrMapConfigSetup("mapper configurator is nil")
		}
		if typesAppeared[cfg.GetType()] && cfg.Singleton() {
			return nil, ErrMapConfigSetup(fmt.Sprintf("duplicate singleton mapper type: %s", cfg.GetType()))
		}
		if namesAppeared[cfg.GetName()] {
			return nil, ErrMapConfigSetup(fmt.Sprintf("duplicate mapper name: %s", cfg.GetName()))
		}
		typesAppeared[cfg.GetType()] = true
		namesAppeared[cfg.GetName()] = true
		mapper, err := cfg.GetMapper()
		if err != nil {
			return nil, ErrMapConfigSetup(fmt.Sprintf("failed to get mapper for config %s: %v", cfg.GetName(), err.Error()))
		}
		mappers = append(mappers, mapper)
	}
	return mappers, nil
}

func getPersistorIndex(persistorName string, mappers []types.Mapper) (int, error) {
	if persistorName == "" {
		return -1, nil
	}
	persistorIdx := findMapperIndex(mappers, persistorName)
	if persistorIdx < 0 {
		return -1, ErrMapConfigSetup(fmt.Sprintf("persistor not found: %s", persistorName))
	}
	if mappers[persistorIdx].Readonly() {
		return -1, ErrMapConfigSetup(fmt.Sprintf("persistor is readonly: %s", persistorName))
	}
	return persistorIdx, nil
}

func findMapperIndex(mappers []types.Mapper, name string) int {
	for idx, mapper := range mappers {
		if mapper.GetName() == name {
			return idx
		}
	}
	return -1
}

func findMapper(mappers []types.Mapper, name string) types.Mapper {
	idx := findMapperIndex(mappers, name)
	if idx < 0 {
		return nil
	}
	return mappers[idx]
}
