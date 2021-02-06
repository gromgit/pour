package console

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
	"strconv"
	"unicode/utf8"
)

type FancyString struct {
	Content, Fmt string
}

func PlainColumn(strs []FancyString) string {
	result := ""
	for _, s := range strs {
		result = result + s.Content + "\n"
	}
	return result
}

func (s FancyString) Print(strfmt string) string {
	return fmt.Sprintf("%s"+strfmt+"%s", s.Fmt, s.Content, Reset())
}

func Columnate(strs []FancyString) string {
	nstrs := len(strs)
	result := ""
	if nstrs > 0 {
		if isatty.IsTerminal(os.Stdin.Fd()) {
			// Now find the screen width
			screenwidth, err := strconv.Atoi(os.Getenv("COLUMNS"))
			if err != nil || screenwidth == 0 {
				screenwidth = 80 // A decent default
			}
			// First find the longest width
			maxwidth := 0
			for _, s := range strs {
				l := utf8.RuneCountInString(s.Content)
				if l > maxwidth {
					maxwidth = l
				}
			}
			maxwidth += 2 // Add two spaces to each column
			strfmt := fmt.Sprintf("%%-%ds", maxwidth)
			cols := screenwidth / maxwidth
			stride := nstrs / cols
			fmt.Fprintf(os.Stderr, "Fmt: %q  Screen: (%d,%d)  Cols: %d  Stride: %d\n", strfmt, 0, screenwidth, cols, stride)
			// Let's run through the elements
			if stride == 0 {
				// Single row, so just walk through strs
				for _, s := range strs {
					result = result + fmt.Sprintf(strfmt, s.Print(strfmt))
				}
				result = result + "\n"
			} else {
				// Striding through the stringies...
				for y := 0; y < stride; y++ {
					for x := 0; x < cols; x++ {
						if x*stride+y < nstrs {
							result = result + fmt.Sprintf(strfmt, strs[x*stride+y].Print(strfmt))
						}
					}
					result = result + "\n"
				}
			}
		} else {
			// Not a terminal, do single column
			result = PlainColumn(strs)
		}
	}
	return result
}
