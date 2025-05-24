package woodis

import (
	"context"
	"github.com/Novicehood/woodis/server"
	"math/rand"
	"sync"
	"time"
)

type WooDis struct {
	sync.Mutex

	dbs        map[int]*RedisDB
	srv        *server.Server
	selectedDB int
	passwords  map[string]string
	scripts    map[string]string
	now        time.Time
	Ctx        context.Context
	CtxCancel  context.CancelFunc
	signal     sync.Cond
	rand       *rand.Rand
}
