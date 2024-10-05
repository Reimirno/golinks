package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/reimirno/golinks/pkg/mapper"
	bolt_mapper "github.com/reimirno/golinks/pkg/mapper/bolt-mapper"
	file_mapper "github.com/reimirno/golinks/pkg/mapper/file-mapper"
	mem_mapper "github.com/reimirno/golinks/pkg/mapper/mem-mapper"
	sql_mapper "github.com/reimirno/golinks/pkg/mapper/sql-mapper"
	"github.com/spf13/viper"
)

type config struct {
	Server serverConfig `mapstructure:"server"`
	Mapper mapperConfig `mapstructure:"mapper"`
}

type serverConfig struct {
	Port struct {
		Redirector string `mapstructure:"redirector"`
		Crud       string `mapstructure:"crud"`
		CrudHttp   string `mapstructure:"crud_http"`
	} `mapstructure:"port"`
	Debug bool `mapstructure:"debug"`
}

type mapperConfig struct {
	Persistor string                    `mapstructure:"persistor"`
	Mappers   []mapperConfigurerWrapper `mapstructure:"mappers"`
}

func NewConfig(configFile string) (*config, error) {
	var config config

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetDefault("Server.Port.Redirector", "8080")
	v.SetDefault("Server.Port.Crud", "8081")
	v.SetDefault("Server.Port.CrudHttp", "8082")
	v.SetDefault("Server.Debug", false)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := v.Unmarshal(&config, viper.DecodeHook((&mapperConfigurerWrapper{}).DecodeMapstructure(nil))); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if config.Server.Debug {
		fmt.Printf("Config loaded: %+v\n", config)
	}

	return &config, nil
}

type mapperConfigurerWrapper struct {
	Type             string `mapstructure:"type"`
	MapperConfigurer mapper.MapperConfigurer
}

func (w *mapperConfigurerWrapper) DecodeMapstructure(config *mapstructure.DecoderConfig) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(mapperConfigurerWrapper{}) {
			return data, nil
		}

		raw, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data type for MapperConfigurerWrapper")
		}

		mapperType, ok := raw["type"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'type' field in mapper")
		}

		var wrapper mapperConfigurerWrapper
		wrapper.Type = mapperType

		switch strings.ToUpper(mapperType) {
		case file_mapper.FileMapperConfigType:
			var fileMapper file_mapper.FileMapperConfig
			if err := mapstructure.Decode(raw, &fileMapper); err != nil {
				return nil, err
			}
			wrapper.MapperConfigurer = &fileMapper
		case bolt_mapper.BoltMapperConfigType:
			var boltMapper bolt_mapper.BoltMapperConfig
			if err := mapstructure.Decode(raw, &boltMapper); err != nil {
				return nil, err
			}
			wrapper.MapperConfigurer = &boltMapper
		case mem_mapper.MemMapperConfigType:
			var memMapper mem_mapper.MemMapperConfig
			if err := mapstructure.Decode(raw, &memMapper); err != nil {
				return nil, err
			}
			wrapper.MapperConfigurer = &memMapper
		case sql_mapper.SqlMapperConfigType:
			var sqlMapper sql_mapper.SqlMapperConfig
			if err := mapstructure.Decode(raw, &sqlMapper); err != nil {
				return nil, err
			}
			wrapper.MapperConfigurer = &sqlMapper
		default:
			return nil, fmt.Errorf("unknown mapper type: %s", mapperType)
		}

		return wrapper, nil
	}
}
