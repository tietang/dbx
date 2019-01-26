package mapping

import "sync"

type EntityCache struct {
	cache sync.Map
}

func (c *EntityCache) Set(key string, val *EntityInfo) {
	c.cache.Store(key, val)
}
func (c *EntityCache) Get(key string) (val *EntityInfo, ok bool) {
	if v, ok := c.cache.Load(key); !ok {
		return nil, false
	} else {
		val = v.(*EntityInfo)
	}
	return val, true
}
