package main

// +build: darwin

func NewRules() []*Rule {
	return []*Rule{
		&Rule{
			SrcExt:  FileExt(".png"),
			Command: []string{"open", "-gF", "-a", "Preview", "$source"},
		},
		&Rule{
			SrcExt:  FileExt(".puml"),
			Command: []string{"make"},
		},
	}
}
