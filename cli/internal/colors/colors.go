package colors

import (
	"os"
	"strings"
)

var (
	// Japanese minimalism inspired colors
	Reset    = "\033[0m"
	Subtle   = "\033[90m"  // Dim gray for secondary text
	Primary  = "\033[37m"  // Clean white for primary text  
	Success  = "\033[32m"  // Soft green for success
	Warning  = "\033[33m"  // Amber for warnings
	Error    = "\033[31m"  // Red for errors
	Accent   = "\033[36m"  // Cyan accent for highlights
	Bold     = "\033[1m"   // Bold for emphasis
)

func init() {
	// Disable colors if not a terminal or NO_COLOR is set
	if os.Getenv("NO_COLOR") != "" || !isTerminal() {
		Reset = ""
		Subtle = ""
		Primary = ""
		Success = ""
		Warning = ""
		Error = ""
		Accent = ""
		Bold = ""
	}
}

func isTerminal() bool {
	// Simple terminal detection
	return os.Getenv("TERM") != ""
}

// Colorize helpers
func Gray(text string) string {
	return Subtle + text + Reset
}

func White(text string) string {
	return Primary + text + Reset
}

func Green(text string) string {
	return Success + text + Reset
}

func Yellow(text string) string {
	return Warning + text + Reset
}

func Red(text string) string {
	return Error + text + Reset
}

func Cyan(text string) string {
	return Accent + text + Reset
}

func BoldWhite(text string) string {
	return Bold + Primary + text + Reset
}

// Status indicators
func SuccessIcon() string {
	return Success + "✓" + Reset
}

func ErrorIcon() string {
	return Error + "✗" + Reset
}

func InfoIcon() string {
	return Accent + "•" + Reset
}

// Table formatting
func TableHeader(text string) string {
	return Bold + Subtle + strings.ToUpper(text) + Reset
}

func TableRow(columns []string, widths []int) string {
	var parts []string
	for i, col := range columns {
		if i == 0 {
			// First column (usually ID) in accent color
			parts = append(parts, Cyan(col))
		} else if i == len(columns)-1 {
			// Last column (usually status/date) in subtle color
			parts = append(parts, Gray(col))
		} else {
			// Middle columns in primary color
			parts = append(parts, Primary+col+Reset)
		}
	}
	return strings.Join(parts, "   ")
}

// Command output formatting
func CommandTitle(text string) string {
	return Bold + Primary + text + Reset
}

func FieldLabel(text string) string {
	return Gray(text + ":")
}

func FieldValue(text string) string {
	return Primary + text + Reset
}