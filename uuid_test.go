package skyd

import (
	"testing"
)

func BenchmarkNewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUUID()
	}
}