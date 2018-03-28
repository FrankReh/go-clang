package ast

type Sources interface {
	Extract(file string, soffset, eoffset int) (string, error)
}
