package types

import (
	"fmt"
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

func (p PathUrlPair) String() string {
	return fmt.Sprintf("'%s' -> '%s'", p.Path, p.Url)
}

func (p *PathUrlPair) Clone() *PathUrlPair {
	return &PathUrlPair{
		Path:     p.Path,
		Url:      p.Url,
		Mapper:   p.Mapper,
		UseCount: p.UseCount,
	}
}

func (p PathUrlPairMap) ToList() PathUrlPairList {
	l := make(PathUrlPairList, 0, len(p))
	for _, pair := range p {
		l = append(l, pair.Clone())
	}
	return l
}

func (p PathUrlPairList) ToMap() PathUrlPairMap {
	m := make(PathUrlPairMap)
	for _, pair := range p {
		m[pair.Path] = pair.Clone()
	}
	return m
}

func (p *PathUrlPairList) Clone() *PathUrlPairList {
	l := make(PathUrlPairList, 0, len(*p))
	for _, pair := range *p {
		l = append(l, pair.Clone())
	}
	return &l
}

func (p *PathUrlPairMap) Clone() *PathUrlPairMap {
	m := make(PathUrlPairMap)
	for path, pair := range *p {
		m[path] = pair.Clone()
	}
	return &m
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
