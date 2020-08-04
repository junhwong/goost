package runtime

func HashCode(data []byte) uint {
	// 算法来源 java
	var h uint
	for _, b := range data {
		h = 31*h + uint(b) // TODO: 数值可能会溢出
	}
	return h
}
