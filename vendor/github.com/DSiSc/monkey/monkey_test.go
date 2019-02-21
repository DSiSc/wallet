package monkey

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func no() bool  { return false }
func yes() bool { return true }

func TestTimePatch(t *testing.T) {
	before := time.Now()
	Patch(time.Now, func() time.Time {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	})
	during := time.Now()
	assert.True(t, Unpatch(time.Now))
	after := time.Now()

	assert.Equal(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), during)
	assert.NotEqual(t, before, during)
	assert.NotEqual(t, during, after)
}

func TestGC(t *testing.T) {
	value := true
	Patch(no, func() bool {
		return value
	})
	defer UnpatchAll()
	runtime.GC()
	assert.True(t, no())
}

func TestSimple(t *testing.T) {
	assert.False(t, no())
	Patch(no, yes)
	assert.True(t, no())
	assert.True(t, Unpatch(no))
	assert.False(t, no())
	assert.False(t, Unpatch(no))
}

func TestGuard(t *testing.T) {
	var guard *PatchGuard
	guard = Patch(no, func() bool {
		guard.Unpatch()
		defer guard.Restore()

		return !no()
	})
	for i := 0; i < 100; i++ {
		assert.True(t, no())
	}
	Unpatch(no)
}

func TestUnpatchAll(t *testing.T) {
	assert.False(t, no())
	Patch(no, yes)
	assert.True(t, no())
	UnpatchAll()
	assert.False(t, no())
}

type s struct{}

func (s *s) yes() bool { return true }

func TestWithInstanceMethod(t *testing.T) {
	i := &s{}

	assert.False(t, no())
	Patch(no, i.yes)
	assert.True(t, no())
	Unpatch(no)
	assert.False(t, no())
}

type f struct{}

func (f *f) no() bool { return false }

func TestOnInstanceMethod(t *testing.T) {
	i := &f{}
	assert.False(t, i.no())
	PatchInstanceMethod(reflect.TypeOf(i), "no", func(_ *f) bool { return true })
	assert.True(t, i.no())
	assert.True(t, UnpatchInstanceMethod(reflect.TypeOf(i), "no"))
	assert.False(t, i.no())
}

func TestNotFunction(t *testing.T) {
	assert.Panics(t, func() {
		Patch(no, 1)
	})
	assert.Panics(t, func() {
		Patch(1, yes)
	})
}

func TestNotCompatible(t *testing.T) {
	assert.Panics(t, func() {
		Patch(no, func() {})
	})
}
