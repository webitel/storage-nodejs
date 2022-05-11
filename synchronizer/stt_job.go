package synchronizer

import (
	"fmt"

	"github.com/webitel/storage/app"
	"github.com/webitel/storage/model"
)

type SttJob struct {
	file model.SyncJob
	app  *app.App
}

func (s *SttJob) Execute() {
	fmt.Println("EXECUTE JOB ", s.file)
}
