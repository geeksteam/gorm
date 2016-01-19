package gorm_test

import (
	"fmt"
	"log"
	"time"

	"github.com/geeksteam/gorm"
)

func ExampleGormUsing() {
	gorm.Run(NewRoutine(func() {
		for i := 0; i <= 2; i++ {
			log.Println("Couner: ", i)
			time.Sleep(time.Second)
		}
		log.Println("Finish")
	}, "task", "Simple counting", 0))

	time.Sleep(4 * time.Second)

	log.Printf("%+v\n", gorm.Routines())

	time.Sleep(1 * time.Second)
}

type Routine struct {
	Started     time.Time // Дата и время с микросекундами (27.07.2015 14:01.123) запуска рутины
	Stopped     time.Time // Дата и время остановки рутины
	Type        string    // Тип горутины ( возможные типы: web handler, task, background, backup )
	Owner       int       // Id пользователя запустившего рутину
	Description string    // Возможное описание рутины ( eg.: Backup files from /tmp/user )
	Status      int       // Статус рутины (  0 - successful done, 1 - running, 2 - stopped with error )
	ExitError   string    // Описание ошибки при завершении рутины с ошибкой
	f           func()    // Тело рутины
}

func NewRoutine(f func(), Type, description string, owner int) *Routine {
	return &Routine{
		Type:        Type,
		Owner:       owner,
		Description: description,
		f:           f,
	}
}

func (r *Routine) Start() {
	r.Started = time.Now()
	r.Status = 1
	defer func() {
		r.Stopped = time.Now()
	}()

	r.f()

	r.Status = 0
}

func (r *Routine) OnError(err error) {
	r.Status = 2
	r.ExitError = err.Error()
}

func (r *Routine) String() string {
	switch r.Status {
	case 0:
		return fmt.Sprintf("Started: %v\nStopped: %v\nType: %v\nDescription: %v\nOwnerID: %v\nStatus: Successful done", r.Started, r.Stopped, r.Type, r.Description, r.Owner)
	case 1:
		return fmt.Sprintf("Started: %v\nType: %v\nDescription: %v\nOwnerID: %v\nStatus: Running", r.Started, r.Type, r.Description, r.Owner)
	case 2:
		return fmt.Sprintf("Started: %v\nStopped: %v\nType: %v\nDescription: %v\nOwnerID: %v\nStatus: Stopped with error '%v'", r.Started, r.Stopped, r.Type, r.Description, r.Owner, r.ExitError)
	}
	return fmt.Sprintf("Started: %v\nType: %v\nDescription: %v\nOwnerID: %v\nStatus: Unknown", r.Started, r.Type, r.Description, r.Owner)
}

func (r *Routine) Finished() bool {
	if r.Status == 1 {
		return false
	}

	return true
}
