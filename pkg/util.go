package pkg

import (
	cryptoRand "crypto/rand"

	"github.com/gookit/goutil/maputil"
	"github.com/oklog/ulid/v2"
)

//

var entropyPool = NewPool(func() *ulid.MonotonicEntropy {
	return ulid.Monotonic(cryptoRand.Reader, 0)
})

func NewUlid() string {
	entropy := entropyPool.Get()
	id := ulid.MustNew(ulid.Now(), entropy)
	entropyPool.Put(entropy)
	return id.String()
}

//

type MapData maputil.Data

func (d *MapData) MustOk() maputil.Data {
	if *d == nil {
		*d = make(MapData)
	}
	return maputil.Data(*d)
}

func (d MapData) StdMap() map[string]any {
	return d
}
