package random

type IRandom interface {
	GenerateAlphaNum(length int) string
}
