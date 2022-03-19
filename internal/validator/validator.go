package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if Errors map is empty
func (v Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds error if the error is not exists in the Errors map
func (v *Validator) AddError(key, msg string) {
	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = msg
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, msg string) {
	if !ok {
		v.AddError(key, msg)
	}
}

// In returns true if a specific value is in a list of strings.
func In(value string, list ...string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all string values in a slice are unique.
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, v := range values {
		uniqueValues[v] = true
	}

	return len(values) == len(uniqueValues)
}
