package utils

type comparableType interface {
	comparable
}

func Includes[T comparableType](arr []T, item T) bool {
	result := false

	Map[T](arr, func(_ int, i T) {
		if i == item {
			result = true
		}
	})

	return result
}
