package metrics

import "sync"

type Registry struct {
	mu        sync.Mutex
	enqueued  int64
	claimed   int64
	acked     int64
	nacked    int64
	reaped    int64
	heartbeat int64
}

func New() *Registry { return &Registry{} }

func (r *Registry) IncEnqueued()  { r.add(&r.enqueued) }
func (r *Registry) IncClaimed()   { r.add(&r.claimed) }
func (r *Registry) IncAcked()     { r.add(&r.acked) }
func (r *Registry) IncNacked()    { r.add(&r.nacked) }
func (r *Registry) IncReaped()    { r.add(&r.reaped) }
func (r *Registry) IncHeartbeat() { r.add(&r.heartbeat) }

func (r *Registry) add(p *int64) {
	r.mu.Lock()
	*p++
	r.mu.Unlock()
}

func (r *Registry) Snapshot() map[string]int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return map[string]int64{
		"enqueued": r.enqueued, "claimed": r.claimed, "acked": r.acked,
		"nacked": r.nacked, "reaped": r.reaped, "heartbeat": r.heartbeat,
	}
}
