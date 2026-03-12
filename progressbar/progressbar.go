package progressbar

import (
	"fmt"
	"strings"
	"time"

	c "github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

var firstRender = true

type ProgressBar struct {
	CurrentDepth int

	StartTime      time.Time
	DepthStartTime time.Time

	TotalAnalyzed int
	TotalSaved    int
	TotalSkipped  int
	TotalInQueue  int

	DepthAnalyzed   int
	DepthTotal      int
	DepthSaved      int
	DepthSkipped    int
	DepthAddedQueue int
}

func NewProgressBar() *ProgressBar {
	now := time.Now()
	return &ProgressBar{
		CurrentDepth:   0,
		StartTime:      now,
		DepthStartTime: now,
		DepthTotal:     1, // A mettre à jour avec la vraie taille si connue
	}
}

func (pb *ProgressBar) SetDepth(depth int) {
	pb.CurrentDepth = depth
	pb.DepthStartTime = time.Now()
	pb.DepthAnalyzed = 0
	pb.DepthSaved = 0
	pb.DepthSkipped = 0
	// pb.DepthTotal = ... (Logique pour déterminer le max à cette profondeur)
}

func (pb *ProgressBar) AddAnalyzed() {
	pb.TotalAnalyzed++
	pb.DepthAnalyzed++
}

func (pb *ProgressBar) AddSaved() {
	pb.TotalSaved++
	pb.DepthSaved++
}

func (pb *ProgressBar) AddSkipped() {
	pb.TotalSkipped++
	pb.DepthSkipped++
}

func (pb *ProgressBar) AddToQueue() {
	pb.TotalInQueue++
	pb.DepthAddedQueue++
}

func formatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (pb *ProgressBar) ResetDepth(depth int, total int) {
	pb.CurrentDepth = depth
	pb.DepthStartTime = time.Now()
	pb.DepthAnalyzed = 0
	pb.DepthSaved = 0
	pb.DepthSkipped = 0
	pb.DepthTotal = total
	pb.DepthAddedQueue = 0
}

func (pb *ProgressBar) Render() {
	totalElapsed := formatDuration(time.Since(pb.StartTime))
	depthElapsed := formatDuration(time.Since(pb.DepthStartTime))

	analyzedLine := formatAnalyzedLine(pb.DepthAnalyzed, pb.DepthTotal)

	if !firstRender {
		fmt.Printf("\033[15A")
	}

	lines := [15]string{
		"┌─── Book Generation ───────────────────────────┐",
		fmt.Sprintf("│  Current depth      : %-24d│", pb.CurrentDepth),
		fmt.Sprintf("│  Total elapsed      : %-24s│", totalElapsed),
		fmt.Sprintf("│  Depth elapsed      : %-24s│", depthElapsed),
		"├─── Total ─────────────────────────────────────┤",
		fmt.Sprintf("│  Analyzed           : %-24d│", pb.TotalAnalyzed),
		fmt.Sprintf("│  Saved              : %-24d│", pb.TotalSaved),
		fmt.Sprintf("│  Skipped            : %-24d│", pb.TotalSkipped),
		fmt.Sprintf("│  In queue           : %-24d│", pb.TotalInQueue),
		"├─── Current Depth ─────────────────────────────┤",
		analyzedLine,
		fmt.Sprintf("│  Saved              : %-24d│", pb.DepthSaved),
		fmt.Sprintf("│  Skipped            : %-24d│", pb.DepthSkipped),
		fmt.Sprintf("│  Added to queue     : %-24d│", pb.DepthAddedQueue),
		"└───────────────────────────────────────────────┘",
	}

	result := strings.Join(lines[:], "\n")
	fmt.Println(result)

	firstRender = false
}

func formatAnalyzedLine(depthAnalyzed int, depthTotalToAnalyze int) string {
	percentage := func(analyzed, total int) string {
		if total == 0 {
			return "0.00%"
		}
		return fmt.Sprintf("%.2f%%", float64(analyzed)/float64(total)*100)
	}

	s := fmt.Sprintf(
		"│  Analyzed           : %d/%d (%s)",
		depthAnalyzed,
		depthTotalToAnalyze,
		c.GreenString(percentage(depthAnalyzed, depthTotalToAnalyze)),
	)

	width := runewidth.StringWidth(s) - 6
	padding := max(49-width, 0)
	return s + strings.Repeat(" ", padding) + "│"
}
