package console

import (
	"fmt"
	"os"
	"strconv"
)

func Columnate(strs []string) string {
	nstrs := len(strs)
	if nstrs == 0 {
		return "" // Nothing to do
	}
	// First find the longest width
	maxwidth := 0
	for _, s := range strs {
		l := len(s)
		if l > maxwidth {
			maxwidth = l + 2 // Have at least two spaces for each column
		}
	}
	strfmt := fmt.Sprintf("%%-%ds", maxwidth)
	// Now find the screen width
	screenwidth, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil || screenwidth == 0 {
		screenwidth = 80 // A decent default
	}
	cols := screenwidth / maxwidth
	stride := nstrs / cols
	fmt.Fprintf(os.Stderr, "Fmt: %q  Cols: %d  Stride: %d\n", strfmt, cols, stride)
	// Let's run through the elements
	result := ""
	if stride == 0 {
		// Single row, so just walk through strs
		for _, s := range strs {
			result = result + fmt.Sprintf(strfmt, s)
		}
		result = result + "\n"
	} else {
		// Striding through the stringies...
		for y := 0; y < stride; y++ {
			for x := 0; x < cols; x++ {
				if x*stride+y < nstrs {
					result = result + fmt.Sprintf(strfmt, strs[x*stride+y])
				}
			}
			result = result + "\n"
		}
	}
	return result
}
