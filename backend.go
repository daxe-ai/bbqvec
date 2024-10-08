package bbq

import (
	"errors"

	"github.com/RoaringBitmap/roaring"
)

type VectorBackend interface {
	PutVector(id ID, v Vector) error
	ComputeSimilarity(targetVector Vector, targetID ID) (float32, error)
	Info() BackendInfo
	Exists(id ID) bool
	Close() error
}

type scannableBackend interface {
	VectorBackend
	ForEachVector(func(ID) error) error
}

type VectorGetter[T any] interface {
	GetVector(id ID) (T, error)
}

type IndexBackend interface {
	SaveBases(bases []Basis, token uint64) (uint64, error)
	LoadBases() ([]Basis, error)

	SaveBitmap(basis int, index int, bitmap *roaring.Bitmap) error
	LoadBitmap(basis, index int) (*roaring.Bitmap, error)
	Sync() error
}

type BackendInfo struct {
	HasIndexData bool
	Dimensions   int
	Quantization string
}

func FullTableScanSearch(be VectorBackend, target Vector, k int) (*ResultSet, error) {
	rs := NewResultSet(k)
	b, ok := be.(scannableBackend)
	if !ok {
		return nil, errors.New("Backend is incompatible")
	}
	err := b.ForEachVector(func(id ID) error {
		sim, err := b.ComputeSimilarity(target, id)
		if err != nil {
			return err
		}
		rs.AddResult(id, sim)
		return nil
	})
	return rs, err
}
