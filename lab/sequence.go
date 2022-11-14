package lab

func MakeRange(from, to, step int) []int {
	var result []int
	for i := from; i < to+1; i += step {
		result = append(result, i)
	}
	return result
}
