package main

import (
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	cors "github.com/itsjamie/gin-cors"
	"github.com/rfashwal/scs-room/internal"
	"github.com/rfashwal/scs-room/internal/config"
	"github.com/rfashwal/scs-room/internal/dbprovider"
	httptransport "github.com/rfashwal/scs-room/internal/transport/http"
	"github.com/rfashwal/scs-room/internal/transport/mq"
	"github.com/rfashwal/scs-utilities/rabbit"
)

func main() {
	dbManager, err := dbprovider.NewDBManager()
	if err != nil {
		log.Fatal("error initializing DB Manager", err)
	}

	conf := config.Config().Manager
	mqManager, err := rabbit.NewRabbitMQManager(conf.RabbitURL())
	if err != nil {
		log.Fatalf("MQ server init: %s", err)
	}

	pub, err := mqManager.InitPublisher()
	if err != nil {
		log.Fatalf("MQ.Publisher init: %s", err)
	}

	err = pub.RabbitConnector.DeclareTopicExchange(conf.TemperatureTopic())
	if err != nil {
		log.Fatalf("MQ.Publisher.DeclareTopicExchange : %s", err)
	}

	svc, err := internal.NewService(pub, conf, dbManager)
	if err != nil {
		log.Fatal("service init err", err)
	}

	go mq.TemperatureObserver(svc, mqManager, conf)

	router, err := httptransport.NewServer(svc)
	if err != nil {
		log.Fatalf("http server init: %s", err)
	}

	manager := config.EurekaManagerInit()
	manager.SendRegistrationOrFail()
	manager.ScheduleHeartBeat(conf.ServiceName(), 10)
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	err = router.Run(conf.Address())

	if err != nil {
		log.Fatalf("router run: %s", err)
	}
}
