package server

type MasterServer interface {
	Start() error
	Stop() error
}
