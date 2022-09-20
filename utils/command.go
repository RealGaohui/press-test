package utils

type Shell interface {
	Command(cmd string) (string, error)
}
