package string

import "strings"

func LastTrimAndRetain(origin string, cut string) string {
	index := strings.LastIndex(origin, cut)
	if index == -1 {
		return origin
	}

	return origin[index+len(cut):]
}
