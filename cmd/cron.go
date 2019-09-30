package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"gitlab.quvideo.com/algo/algo-resume-tool/common"
	"gitlab.quvideo.com/algo/algo-resume-tool/orm"
)

func runCron() {
	common.Log.Info("startup cron task... ")

	err := orm.NewDB(false, viper.GetBool("database.debug"), viper.GetString("database.addr"))
	if err != nil {
		common.Log.Error("open db error")
		return
	}
	common.Log.Info("db open")

	c := cron.New()
	defer c.Stop()
	cronOn := viper.GetString("tool.checkInterval")
	if cronOn == "" {
		cronOn = "@every 1h"
	}
	if _, err := c.AddJob(cronOn, query{}); err != nil {
		common.Log.Error("starting cronOn job err: ", err)
		return
	}

	c.Start()
	select {}
}

type query struct{}

func (t query) Run() {
	common.Log.Debug("in every 1min query task")

	task := &orm.JimuTask{}
	tasks := []orm.JimuTask{}

	var ct time.Time
	for {
		var tx *gorm.DB
		if AppName != "" {
			common.Log.Debug("query app name: ", AppName)
			tx = orm.DB.Model(task).Order("create_at").Where("is_callback_retry = true AND create_at > ? AND app_name = ?", ct, AppName).Limit(viper.GetUint("tool.fetch")).Find(&tasks)
		} else {
			tx = orm.DB.Model(task).Order("create_at").Where("is_callback_retry = true AND create_at > ?", ct).Limit(viper.GetUint("tool.fetch")).Find(&tasks)
		}
		if tx.Error != nil {
			common.Log.Errorf("db query failed")
			break
		}
		for _, v := range tasks {
			common.Log.Debugf("task_id: %d, callback_addr: %s, create_at: %v\n", v.TaskId, v.CallbackAddr, v.CreateAt)
			go t.rePost(v)
		}
		if len(tasks) == 0 {
			break
		}
		ct = tasks[len(tasks)-1].CreateAt
		fmt.Printf("last create_at: %v", ct)
		tasks = tasks[:0]
	}
}

func (t query) rePost(task orm.JimuTask) {
	content := task.AlgoResponseBuffer
	var rep common.AlgoVideoPipelineCallbackResponse

	httpReq, err := http.NewRequest("POST", task.CallbackAddr, bytes.NewBuffer(content))
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	httpRep, err := client.Do(httpReq)
	if err != nil {
		common.Log.Errorf("task %d http.Do() failed: %v", task.TaskId, err)
		return
	}
	defer httpRep.Body.Close()
	if httpRep.StatusCode != 200 {
		common.Log.Errorf("task %d HttpCallback:%s StatusCode:%d not success", task.TaskId, task.CallbackAddr, httpRep.StatusCode)
		return
	}
	httpBody, err := ioutil.ReadAll(httpRep.Body)
	if err != nil {
		common.Log.Errorf("task %d ReadAll() failed: %v", task.TaskId, err)
		return
	}

	common.Log.Debugf("VideoPipelineCallback httpReq.Body:%s httpRep.Body: %s", string(content), string(httpBody))
	err = json.Unmarshal(httpBody, &rep)
	if err != nil {
		common.Log.Errorf("task id %d json.Unmashal() failed: %v", task.TaskId, err)
		t.updateDB(&task, orm.CallbackFailed, false)
		return
	}
	if rep.ErrorCode != common.API_CALLBACK_OK {
		common.Log.Errorf("task %d HttpCallback:%s not API_CALLBACK_OK ErrorCode:%s ErrorMsg:%s", task.TaskId, task.CallbackAddr, rep.ErrorCode, rep.ErrorMsg)
		t.updateDB(&task, orm.CallbackFailed, false)
		return
	}
	t.updateDB(&task, orm.CallbackSuccessed, false)
	common.Log.Info("task %d success reposted", task.TaskId)
}

func (t query) updateDB(task *orm.JimuTask, status orm.CallbackStatus, retry bool) {
	tx := orm.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		common.Log.Error("db begin error")
		return
	}

	db := orm.DB.Model(task).Where("task_id = ?", task.TaskId).Updates(map[string]interface{}{
		"is_callback_retry": retry,
		"callback_status":   status,
	})

	if db.Error != nil {
		common.Log.Error("db updates error")
		tx.Rollback()
		return
	}
	if tx.Commit().Error != nil {
		common.Log.Errorf("db commit error")
		tx.Rollback()
		return
	}
}
