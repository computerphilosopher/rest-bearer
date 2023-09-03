package maplock

import (
	"fmt"
	"sync"

	"github.com/computerphilosopher/rest-bearer/pkg/gittypes"
)

type MapLock struct {
	locks sync.Map
}

var instance *MapLock = nil
var once sync.Once

func GetMapLock() *MapLock {
	once.Do(func() {
		instance = &MapLock{}
	})
	return instance
}

func (mapLock *MapLock) Load(repository gittypes.Repository) (*sync.RWMutex, error) {
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
