package structures

import (
	"errors"
	"io"

	"github.com/SpongeData-cz/gonatus"
)

// RAMTape is Taper implemented by Slice
type RAMTape[T comparable] struct {
	gonatus.Gobject
	Slice  []T
	Offset int
	closed bool
}

func NewRAMTape[T comparable](conf gonatus.Conf) *RAMTape[T] {
	//TODO initial capacity?
	ego := &RAMTape[T]{}
	ego.Init(ego, conf)
	return ego
}

//TODO implementations of Taper methods - shoulld we use Taper documentation or is there anything extra to write down?

// TODO do we really want Seeker from IO? if so then offsets should be int64 everywhere
func (ego *RAMTape[T]) Seek(offset int, whence int) (int, error) {

	switch whence {
	case io.SeekStart:
		if offset < 0 {
			return -1, errors.New("Negative offset")
		}
		if offset >= len(ego.Slice) {
			return -1, errors.New("Too large offset")
		}
		ego.Offset = offset //TODO ugly conversion
		break
	case io.SeekCurrent:
		newOffset := ego.Offset + offset
		if newOffset < 0 {
			return -1, errors.New("Negative offset")
		}
		if newOffset >= len(ego.Slice) {
			return -1, errors.New("Too large offset")
		}
		ego.Offset = newOffset
		break
	case io.SeekEnd:
		newOffset := len(ego.Slice) - 1 - offset
		if newOffset < 0 {
			return -1, errors.New("Negative offset")
		}
		if newOffset >= len(ego.Slice) {
			return -1, errors.New("Too large offset")
		}
		ego.Offset = newOffset
	}
	return ego.Offset, nil
}

func (ego *RAMTape[T]) Close() error {
	ego.closed = true
	return nil
}

func (ego *RAMTape[T]) Append(item ...T) {
	ego.Slice = append(ego.Slice, item...)
}

func (ego *RAMTape[T]) Read(p []T) (n int, err error) {
	itemsRed := 0
	egoSliceLen := len(ego.Slice)
	bufferLen := len(p)
	for i := ego.Offset; i-ego.Offset < bufferLen && i < egoSliceLen; i++ {
		p[itemsRed] = ego.Slice[i]
		itemsRed++
	}
	ego.Offset += itemsRed
	return itemsRed, nil //TODO if there is nothing more to read, should we return "EOF" error? Or only if there is nothing less and moreover we are closed?

	//TODO should read move offset?

}

func (ego *RAMTape[T]) Filter(dest Taper[T], fn func(T) bool) error {
	var res []T
	for _, it := range ego.Slice {
		if fn(it) {
			res = append(res, it)
		}
	}
	dest.Append(res...)
	return nil
}

func (ego *RAMTape[T]) Closed() bool {
	return ego.closed
}