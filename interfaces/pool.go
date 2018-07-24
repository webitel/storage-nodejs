package interfaces

type TaskInterface interface {
	Execute()
}

type PoolInterface interface {
	Resize(n int)
	Close()
	Wait()
	Exec(task TaskInterface)
}
