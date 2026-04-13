package utilx

import "strings"

func StringsJoin(sep string, elements ...string) string {
	return strings.Join(elements, sep)
}
