package rbac

func IntToBinary(n, length int) []int {
	binary := make([]int, length)
	for i := 0; i < length; i++ {
		binary[i] = n % 2
		n /= 2
	}
	return binary
}
