package utils

type MapCallback[T any] func(index int, item T)

func Map[T any](arr []T, callback MapCallback[T]) {
	for index, it := range arr {
		callback(index, it)
	}
}
