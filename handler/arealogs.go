package handler

import (
	"context"
	"github.com/PonyWilliam/go-arealogs/models"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"

	arealogs "github.com/PonyWilliam/go-arealogs/proto/arealogs"
)


type Arealogs struct{
	Db *gorm.DB
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Arealogs) AddLog(ctx context.Context,req *arealogs.ALog,rsp *arealogs.Status) error{
	timeunix := time.Now().Unix()
	if err := e.Db.Create(&models.AreaLogs{PID: req.PID,WID: req.WID,Content: req.Content,Time: strconv.FormatInt(timeunix,10)}).Error;err!=nil{
		return err
	}
	rsp.Result = true
	rsp.Response = "success"
	return nil
}
func (e *Arealogs) FindAll(ctx context.Context, req *arealogs.Null, rsp *arealogs.Logs) error {
	var res []*models.AreaLogs
	if err := e.Db.Find(&res).Error;err!=nil{
		return err
	}
	for _,v := range res{
		temp := Swap(v)
		rsp.Logs = append(rsp.Logs,temp)
	}
	return nil
}
func (e *Arealogs) FindByID(ctx context.Context, req *arealogs.Id, rsp *arealogs.Log) error {
	res := &models.AreaLogs{}
	if err := e.Db.First(&res).Where("id = ?",req.Id).Error;err!=nil{
		return err
	}
	rsp = Swap(res)
	return nil
}
func (e *Arealogs) FindByWID(ctx context.Context, req *arealogs.Worker, rsp *arealogs.Logs) error {
	var res []*models.AreaLogs
	if err := e.Db.Find(&res).Where("wid = ?",req.Id).Error;err!=nil{
		return err
	}
	for _,v := range res{
		temp := Swap(v)
		rsp.Logs = append(rsp.Logs,temp)
	}
	return nil
}
func (e *Arealogs) FindByAID(ctx context.Context, req *arealogs.Area, rsp *arealogs.Logs) error {
	var res []*models.AreaLogs
	if err := e.Db.Find(&res).Where("aid = ?",req.Aid).Error;err!=nil{
		return err
	}
	for _,v := range res{
		temp := Swap(v)
		rsp.Logs = append(rsp.Logs,temp)
	}
	return nil
}
func Swap(req *models.AreaLogs) *arealogs.Log{
	rsp := &arealogs.Log{}
	rsp.ID = req.ID
	rsp.WID = req.WID
	rsp.Time = req.Time
	rsp.Content = req.Content
	rsp.AreaID = req.AreaID
	rsp.PID = req.PID
	return rsp
}