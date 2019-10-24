package common

type List struct {
	Tracks []Track
}

func (l List) Len() int {
	return len(l.Tracks)
}

func (l List) Eq(list List) bool {
	if len(l.Tracks) != len(list.Tracks) {
		return false
	}

	for i := 0; i < len(l.Tracks); i++ {
		if !l.Tracks[i].Eq(list.Tracks[i]) {
			return false
		}
	}

	return true
}
