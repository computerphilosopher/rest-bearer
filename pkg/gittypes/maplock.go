package gittypes

import (
	"fmt"
	"sync"
)

type MapLock struct {
	locks sync.Map
}

func (mapLock *MapLock) Load(repository Repository) (*sync.RWMutex, error) {
	key := repository.ToString()
	raw, exist := mapLock.locks.Load(key)
	if !exist {
		mutex := &sync.RWMutex{}
		mapLock.locks.Store(key, mutex)
		return mutex, nil
	}

	mutex, ok := raw.(*sync.RWMutex)
	if !ok {
		return nil,
			fmt.Errorf("%s's lock type is wrong expected: mutex actual: %T",
				key, raw)
	}

	return mutex, nil
}
