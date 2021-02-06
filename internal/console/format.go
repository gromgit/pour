package console

import (
	"fmt"
	"github.com/gromgit/pour/internal/config"
	"os"
	"strings"
	"unicode/utf8"
)

type FancyString struct {
	Content, Fmt string
}

// FancyString slice support routines
type FancyStrSlice []FancyString

func (a FancyStrSlice) Len() int           { return len(a) }
func (a FancyStrSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FancyStrSlice) Less(i, j int) bool { return a[i].Content < a[j].Content }

func PlainColumn(strs FancyStrSlice) string {
	result := ""
	for _, s := range strs {
		result = result + s.Content + "\n"
	}
	return result
}

func (strs FancyStrSlice) List() string {
	var results []string
	for _, s := range strs {
		results = append(results, s.Print("%s"))
	}
	return strings.Join(results, ", ")
}

func (s FancyString) Print(strfmt string) string {
	return fmt.Sprintf("%s"+strfmt+"%s", s.Fmt, s.Content, Reset())
}

func (strs FancyStrSlice) Columnate() string {
	nstrs := len(strs)
	result := ""
	if nstrs > 0 {
		if config.Fancy {
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
			var cols int
			if maxwidth > config.ScreenWidth {
				cols = 1
			} else {
				cols = config.ScreenWidth / maxwidth
			}
			stride := (nstrs + cols - 1) / cols
			fmt.Fprintf(os.Stderr, "Fmt: %q  Screen: (%d,%d)  Cols: %d  Stride: %d\n", strfmt, 0, config.ScreenWidth, cols, stride)
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
