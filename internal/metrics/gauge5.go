package metrics

import "sync"

type Gauge5 struct {
	mu    sync.Mutex
	depth int64
	inflight int64
	acked int64
	nacked int64
	reaped int64
}

func NewGauge5() *Gauge5 { return &Gauge5{} }

func (g *Gauge5) SetDepth(n int64) {
	g.mu.Lock()
	g.depth = n
	g.mu.Unlock()
}

func (g *Gauge5) IncInflight() { g.add(&g.inflight) }
func (g *Gauge5) DecInflight() { g.addN(&g.inflight, -1) }
func (g *Gauge5) IncAcked() { g.add(&g.acked) }
func (g *Gauge5) IncNacked() { g.add(&g.nacked) }
func (g *Gauge5) IncReaped() { g.add(&g.reaped) }

func (g *Gauge5) add(p *int64) { g.addN(p, 1) }

func (g *Gauge5) addN(p *int64, n int64) {
	g.mu.Lock()
	*p += n
	g.mu.Unlock()
}

func (g *Gauge5) Snapshot() map[string]int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return map[string]int64{"depth": g.depth, "inflight": g.inflight, "acked": g.acked, "nacked": g.nacked, "reaped": g.reaped}
}
