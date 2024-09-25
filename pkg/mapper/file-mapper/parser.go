package file_mapper

import (
	"github.com/reimirno/golinks/pkg/types"
	"github.com/spf13/viper"
)

type pathUrlPairWrapper struct {
	Data []types.PathUrlPair `yaml:"data" json:"data"`
}

func parseFile(file string) (types.PathUrlPairList, error) {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var parsed pathUrlPairWrapper
	if err := v.Unmarshal(&parsed); err != nil {
		return nil, err
	}
	pairs := make(types.PathUrlPairList, len(parsed.Data))
	for i, pair := range parsed.Data {
		pairs[i] = &pair
	}
	return pairs, nil
}
