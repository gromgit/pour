package console

const (
	CSI            = "\033["
	RESET          = "0"
	BOLD_ON        = "1"
	BOLD_OFF       = "22"
	DIM_ON         = "2"
	DIM_OFF        = BOLD_OFF
	UNDERSCORE_ON  = "4"
	UNDERSCORE_OFF = "24"
	BLINK_ON       = "5"
	BLINK_OFF      = "25"
	REVERSE_ON     = "7"
	REVERSE_OFF    = "27"
)

var (
	Bold       = Set(BOLD_ON)
	Dim        = Set(DIM_ON)
	Underscore = Set(UNDERSCORE_ON)
	Off        = Reset()
)

func Set(attrs ...string) string {
	if len(attrs) > 0 {
		result := CSI + attrs[0]
		for _, attr := range attrs[1:] {
			result = result + ";" + attr
		}
		return result + "m"
	} else {
		return ""
	}
}

func Reset() string {
	return CSI + RESET + "m"
}
