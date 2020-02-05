package kompress

func getTestData() [][]byte {
	data := [][]byte{
		{1}, {1},
		{1, 2}, {1, 2},
		{1, 2, 3}, {1, 2, 3},
		{1, 2, 3, 4}, {1, 2, 3, 4},

		{}, {},

		{0}, {0, 0},
		{0, 0}, {0, 1, 0},
		{0, 0, 0}, {0, 2, 0},
		{0, 0, 0, 0}, {0, 3, 0},

		{1, 3, 3, 3}, {1, 0, 2, 3},
		{1, 3, 3, 3, 4}, {1, 0, 2, 3, 4},
		{1, 3, 3, 3, 0}, {1, 0, 2, 3, 0, 0},
		{0, 1, 3, 3, 3, 0}, {0, 0, 1, 0, 2, 3, 0, 0},
		{0, 3, 3, 3, 0}, {0, 0, 0, 2, 3, 0, 0},

		{1, 8, 8, 5}, {1, 8, 8, 5},
		{8, 8, 5}, {8, 8, 5},
		{8, 8, 0}, {8, 8, 0, 0},
		{5, 8, 8}, {5, 8, 8},
		{0, 8, 8}, {0, 0, 8, 8},

		{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}, {1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3},
	}

	// Add some long sequences
	seq := []byte{22, 22}
	for i := 2; i <= 255; i++ {
		seq = append(seq, 22)
		data = append(data, seq)
		data = append(data, []byte{0, byte((i % 256)), 22})
	}
	// 256 and beyond ...
	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 22})

	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 22, 22})

	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 0, 2, 22})

	seq = append(seq, 22)
	data = append(data, seq)
	data = append(data, []byte{0, 255, 22, 0, 3, 22})

	return data
}
