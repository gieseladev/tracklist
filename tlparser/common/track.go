package common

type Track struct {
	StartOffsetMS uint32
	EndOffsetMS   uint32
	Name          string
}

func (t Track) HasEndOffset() bool {
	return t.EndOffsetMS != 0
}

func (t Track) DurationMS() uint32 {
	if t.EndOffsetMS >= t.StartOffsetMS {
		return 0
	}

	return t.EndOffsetMS - t.StartOffsetMS
}

func (t Track) Eq(track Track) bool {
	return t.StartOffsetMS == track.StartOffsetMS &&
		t.EndOffsetMS == track.EndOffsetMS &&
		t.Name == track.Name
}
