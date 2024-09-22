package sql_mapper

import (
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
	"gorm.io/gorm"
)

type SqlMapper struct {
	name string
	db   *gorm.DB
}

var _ mapper.Mapper = (*SqlMapper)(nil)

func (m *SqlMapper) GetName() string {
	return m.name
}

func (m *SqlMapper) GetType() string {
	return SqlMapperConfigType
}

func (m *SqlMapper) Teardown() error {
	return nil
}

func (m *SqlMapper) GetUrl(path string) (*types.PathUrlPair, error) {
	var pair types.PathUrlPair
	err := m.db.Where("path = ?", path).Take(&pair).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	mapper.Sanitize(m, &pair)
	return &pair, nil
}

func (m *SqlMapper) PutUrl(pair *types.PathUrlPair) (*types.PathUrlPair, error) {
	err := m.db.Save(pair).Error
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (m *SqlMapper) DeleteUrl(path string) error {
	err := m.db.Where("path = ?", path).Delete(&types.PathUrlPair{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *SqlMapper) ListUrls() (types.PathUrlPairList, error) {
	var pairs types.PathUrlPairList
	err := m.db.Find(&pairs).Error
	if err != nil {
		return nil, err
	}
	for _, pair := range pairs {
		mapper.Sanitize(m, pair)
	}
	return pairs, nil
}

func (m *SqlMapper) Readonly() bool {
	return false
}