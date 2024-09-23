package types

import (
	"sort"
)

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

func (p *PathUrlPair) Equals(other *PathUrlPair) bool {
	// ignore Mapper and UseCount
	if p == nil && other == nil {
		return true
	}
	if p == nil || other == nil {
		return false
	}
	return p.Path == other.Path && p.Url == other.Url
}

func (m *PathUrlPairMap) Equals(other *PathUrlPairMap) bool {
	if m == nil && other == nil {
		return true
	}
	if m == nil || other == nil {
		return false
	}
	if len(*m) != len(*other) {
		return false
	}
	for key, pair := range *m {
		otherPair, ok := (*other)[key]
		if !ok || !pair.Equals(otherPair) {
			return false
		}
	}
	return true
}

func (l *PathUrlPairList) Equals(other *PathUrlPairList) bool {
	if l == nil && other == nil {
		return true
	}
	if l == nil || other == nil {
		return false
	}
	if len(*l) != len(*other) {
		return false
	}
	// orderless comparison
	sort.SliceStable(*l, func(i, j int) bool {
		return (*l)[i].Path < (*l)[j].Path
	})
	sort.SliceStable(*other, func(i, j int) bool {
		return (*other)[i].Path < (*other)[j].Path
	})
	for i, pair := range *l {
		if *pair != *(*other)[i] {
			return false
		}
	}
	return true
}

type Service interface {
	Start(errChan chan<- error)
	Stop() error
	GetName() string
}
