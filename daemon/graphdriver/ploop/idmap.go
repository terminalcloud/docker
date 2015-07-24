// +build linux

package ploop

import (
	"sync"
)

type IdMap struct {
	sync.Mutex
	ids map[string]string
}
