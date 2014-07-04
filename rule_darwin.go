package main

// +build: darwin

func NewRules() []*Rule {
	return []*Rule{
		&Rule{
			Matcher: FileExt(".png"),
			Command: []string{"open", "-gF", "-a", "Preview", "$source"},
		},
		&Rule{
			Matcher: FileExt(".puml"),
			Command: []string{"make"},
		},
	}
}
