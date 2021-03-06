package config

import (
	"fmt"
)

type ScalarConfigProvider struct {
	weight int
	name   string
	kv     *KeyValue
}

var _ Provider = (*ScalarConfigProvider)(nil)

func NewScalarConfigProvider(kv *KeyValue, repo *Repository, weight int) (*ScalarConfigProvider, error) {
	p := &ScalarConfigProvider{
		weight: weight,
		kv:     kv,
		name:   fmt.Sprintf("scalar-provider-%s", kv.Key),
	}
	repo.RegisterKey(kv.Key, p)

	return p, nil
}

func (s *ScalarConfigProvider) Name() string {
	return s.name
}

func (s *ScalarConfigProvider) Depends() []string {
	return []string{}
}

func (*ScalarConfigProvider) SetUp(repo *Repository) error {
	return nil
}

func (*ScalarConfigProvider) TearDown(*Repository) error {
	return nil
}

func (s *ScalarConfigProvider) Get(key Key) (*KeyValue, bool) {
	if key.Equals(s.kv.Key) {
		return s.kv, true
	}

	return nil, false
}

func (s *ScalarConfigProvider) Weight() int {
	return s.weight
}
