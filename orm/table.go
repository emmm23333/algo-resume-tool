package orm

import (
	"time"
)

type TaskStatus int
type CallbackStatus int

const (
	TaskNone = iota
	TaskProcessing
	TaskFailed
	TaskSuccessed = TaskStatus(200)
)

const (
	CallbackNone CallbackStatus = iota
	CallbackFailed
	CallbackSuccessed = CallbackStatus(200)
)

type JimuTask struct {
	//TaskStatus     TaskStatus `gorm:"index:idx_task_code"`
	// Api
	//AppName        string `gorm:"type:varchar(64);index:idx_task_appname"`
	TaskId         uint64         `gorm:"AUTO_INCREMENT;primary_key"`
	TaskStatus     TaskStatus     // chenggong shibai
	CallbackStatus CallbackStatus // chenggong shibai
	//AlgoTaskId     uint64 `gorm:"unique_index:idx_task_id"`
	AlgoTaskId      uint64 `gorm:"index:idx_task_id"`
	JobId           string `gorm:"type:varchar(64)"`
	IsCallbackRetry bool
	CallbackNumber  int
	CallbackAddr    string `gorm:"type:varchar(256)"`
	AppName         string `gorm:"type:varchar(64);index:idx_task_appname"`
	ApiType         string `gorm:"type:varchar(64);index:idx_task_appname"`

	ProcessCostMs uint
	//AlgoProcessCode   string `gorm:"type:varchar(64);index:idx_task_code"`
	//AlgoProcessMsg    string `gorm:"type:varchar(1024)"`
	CreateAt time.Time
	UpdateAt time.Time

	AlgoRequestBuffer  []byte
	AlgoResponseBuffer []byte
}

func (t JimuTask) TableName() string {
	return "jimu_task"
}
