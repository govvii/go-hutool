package list

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"sync"
)

// List 是一个通用的列表类型
type List[T any] struct {
	items []T
	mutex sync.RWMutex
}

// NewEmpty 创建一个新的空 List
func NewEmpty[T any]() *List[T] {
	return &List[T]{items: []T{}}
}

// New 创建一个新的 空List
func New[T any](items ...T) *List[T] {
	return &List[T]{items: items}
}

// Add 向列表末尾添加一个元素
func (l *List[T]) Add(item T) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = append(l.items, item)
}

// AddAll 向列表末尾添加多个元素
func (l *List[T]) AddAll(items ...T) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = append(l.items, items...)
}

// Get 获取指定索引的元素
func (l *List[T]) Get(index int) (T, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, errors.New("索引越界")
	}
	return l.items[index], nil
}

// Set 设置指定索引的元素
func (l *List[T]) Set(index int, item T) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if index < 0 || index >= len(l.items) {
		return errors.New("索引越界")
	}
	l.items[index] = item
	return nil
}

// Remove 移除指定索引的元素
func (l *List[T]) Remove(index int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if index < 0 || index >= len(l.items) {
		return errors.New("索引越界")
	}
	l.items = append(l.items[:index], l.items[index+1:]...)
	return nil
}

// Clear 清空列表
func (l *List[T]) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = []T{}
}

// Size 返回列表的大小
func (l *List[T]) Size() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return len(l.items)
}

// IsEmpty 检查列表是否为空
func (l *List[T]) IsEmpty() bool {
	return l.Size() == 0
}

// Contains 检查列表是否包含指定元素
func (l *List[T]) Contains(item T) bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for _, v := range l.items {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

// IndexOf 返回指定元素在列表中的索引，如果不存在则返回 -1
func (l *List[T]) IndexOf(item T) int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for i, v := range l.items {
		if reflect.DeepEqual(v, item) {
			return i
		}
	}
	return -1
}

// ToSlice 返回列表的切片副本
func (l *List[T]) ToSlice() []T {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return append([]T{}, l.items...)
}

// ForEach 对列表中的每个元素执行指定操作
func (l *List[T]) ForEach(f func(T)) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for _, item := range l.items {
		f(item)
	}
}

// Filter 返回满足条件的元素列表
func (l *List[T]) Filter(f func(T) bool) *List[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	result := New[T]()
	for _, item := range l.items {
		if f(item) {
			result.Add(item)
		}
	}
	return result
}

// Map 将列表中的每个元素映射到新的值
func (l *List[T]) Map(f func(T) T) *List[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	result := New[T]()
	for _, item := range l.items {
		result.Add(f(item))
	}
	return result
}

// Reduce 将列表归约为单个值
func (l *List[T]) Reduce(f func(T, T) T, initial T) T {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	result := initial
	for _, item := range l.items {
		result = f(result, item)
	}
	return result
}

// Sort 对列表进行排序（要求 T 类型实现了 sort.Interface）
func (l *List[T]) Sort(less func(i, j T) bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	sort.Slice(l.items, func(i, j int) bool {
		return less(l.items[i], l.items[j])
	})
}

// Reverse 反转列表
func (l *List[T]) Reverse() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for i, j := 0, len(l.items)-1; i < j; i, j = i+1, j-1 {
		l.items[i], l.items[j] = l.items[j], l.items[i]
	}
}

// Shuffle 随机打乱列表
func (l *List[T]) Shuffle() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	rand.Shuffle(len(l.items), func(i, j int) {
		l.items[i], l.items[j] = l.items[j], l.items[i]
	})
}

// String 返回列表的字符串表示
func (l *List[T]) String() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return fmt.Sprintf("%v", l.items)
}

// Equals 比较两个列表是否相等
func (l *List[T]) Equals(other *List[T]) bool {
	if l == other {
		return true
	}
	if l.Size() != other.Size() {
		return false
	}
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	other.mutex.RLock()
	defer other.mutex.RUnlock()
	for i, v := range l.items {
		if !reflect.DeepEqual(v, other.items[i]) {
			return false
		}
	}
	return true
}

// Clone 创建列表的深拷贝
func (l *List[T]) Clone() *List[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return New(l.items...)
}

// Concat 连接两个列表
func (l *List[T]) Concat(other *List[T]) *List[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	other.mutex.RLock()
	defer other.mutex.RUnlock()
	return New(append(l.items, other.items...)...)
}

// Unique 返回去重后的列表
func (l *List[T]) Unique() *List[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	seen := make(map[interface{}]struct{})
	result := New[T]()
	for _, item := range l.items {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result.Add(item)
		}
	}
	return result
}
