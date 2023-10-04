package forms

// customer errors type for holding validation errors from forms
// of type key=string; value=slice of strings
type errors map[string][]string

// method for error type to add error messages for a given field to the map
func (e errors) Add(field string, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	} else {
		return es[0]
	}
}
