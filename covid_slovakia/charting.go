package covid_slovakia

import (
	"github.com/guptarohit/asciigraph"
)

// PrintLineASCII Prints the actual chart, gets the data as well as the labelstring to put in.
func PrintLineASCII(data []float64, dateStringStart, dateStringEnd string) string {
	graphLabel := NormalizeXAxis(dateStringStart, dateStringEnd)
	var chart string

	// The length 22 is the magical number for correct formatting on iPhones (tested on iPhone 11 Pro and iPhone SE 2020)
	// that's why we wrap it around when it's bigger than that. Height is 15
	chart = asciigraph.Plot(data, asciigraph.Width(22), asciigraph.Height(15), asciigraph.Caption(graphLabel))

	return chart
}

// NormalizeXAxis Inserts the X axis-like line at least, since there is no X axis...
func NormalizeXAxis(startDate, endDate string) string {
	return "―――――――――――――――――――――\n\t" + startDate + " <-> " + endDate
}

// GetGraphReadyForDiscordPrint Simplifies the printout for reuse
func GetGraphReadyForDiscordPrint(input string) string {
	return "**\n```go\n" + input + "```"
}
