package main

import (
	"github.com/PonyWilliam/go-arealogs/handler"
	"github.com/PonyWilliam/go-arealogs/models"
	arealogs "github.com/PonyWilliam/go-arealogs/proto/arealogs"
	common "github.com/PonyWilliam/go-common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	"strconv"
	"time"
)
func main() {
	// New Service
	consulConfig,err := common.GetConsualConfig("1.116.62.214",8500,"/micro/config")
	//配置中心
	if err != nil{
		log.Fatal(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(
		func(options *registry.Options){
			options.Addrs = []string{"1.116.62.214"}
			options.Timeout = time.Second * 10
		},
	)
	srv := micro.NewService(
		micro.Name("go.micro.service.arealogs"),
		micro.Version("latest"),
		micro.Address(":8084"),
		micro.Registry(consulRegistry),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(common.QPS)),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	mysqlInfo := common.GetMysqlFromConsul(consulConfig,"mysql")
	db,err := gorm.Open("mysql",
		mysqlInfo.User+":"+mysqlInfo.Pwd+"@tcp("+mysqlInfo.Host + ":"+ strconv.FormatInt(mysqlInfo.Port,10) +")/"+mysqlInfo.DataBase+"?charset=utf8&parseTime=True&loc=Local",
	)

	if err != nil{
		log.Error(err)
	}
	defer db.Close()
	db.SingularTable(true)
	//初始化
	if !db.HasTable(&models.AreaLogs{}){
		db.CreateTable(&models.AreaLogs{})
	}

	// Initialise service
	srv.Init()

	// Register Handler
	_ = arealogs.RegisterAreaLogsHandler(srv.Server(), &handler.Arealogs{Db: db})
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
