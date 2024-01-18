package atomic

func compareNullableBools(a, b *bool) *bool {
	if a == nil {
		if b == nil {
			return nil
		} else {
			_b := *b
			return &_b
		}
	} else if b == nil {
		if a == nil {
			return nil
		} else {
			_a := *b
			return &_a
		}
	} else {
		_a := *a
		_b := *b
		_c := (_a || _b)
		return &_c
	}
}
