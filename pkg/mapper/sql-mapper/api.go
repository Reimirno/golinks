package sql_mapper

import (
	"fmt"

	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	SqlMapperConfigType = "SQL"
)

var _ mapper.MapperConfigurer = (*SqlMapperConfig)(nil)

type SqlMapperConfig struct {
	Name   string `mapstructure:"name"`
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

func (m *SqlMapperConfig) GetName() string {
	return m.Name
}

func (m *SqlMapperConfig) GetType() string {
	return SqlMapperConfigType
}

func (m *SqlMapperConfig) GetMapper() (mapper.Mapper, error) {
	var db *gorm.DB
	var err error
	switch m.Driver {
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(m.DSN), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(m.DSN), &gorm.Config{})
	case "postgres":
		// disable prepared statement cache, otherwise auto-migration fails on the second run forward for some reason
		db, err = gorm.Open(postgres.New(postgres.Config{DSN: m.DSN, PreferSimpleProtocol: true}), &gorm.Config{})
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(m.DSN), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported driver: %s", m.Driver)
	}
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&types.PathUrlPair{})
	if err != nil {
		return nil, err
	}
	return &SqlMapper{name: m.Name, db: db}, nil
}

func (m *SqlMapperConfig) Singleton() bool {
	return false
}

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
