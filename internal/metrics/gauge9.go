package metrics

import "sync"

type Gauge9 struct {
	mu    sync.Mutex
	depth int64
	inflight int64
	acked int64
	nacked int64
	reaped int64
}

func NewGauge9() *Gauge9 { return &Gauge9{} }

func (g *Gauge9) SetDepth(n int64) {
	g.mu.Lock()
	g.depth = n
	g.mu.Unlock()
}

func (g *Gauge9) IncInflight() { g.add(&g.inflight) }
func (g *Gauge9) DecInflight() { g.addN(&g.inflight, -1) }
func (g *Gauge9) IncAcked() { g.add(&g.acked) }
func (g *Gauge9) IncNacked() { g.add(&g.nacked) }
func (g *Gauge9) IncReaped() { g.add(&g.reaped) }

func (g *Gauge9) add(p *int64) { g.addN(p, 1) }

func (g *Gauge9) addN(p *int64, n int64) {
	g.mu.Lock()
	*p += n
	g.mu.Unlock()
}

func (g *Gauge9) Snapshot() map[string]int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return map[string]int64{"depth": g.depth, "inflight": g.inflight, "acked": g.acked, "nacked": g.nacked, "reaped": g.reaped}
}
