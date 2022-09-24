package plugin

import (
	"hash"
)

type fakeHash struct{}

func (f *fakeHash) Write(b []byte) (int, error) {
	return len(b), nil
}

func (f *fakeHash) Sum(b []byte) []byte {
	return append(b, []byte("fake")...)
}

func (f *fakeHash) Reset() {
}

func (f *fakeHash) Size() int {
	return 4
}

func (f *fakeHash) BlockSize() int {
	return 4096
}

var _ hash.Hash = &fakeHash{}
