package cassette

import (
	"fmt"
	"strings"
)

type readoutDef struct {
	Label       string
	Unit        string
	Baseline    float64
	Variance    float64
	IntervalMs  int     // 0 = static (no Alpine update)
	FaultBelow  float64 // 0 = no fault threshold
	WarnBelow   float64 // 0 = no warn threshold
	StaticValue string  // used when IntervalMs == 0
	StaticLabel string  // status label for static readouts
}

var instrumentReadouts = []readoutDef{
	{Label: "CABIN PRESSURE", Unit: "kPa", Baseline: 101.3, Variance: 0.2, IntervalMs: 2500},
	{Label: "O\u2082 PARTIAL PRESS.", Unit: "kPa", Baseline: 21.3, Variance: 0.15, IntervalMs: 2500},
	{Label: "CABIN TEMP.", Unit: "\u00b0C", Baseline: 19.7, Variance: 0.25, IntervalMs: 2500},
	{Label: "REACTOR OUTPUT", Unit: "%", Baseline: 98.7, Variance: 1.0, IntervalMs: 2500, FaultBelow: 80},
	{Label: "VELOCITY", Unit: "km/s", StaticValue: "12.4", StaticLabel: "CRUISE"},
	{Label: "FUEL REMAINING", Unit: "%", Baseline: 67.4, Variance: 0.1, IntervalMs: 3000, WarnBelow: 30},
}

func (r readoutDef) isDynamic() bool { return r.IntervalMs > 0 }

func (r readoutDef) alpineData() string {
	if !r.isDynamic() {
		return ""
	}
	return fmt.Sprintf("{ v: %g }", r.Baseline)
}

func (r readoutDef) alpineInit() string {
	if !r.isDynamic() {
		return ""
	}
	return fmt.Sprintf("setInterval(()=>{ v=parseFloat((%g+(Math.random()-0.5)*%g).toFixed(1)) },%d)",
		r.Baseline, r.Variance*2, r.IntervalMs)
}

func (r readoutDef) hasThreshold() bool {
	return r.FaultBelow > 0 || r.WarnBelow > 0
}

// valueClassExpr returns the Alpine :class expression for the readout value.
func (r readoutDef) valueClassExpr() string {
	if r.FaultBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-readout-value cass-readout-value-danger' : 'cass-readout-value'", r.FaultBelow)
	}
	if r.WarnBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-readout-value cass-readout-value-danger' : 'cass-readout-value'", r.WarnBelow)
	}
	return ""
}

// statusClassExpr returns the Alpine :class expression for the status label.
func (r readoutDef) statusClassExpr() string {
	if r.FaultBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-value-danger' : 'cass-value-ok'", r.FaultBelow)
	}
	if r.WarnBelow > 0 {
		return fmt.Sprintf("v < %g ? 'cass-value-warn' : 'cass-value-ok'", r.WarnBelow)
	}
	return ""
}

// statusTextExpr returns the Alpine x-text expression for the status label.
func (r readoutDef) statusTextExpr() string {
	if r.FaultBelow > 0 {
		return fmt.Sprintf("v < %g ? 'FAULT' : 'NOMINAL'", r.FaultBelow)
	}
	if r.WarnBelow > 0 {
		return fmt.Sprintf("v < %g ? 'LOW' : 'NOMINAL'", r.WarnBelow)
	}
	return ""
}

type barGaugeDef struct {
	Label      string
	AlpineVar  string
	InitialVal int
	FillClass  string
	ValueClass string
}

var barGauges = []barGaugeDef{
	{"FUEL CELLS", "fuel", 67, "", "cass-value"},
	{"OXYGEN RESERVES", "o2", 89, "cass-progress-fill-green", "cass-value"},
	{"COOLANT LEVEL", "coolant", 23, "cass-progress-fill-red", "cass-value-danger"},
	{"POWER DISTRIBUTION", "power", 78, "", "cass-value"},
}

// barGaugeAlpineData returns the x-data expression for the bar gauge panel.
func barGaugeAlpineData() string {
	var b strings.Builder
	b.WriteString("{ ")
	for i, g := range barGauges {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s: %d", g.AlpineVar, g.InitialVal)
	}
	b.WriteString(" }")
	return b.String()
}

// barGaugeRefreshExpr returns the Alpine @click expression to randomize all gauges.
func barGaugeRefreshExpr() string {
	var b strings.Builder
	for i, g := range barGauges {
		if i > 0 {
			b.WriteString("; ")
		}
		fmt.Fprintf(&b, "%s=Math.floor(Math.random()*60+30)", g.AlpineVar)
	}
	return b.String()
}
