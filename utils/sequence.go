package utils

func MakeRange(from, to, step int) []int {
	var result []int
	for i := from; i < to+1; i += step {
		result = append(result, i)
	}
	return result
}

func Repeat[K int | string](value K, times int) []K {
	var result []K
	for i := 0; i < times; i++ {
		result = append(result, value)
	}
	return result
}
