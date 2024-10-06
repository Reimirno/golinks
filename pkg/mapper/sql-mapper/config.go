package sql_mapper

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/reimirno/golinks/pkg/types"
)

const (
	SqlMapperConfigType = "SQL"
)

var _ types.MapperConfigurer = (*SqlMapperConfig)(nil)

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

func (m *SqlMapperConfig) GetMapper() (types.Mapper, error) {
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
