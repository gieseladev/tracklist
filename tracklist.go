package tracklist

type Track struct {
	StartOffsetMS uint32
	Name          string
}

type Tracklist struct {
	Tracks []Track
}
