package file_mapper

import (
	"fmt"

	"github.com/spf13/viper"
	"reimirno.com/golinks/pkg/types"
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
	fmt.Println(pairs)
	return pairs, nil
}
