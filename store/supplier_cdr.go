package store

type CdrSupplier struct {
	DatabaseLayer LayeredStoreDatabaseLayer
	ElasticLayer  *ElasticSupplier
}

func (self *CdrSupplier) Search() StoreChannel {
	return Do(func(result *StoreResult) {

	})
}
func (self *CdrSupplier) Scroll() StoreChannel {
	return Do(func(result *StoreResult) {

	})
}
