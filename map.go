package main

import (
	"math/big"
	"sync"
	"unsafe"
	"sync/atomic"
)

type Map struct {
	sync.Mutex
	m map[string]**big.Int
}

func (mymap *Map) AddBigInt(key string, value big.Int) big.Int {

	zero := big.NewInt(0)

	mymap.Lock()
	addr, ok := mymap.m[key]
	if !ok {
		mymap.m[key] = &zero
	}
	mymap.Unlock()

	addr, _ = mymap.m[key]

	done := false
	result := big.NewInt(0)
	for !done {
		oldValue := *addr
		result.Add(oldValue, &value)
		done = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(addr)), unsafe.Pointer(oldValue), unsafe.Pointer(result))
	}
	return *result
}