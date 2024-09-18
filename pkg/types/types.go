package types

type PathUrlPairMap map[string]*PathUrlPair

type PathUrlPairList []*PathUrlPair

type PathUrlPair struct {
	Path     string `yaml:"path" json:"path" gorm:"primaryKey"`
	Url      string `yaml:"url" json:"url" gorm:"not null"`
	Mapper   string `gorm:"-"`
	UseCount int    `gorm:"not null;default:0"`
}

func (p PathUrlPairMap) ToList() PathUrlPairList {
	l := make(PathUrlPairList, 0, len(p))
	for _, pair := range p {
		l = append(l, pair)
	}
	return l
}

func (p PathUrlPairList) ToMap() PathUrlPairMap {
	m := make(PathUrlPairMap)
	for _, pair := range p {
		m[pair.Path] = pair
	}
	return m
}

type Service interface {
	Start(errChan chan<- error)
	Stop() error
	GetName() string
}
