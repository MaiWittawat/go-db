package redis

import "fmt"

type KeyGenerator struct {
	prefix string	
}

func NewKeyGenerator(prefix string) *KeyGenerator {
	return &KeyGenerator{prefix: prefix}
}

func (k *KeyGenerator) KeyList() string {
	return k.prefix
}

func (k *KeyGenerator) KeyID(id string) string {
	return fmt.Sprintf("%s:%s", k.prefix, id)
}

func (k *KeyGenerator) KeyField(field string, value string) string {
	return fmt.Sprintf("%s:%s", k.prefix, field)
}