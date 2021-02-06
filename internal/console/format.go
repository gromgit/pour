package console

import (
	"fmt"
	"github.com/gromgit/goncurses"
	"github.com/mattn/go-isatty"
	"os"
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
	if nstrs == 0 {
		goto nothing
	}
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		// Not a terminal, do single column
		return PlainColumn(strs)
	}
	if stdscr, err := goncurses.Init(); err != nil {
		panic(err)
	} else {
		defer goncurses.End()
		// On a terminal, so let's get fancy
		screenheight, screenwidth := stdscr.MaxYX()
		if screenwidth == 0 {
			panic(fmt.Sprintf("ARGH! MaxYX = (%d,%d)", screenheight, screenwidth))
		}
		// First find the longest width
		maxwidth := 0
		for _, s := range strs {
			l := utf8.RuneCountInString(s.Content)
			if l > maxwidth {
				maxwidth = l + 2 // Have at least two spaces for each column
			}
		}
		strfmt := fmt.Sprintf("%%-%ds", maxwidth)
		cols := screenwidth / maxwidth
		stride := nstrs / cols
		fmt.Fprintf(os.Stderr, "Fmt: %q  Screen: (%d,%d)  Cols: %d  Stride: %d\n", strfmt, screenheight, screenwidth, cols, stride)
		// Let's run through the elements
		result := ""
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
		return result
	}
nothing:
	return "" // Nothing to do
}
