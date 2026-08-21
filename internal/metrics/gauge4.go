package metrics

import "sync"

type Gauge4 struct {
	mu    sync.Mutex
	depth int64
	inflight int64
	acked int64
	nacked int64
	reaped int64
}

func NewGauge4() *Gauge4 { return &Gauge4{} }

func (g *Gauge4) SetDepth(n int64) {
	g.mu.Lock()
	g.depth = n
	g.mu.Unlock()
}

func (g *Gauge4) IncInflight() { g.add(&g.inflight) }
func (g *Gauge4) DecInflight() { g.addN(&g.inflight, -1) }
func (g *Gauge4) IncAcked() { g.add(&g.acked) }
func (g *Gauge4) IncNacked() { g.add(&g.nacked) }
func (g *Gauge4) IncReaped() { g.add(&g.reaped) }

func (g *Gauge4) add(p *int64) { g.addN(p, 1) }

func (g *Gauge4) addN(p *int64, n int64) {
	g.mu.Lock()
	*p += n
	g.mu.Unlock()
}

func (g *Gauge4) Snapshot() map[string]int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return map[string]int64{"depth": g.depth, "inflight": g.inflight, "acked": g.acked, "nacked": g.nacked, "reaped": g.reaped}
}
