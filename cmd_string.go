package woodis

import (
	"github.com/Novicehood/woodis/server"
	"time"
)

// command args
type optString struct {
	key   string
	value string
	ttl   time.Duration
	nx    bool
	ex    bool
	get   bool
}

// register the cmd func in server
func commandString(wd *WooDis) {

}

func (wd *WooDis) cmdSet(p *server.Peer, cmd string, args []string) {

}
