package lazy

import "sync"

// Lazy 是一个泛型懒加载器，负责按需初始化类型 T
type Lazy[T any] struct {
	once  sync.Once
	value T
	err   error
	init  func() (T, error)
}

// NewLazy 创建一个懒加载器，传入初始化函数
func NewLazy[T any](init func() (T, error)) *Lazy[T] {
	return &Lazy[T]{init: init}
}

// Get 获取值（如果尚未初始化则调用 init 函数）
func (l *Lazy[T]) Get() (T, error) {
	l.once.Do(func() {
		l.value, l.err = l.init()
	})
	return l.value, l.err
}
