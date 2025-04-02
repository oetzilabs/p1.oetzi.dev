package models

type Footer struct {
	Commands []FooterCommand
}

type FooterCommand struct {
	Key   string
	Value string
}
