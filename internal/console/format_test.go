package console

import (
	"os"
	"strconv"
	"testing"
)

func TestColumnate(t *testing.T) {
	cases := []struct {
		width int
		strs  []string
		want  string
	}{
		{70, []string{"a", "this is b", "c is for cookie", "d's good enough for me"}, `a                       c is for cookie         
this is b               d's good enough for me  
`},
		{160, []string{"a", "this is b", "c is for cookie", "d's good enough for me"}, `a                       this is b               c is for cookie         d's good enough for me  
`},
	}
	for _, c := range cases {
		os.Setenv("COLUMNS", strconv.Itoa(c.width))
		got := Columnate(c.strs)
		if got != c.want {
			t.Errorf("Columnate(%q)[%d] == %q, want %q", c.strs, c.width, got, c.want)
		}
	}
}
