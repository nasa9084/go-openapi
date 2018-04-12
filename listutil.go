package openapi

func car(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	return ss[0]
}

func cdr(ss []string) []string {
	if len(ss) == 0 {
		return nil
	}
	return ss[1:]
}
