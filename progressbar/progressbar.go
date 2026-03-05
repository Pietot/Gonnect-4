package progressbar

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const DisplayLines = 15

type Progress struct {
	mu sync.RWMutex

	currentDepth     int
	totalSaved       int64
	totalSkipped     int64
	totalAnalyzed    int64
	totalEnqueued    int64
	totalDequeued    int64
	initialQueueSize int64

	depthSaved      map[int]int64
	depthSkipped    map[int]int64
	depthAnalyzed   map[int]int64
	depthEnqueued   map[int]int64
	depthDispatched map[int]int64

	globalStart time.Time
	depthStart  time.Time

	rendered int32 // atomic: 0 = not yet rendered
}

func New() *Progress {
	now := time.Now()
	return &Progress{
		depthSaved:      make(map[int]int64),
		depthSkipped:    make(map[int]int64),
		depthAnalyzed:   make(map[int]int64),
		depthEnqueued:   make(map[int]int64),
		depthDispatched: make(map[int]int64),
		globalStart:     now,
		depthStart:      now,
	}
}

func (p *Progress) InitQueueSize(n int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.initialQueueSize = n
}

func (p *Progress) RecordSaved(depth int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.totalSaved++
	p.totalAnalyzed++
	p.depthSaved[depth]++
	p.depthAnalyzed[depth]++
	if depth > p.currentDepth {
		p.currentDepth = depth
		p.depthStart = time.Now()
	}
}

func (p *Progress) RecordSkipped(depth int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.totalSkipped++
	p.totalAnalyzed++
	p.depthSkipped[depth]++
	p.depthAnalyzed[depth]++
	if depth > p.currentDepth {
		p.currentDepth = depth
		p.depthStart = time.Now()
	}
}

func (p *Progress) RecordEnqueued(depth int, count int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.totalEnqueued += count
	p.depthEnqueued[depth] += count
}

func (p *Progress) RecordDequeued(count int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.totalDequeued += count
}

func (p *Progress) RecordDispatched(depth int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.depthDispatched[depth]++
}

func (p *Progress) Render() {
	p.mu.RLock()
	d := p.currentDepth
	totalSaved := p.totalSaved
	totalSkipped := p.totalSkipped
	totalAnalyzed := p.totalAnalyzed
	queueSize := p.initialQueueSize + p.totalEnqueued - p.totalDequeued
	dSaved := p.depthSaved[d]
	dSkipped := p.depthSkipped[d]
	dAnalyzed := p.depthAnalyzed[d]
	dDispatched := p.depthDispatched[d]
	dEnqueued := p.depthEnqueued[d]
	globalElapsed := time.Since(p.globalStart)
	depthElapsed := time.Since(p.depthStart)
	p.mu.RUnlock()

	if atomic.SwapInt32(&p.rendered, 1) == 1 {
		fmt.Printf("\033[%dA", DisplayLines)
	}

	analyzedLine := formatAnalyzedLine(d, dAnalyzed, dDispatched)

	lines := [DisplayLines]string{
		"┌─── Book Generation ───────────────────────────┐",
		fmt.Sprintf("│  Current depth          : %-20d│", d),
		fmt.Sprintf("│  Total elapsed          : %-20s│", formatDuration(globalElapsed)),
		fmt.Sprintf("│  Depth elapsed          : %-20s│", formatDuration(depthElapsed)),
		"├───────────────────────────────────────────────┤",
		fmt.Sprintf("│  Total saved            : %-20d│", totalSaved),
		fmt.Sprintf("│  Total skipped          : %-20d│", totalSkipped),
		fmt.Sprintf("│  Total analyzed         : %-20d│", totalAnalyzed),
		fmt.Sprintf("│  Total in queue         : %-20d│", queueSize),
		"├───────────────────────────────────────────────┤",
		fmt.Sprintf("│  [D:%d] Saved           : %-21d│", d, dSaved),
		fmt.Sprintf("│  [D:%d] Skipped         : %-21d│", d, dSkipped),
		analyzedLine,
		fmt.Sprintf("│  [D:%d] Added to queue  : %-21d│", d, dEnqueued),
		"└───────────────────────────────────────────────┘",
	}

	for _, line := range lines {
		fmt.Printf("\033[2K%s\n", line)
	}
}

func formatAnalyzedLine(d int, analyzed, dispatched int64) string {
	if dispatched <= 0 {
		return fmt.Sprintf("│  [D:%d] Analyzed        : %-21d│", d, analyzed)
	}
	pct := analyzed * 100 / dispatched
	pctStr := fmt.Sprintf("%d%%", pct)
	visible := fmt.Sprintf("%d/%d (%s)", analyzed, dispatched, pctStr)
	colored := fmt.Sprintf("%d/%d (\033[32m%s\033[0m)", analyzed, dispatched, pctStr)
	padding := max(21-len(visible), 0)
	return fmt.Sprintf("│  [D:%d] Analyzed        : %s%s│", d, colored, strings.Repeat(" ", padding))
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
