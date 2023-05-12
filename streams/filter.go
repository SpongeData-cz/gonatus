package streams

/*
Two-sided stream, filters the data with the given filter function.

Extends:
  - InputStreamer,
  - OutputStreamer.

Type parameters:
  - T - type of the data.
*/
type FilterStreamer[T any] interface {
	InputStreamer[T]
	OutputStreamer[T]
}

/*
Transform stream.

Extends:
  - stream.

Implements:
  - InputStreamer,
  - OutputStreamer.

Type parameters:
  - T - type of the data.
*/
type filterStream[T any] struct {
	stream
	source InputStreamer[T]
	filter func(e T) bool // filter function
}

/*
Filter stream constructor.

Parameters:
  - filter - filter function.

Type parameters:
  - T - type of the data..

Returns:
  - pointer to the created filter stream.
*/
func NewFilterStream[T any](filter func(e T) bool) FilterStreamer[T] {
	ego := &filterStream[T]{
		filter: filter,
	}
	ego.init(ego)
	return ego
}

func (ego *filterStream[T]) get() (value T, valid bool, err error) {
	for true {
		value, valid, err = ego.source.get()
		closed := ego.source.Closed()
		if closed {
			ego.close()
		}
		if valid && ego.filter(value) {
			return
		}
		if closed {
			break
		}
	}
	valid = false
	return
}

func (ego *filterStream[T]) setSource(s InputStreamer[T]) {
	ego.source = s
}

func (ego *filterStream[T]) Pipe(s OutputStreamer[T]) InputStreamer[T] {
	return pipe[T](ego, s)
}

func (ego *filterStream[T]) Split(s SplitStreamer[T]) (positiveStream InputStreamer[T], negativeStream InputStreamer[T]) {
	return split[T](ego, s)
}

func (ego *filterStream[T]) Duplicate(s DuplicationStreamer[T]) (stream1 InputStreamer[T], stream2 InputStreamer[T]) {
	return duplicate[T](ego, s)
}
