package slices

type accepted interface {
	~string | int | uint
}

func Contains[T accepted](slice []T, sortKey T) bool {
	return Index(slice, sortKey) >= 0
}

func Index[T accepted](slice []T, sortKey T) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == sortKey {
			return i
		}
	}

	return -1
}
