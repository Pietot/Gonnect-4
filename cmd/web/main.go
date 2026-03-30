//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"

	"github.com/Pietot/Gonnect-4/config"
	"github.com/Pietot/Gonnect-4/grid"
)

func analyze(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return map[string]any{"ok": false, "error": "missing sequence"}
	}
	sequence := args[0].String()
	config.IsBookEnabled = true

	g, err := grid.InitGrid(sequence)
	if err != nil {
		return map[string]any{"ok": false, "error": err.Error()}
	}

	analysis, _ := g.Analyze()

	return map[string]any{
		"ok": true,
		"analysis": map[string]any{
			"scores":         analysis.Scores,
			"remainingMoves": analysis.RemainingMoves,
		},
	}
}

func main() {
	js.Global().Set("gonnectAnalyze", js.FuncOf(analyze))
	select {}
}
