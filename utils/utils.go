package utils

func Contains[T comparable](list []T, item T) bool {
	for _, x := range list {
		if x == item {
			return true
		}
	}
	return false
}

func IndexOf[T comparable](list []T, item T) int {
	for i, x := range list {
		if x == item {
			return i
		}
	}
	return -1
}

func Remove[T any](list []T, index int) []T {
	list[index] = list[len(list)-1]
	list = list[:len(list)-1]
	return list
}

type HasID interface {
	GetID() string
}
