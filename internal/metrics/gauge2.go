package metrics

import "sync"

type Gauge2 struct {
	mu    sync.Mutex
	depth int64
	inflight int64
	acked int64
	nacked int64
	reaped int64
}

func NewGauge2() *Gauge2 { return &Gauge2{} }

func (g *Gauge2) SetDepth(n int64) {
	g.mu.Lock()
	g.depth = n
	g.mu.Unlock()
}

func (g *Gauge2) IncInflight() { g.add(&g.inflight) }
func (g *Gauge2) DecInflight() { g.addN(&g.inflight, -1) }
func (g *Gauge2) IncAcked() { g.add(&g.acked) }
func (g *Gauge2) IncNacked() { g.add(&g.nacked) }
func (g *Gauge2) IncReaped() { g.add(&g.reaped) }

func (g *Gauge2) add(p *int64) { g.addN(p, 1) }

func (g *Gauge2) addN(p *int64, n int64) {
	g.mu.Lock()
	*p += n
	g.mu.Unlock()
}

func (g *Gauge2) Snapshot() map[string]int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return map[string]int64{"depth": g.depth, "inflight": g.inflight, "acked": g.acked, "nacked": g.nacked, "reaped": g.reaped}
}
