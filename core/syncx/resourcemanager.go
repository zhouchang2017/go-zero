package syncx

import (
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"sync"

	"github.com/zeromicro/go-zero/core/errorx"
)

var (
	resourceManagers []*ResourceManager
	_lock            sync.Mutex
	once             sync.Once
)

// CloseAllResource 关闭所有打开的资源
// 可以在main函数里注册在proc.AddShutdownListener回调内
func CloseAllResource() {
	once.Do(func() {
		for _, manager := range resourceManagers {
			if manager.resources != nil {
				_ = manager.Close()
			}
		}
	})
}

// A ResourceManager is a manager that used to manage resources.
type ResourceManager struct {
	resources    map[string]io.Closer
	singleFlight SingleFlight
	lock         sync.RWMutex
}

// NewResourceManager returns a ResourceManager.
func NewResourceManager() *ResourceManager {
	_lock.Lock()
	defer _lock.Unlock()
	resource := &ResourceManager{
		resources:    make(map[string]io.Closer),
		singleFlight: NewSingleFlight(),
	}
	resourceManagers = append(resourceManagers, resource)
	return resource
}

// Close closes the manager.
// Don't use the ResourceManager after Close() called.
func (manager *ResourceManager) Close() error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	var be errorx.BatchError
	for name, resource := range manager.resources {
		var err error
		if err = resource.Close(); err != nil {
			logx.GlobalLogger().Errorf("%s resource close err: %s", name, err.Error())
			be.Add(err)
		}
		logx.GlobalLogger().Infof("%s resource close success", name)
	}

	// release resources to avoid using it later
	manager.resources = nil

	return be.Err()
}

// GetResource returns the resource associated with given key.
func (manager *ResourceManager) GetResource(key string, create func() (io.Closer, error)) (io.Closer, error) {
	val, err := manager.singleFlight.Do(key, func() (interface{}, error) {
		manager.lock.RLock()
		resource, ok := manager.resources[key]
		manager.lock.RUnlock()
		if ok {
			return resource, nil
		}

		resource, err := create()
		if err != nil {
			logx.GlobalLogger().Errorf("%s resource open err: %s", key, err.Error())
			return nil, err
		}
		logx.GlobalLogger().Infof("%s resource open success", key)

		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.resources[key] = resource

		return resource, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(io.Closer), nil
}

// Inject injects the resource associated with given key.
func (manager *ResourceManager) Inject(key string, resource io.Closer) {
	manager.lock.Lock()
	manager.resources[key] = resource
	manager.lock.Unlock()
}
