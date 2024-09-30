package _map

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Map 是一个通用的映射类型
type Map[K comparable, V any] struct {
	items map[K]V
	mutex sync.RWMutex
}

// New 创建一个新的 Map
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		items: make(map[K]V),
	}
}

// Put 向映射中添加键值对
func (m *Map[K, V]) Put(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.items[key] = value
}

// Get 获取指定键的值
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, exists := m.items[key]
	return value, exists
}

// Remove 移除指定键的值
func (m *Map[K, V]) Remove(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.items, key)
}

// Clear 清空映射
func (m *Map[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.items = make(map[K]V)
}

// Size 返回映射的大小
func (m *Map[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.items)
}

// IsEmpty 检查映射是否为空
func (m *Map[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// ContainsKey 检查映射是否包含指定键
func (m *Map[K, V]) ContainsKey(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.items[key]
	return exists
}

// Keys 返回映射中所有键的切片
func (m *Map[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回映射中所有值的切片
func (m *Map[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	values := make([]V, 0, len(m.items))
	for _, v := range m.items {
		values = append(values, v)
	}
	return values
}

// ForEach 对映射中的每个键值对执行指定操作
func (m *Map[K, V]) ForEach(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for k, v := range m.items {
		f(k, v)
	}
}

// Filter 返回满足条件的键值对的新映射
func (m *Map[K, V]) Filter(f func(K, V) bool) *Map[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	result := New[K, V]()
	for k, v := range m.items {
		if f(k, v) {
			result.Put(k, v)
		}
	}
	return result
}

// Map 将映射中的每个值映射到新的值
func (m *Map[K, V]) Map(f func(K, V) V) *Map[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	result := New[K, V]()
	for k, v := range m.items {
		result.Put(k, f(k, v))
	}
	return result
}

// Reduce 将映射归约为单个值
func (m *Map[K, V]) Reduce(f func(V, K, V) V, initial V) V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	result := initial
	for k, v := range m.items {
		result = f(result, k, v)
	}
	return result
}

// MarshalJSON 实现 json.Marshaler 接口
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return json.Marshal(m.items)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return json.Unmarshal(data, &m.items)
}

// String 返回映射的字符串表示
func (m *Map[K, V]) String() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return fmt.Sprintf("%v", m.items)
}

// Equals 比较两个映射是否相等
func (m *Map[K, V]) Equals(other *Map[K, V]) bool {
	if m == other {
		return true
	}
	if m.Size() != other.Size() {
		return false
	}
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	other.mutex.RLock()
	defer other.mutex.RUnlock()
	for k, v := range m.items {
		if otherV, exists := other.items[k]; !exists || !reflect.DeepEqual(v, otherV) {
			return false
		}
	}
	return true
}

// Clone 创建映射的深拷贝
func (m *Map[K, V]) Clone() *Map[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	clone := New[K, V]()
	for k, v := range m.items {
		clone.Put(k, v)
	}
	return clone
}

// Merge 合并两个映射，如果有冲突则使用提供的解决函数
func (m *Map[K, V]) Merge(other *Map[K, V], resolver func(V, V) V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	other.mutex.RLock()
	defer other.mutex.RUnlock()
	for k, v := range other.items {
		if existingV, exists := m.items[k]; exists {
			m.items[k] = resolver(existingV, v)
		} else {
			m.items[k] = v
		}
	}
}

// GetOrDefault 获取指定键的值，如果不存在则返回默认值
func (m *Map[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, exists := m.Get(key); exists {
		return value
	}
	return defaultValue
}

// PutIfAbsent 如果键不存在，则添加键值对
func (m *Map[K, V]) PutIfAbsent(key K, value V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, exists := m.items[key]; !exists {
		m.items[key] = value
		return true
	}
	return false
}

// ComputeIfAbsent 如果键不存在，则计算并添加新值
func (m *Map[K, V]) ComputeIfAbsent(key K, mappingFunction func(K) V) (V, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if v, exists := m.items[key]; exists {
		return v, nil
	}
	newValue := mappingFunction(key)
	m.items[key] = newValue
	return newValue, nil
}

// ToJSON 将映射转换为 JSON 字符串
func (m *Map[K, V]) ToJSON() (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	data, err := json.Marshal(m.items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从 JSON 字符串创建映射
func FromJSON[K comparable, V any](jsonStr string) (*Map[K, V], error) {
	var items map[K]V
	err := json.Unmarshal([]byte(jsonStr), &items)
	if err != nil {
		return nil, err
	}
	m := New[K, V]()
	m.items = items
	return m, nil
}

// GetOrCompute 获取键对应的值，如果不存在则计算并存储
func (m *Map[K, V]) GetOrCompute(key K, computeFunc func() V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if v, exists := m.items[key]; exists {
		return v
	}
	newValue := computeFunc()
	m.items[key] = newValue
	return newValue
}

// Update 更新指定键的值
func (m *Map[K, V]) Update(key K, updateFunc func(V) V) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if v, exists := m.items[key]; exists {
		m.items[key] = updateFunc(v)
		return nil
	}
	return errors.New("键不存在")
}
