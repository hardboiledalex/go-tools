package utils

import "strings"

const FillWidth = 90

type HeaderLevel int

const (
	HeaderLine_CENTER HeaderLevel = iota
	HeaderLine_L1
	HeaderLine_L2
	HeaderLine_L3
)

func HeaderLine(message string, level ...HeaderLevel) string {
	n := 0
	var iLevel HeaderLevel
	if len(level) == 0 {
		iLevel = HeaderLine_CENTER
	} else {
		iLevel = level[0]
	}
	switch iLevel {
	case HeaderLine_CENTER:
		n = (FillWidth - len(message) - 4) / 2
	case HeaderLine_L1:
		n = 10
		break
	case HeaderLine_L2:
		n = 5
		break
	case HeaderLine_L3:
		n = 3
		break
	default:
		break
	}
	if n < 0 {
		n = 0
	}
	var result strings.Builder
	headerSeparator := strings.Repeat("#", FillWidth)

	result.WriteString("#")
	result.WriteString(strings.Repeat(" ", n))
	if n > 0 {
		result.WriteString(" ")
	}
	result.WriteString(message)
	if result.Len()+1 < FillWidth {
		result.WriteString(" ")
		result.WriteString(strings.Repeat(" ", FillWidth-result.Len()-1))
		result.WriteString("#")
	}
	return headerSeparator + "\n" + result.String() + "\n" + headerSeparator
}
