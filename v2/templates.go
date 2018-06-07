package v2

import "net/textproto"

func Exists(value string) bool {
	return value != ""
}

func Missing(value string) bool {
	return value == ""
}

func ParseRequestValues(rawValues map[string]string) (TemplatedValues, map[string]string) {
	t := make(TemplatedValues)
	remainder := make(map[string]string)

	for k, v := range rawValues {
		switch v {
		case "${exists}":
			t[textproto.CanonicalMIMEHeaderKey(k)] = Exists
		case "${missing}":
			t[textproto.CanonicalMIMEHeaderKey(k)] = Missing
		default:
			remainder[k] = v
		}
	}
	return t, remainder
}
