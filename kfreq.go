package kompress

type kfreq struct {
	// occurences of each byte so far
	freq [256]int
	// sum of all occurrence
	sum int
	// currently selected escape byte
	// you may read it directly
	// but only use update to update it.
	esc byte
}

// update the frequency table with the byte provided
// Will also update the esc value, as needed.
func (f *kfreq) update(b byte) {
	f.sum++
	f.freq[b]++
	if b == f.esc {
		// We need to revalidate the esc value
		for i := range f.freq {
			if f.freq[f.esc] > f.freq[i] {
				// If we found a lower occurence value, we retain it.
				// There acn not be better values, only same values,
				// so we can safely return immediately.
				f.esc = byte(i)
				// DEBUG
				// fmt.Println("\nEsc was changed to ", f.esc)
				return
			}
		}
	}
}
