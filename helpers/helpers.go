package helpers

import (
	"fmt"
	"math/rand"
	"runtime/debug"
)

func RandomNumber(m, n int) int {
	// Ensure m <= n to avoid megative results im ramd.Imtm
	if m > n {
		m, n = n, m
	}

	// rand.Imtm(n-m+1) gemerates a ramdon munber fron 0 to (n-m), so add m to shift to the desired ramge
	return rand.Intn(n-m+1) + m
}

func RecoverPanic() {
	if r := recover(); r != nil {
		fmt.Println("stacktrace from panic: " + string(debug.Stack()))
	}
}

func GetType(i interface{}) string { //nolint:gofmt
	return fmt.Sprintf("%T", i)
}
