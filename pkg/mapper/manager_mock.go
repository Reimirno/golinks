package mapper

import (
	"github.com/reimirno/golinks/pkg/types"
	"github.com/stretchr/testify/mock"
)

// MockMapper is a mock implementation of the Mapper interface for testing purposes.
// It uses a in-memory map to store PathUrlPair objects.
// Supports all operations, but does persist anything.
type MockMapper struct {
	mock.Mock
	Pairs      types.PathUrlPairMap
	IsReadOnly bool
	Name       string
}

func (m *MockMapper) GetType() string {
	return "mock"
}

func (m *MockMapper) GetName() string {
	return m.Name
}

func (m *MockMapper) Readonly() bool {
	return m.IsReadOnly
}

func (m *MockMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	if pair, ok := m.Pairs[path]; ok {
		return pair, nil
	}
	return nil, nil
}

func (m *MockMapper) ListUrls() (types.PathUrlPairList, error) {
	return m.Pairs.ToList(), nil
}

func (m *MockMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	if m.IsReadOnly {
		return nil, ErrOperationNotSupported("put")
	}
	Sanitize(m, pair)
	m.Pairs[pair.Path] = pair
	return pair, nil
}

func (m *MockMapper) DeleteUrl(path string) error {
	if m.IsReadOnly {
		return ErrOperationNotSupported("delete")
	}
	delete(m.Pairs, path)
	return nil
}

func (m *MockMapper) Teardown() error {
	return nil
}

var _ Mapper = (*MockMapper)(nil)

// MockMapperConfigurer is a mock implementation of the MapperConfigurer interface for testing purposes.
// It just returns a new MockMapper instance.
type MockMapperConfigurer struct {
	mock.Mock
	Name         string
	IsSingleton  bool
	IsReadOnly   bool
	StarterPairs types.PathUrlPairMap
}

func (m *MockMapperConfigurer) GetType() string {
	return "mock"
}

func (m *MockMapperConfigurer) GetName() string {
	return m.Name
}

func (m *MockMapperConfigurer) GetMapper() (Mapper, error) {
	mapper := new(MockMapper)
	mapper.Pairs = m.StarterPairs
	for _, pair := range mapper.Pairs {
		pair.Mapper = m.Name
	}
	mapper.IsReadOnly = m.IsReadOnly
	mapper.Name = m.Name
	if mapper.Pairs == nil {
		mapper.Pairs = make(types.PathUrlPairMap)
	}
	return mapper, nil
}

func (m *MockMapperConfigurer) Singleton() bool {
	return m.IsSingleton
}

var _ MapperConfigurer = (*MockMapperConfigurer)(nil)

func CloneConfigurers(configurers []*MockMapperConfigurer) []MapperConfigurer {
	copied := make([]MapperConfigurer, len(configurers))
	for i, configurer := range configurers {
		element := &MockMapperConfigurer{
			Name:         configurer.Name,
			IsSingleton:  configurer.IsSingleton,
			IsReadOnly:   configurer.IsReadOnly,
			StarterPairs: make(types.PathUrlPairMap),
		}
		for path, pair := range configurer.StarterPairs {
			element.StarterPairs[path] = pair
		}
		copied[i] = element
	}
	return copied
}