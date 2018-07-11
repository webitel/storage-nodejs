package model

type FileBackendSettings struct {
	Id         int
	Domain     *string
	TypeId     int
	Name       string
	Parameters map[string]string
}
