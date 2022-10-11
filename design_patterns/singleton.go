package dp

import (
	"sync"
)

// GLI is singleton
var GLI sync.Mutex

func init() {
	GLI = sync.Mutex{}
}
