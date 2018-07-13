package main

import (
	"math/big"
	"sync"
)

type Map struct {
	sync.Mutex
	m map[string]big.Int
}

func (mymap *Map) AddBigInt(key string, value big.Int) big.Int {

	zero := big.NewInt(0)

	mymap.Lock()
	addr, ok := mymap.m[key]
	if !ok {
		mymap.m[key] = *zero
	}
	mymap.Unlock()

	addr, _ = mymap.m[key]

	result := big.NewInt(0)
	result.Add(&addr, &value)
	mymap.m[key] = *result
	return *result
}

func (mymap *Map) Clone() Map {

	dbCopy := Map{
		m: map[string]big.Int{},
	}

	/* Copy Content from Map1 to Map2*/
	for index,element := range mymap.m{
		dbCopy.m[index] = element
	}

	return dbCopy
}