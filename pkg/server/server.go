/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	driver "github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/lizw91/otel-demo/global"
)

type ServerOptions struct {
	WorkerAddr string

	MysqlAddr     string
	MysqlUserName string
	MysqlPassword string
	MysqlDBName   string
}

var (
	db     *gorm.DB
	tracer = otel.Tracer("otel-demo")
	Opts   = &ServerOptions{}
)

const serviceName = "server"

func Init() error {
	global.InitLoggerFactory(serviceName)

	cfg := &driver.Config{
		User:                 Opts.MysqlUserName,
		Passwd:               Opts.MysqlPassword,
		Net:                  "tcp",
		Addr:                 Opts.MysqlAddr,
		DBName:               Opts.MysqlDBName,
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
	}
	dsn := cfg.FormatDSN()
	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}
	return db.Use(otelgorm.NewPlugin(otelgorm.WithoutMetrics()))
}

func Run() error {
	router := gin.New()
	router.Use(
		otelgin.Middleware("server", otelgin.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/healthz"
		})),
	)
	router.GET("/users/:id", func(c *gin.Context) {
		ctx := c.Request.Context()
		id := c.Param("id")
		global.LoggerF.For(ctx).Info("request param", zap.String("id", id))
		name, err := getUser(ctx, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		echoResp, err := doHelloRequest(ctx, Opts.WorkerAddr, id, name)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, echoResp)
	})
	return router.Run(":8080")
}

func doHelloRequest(ctx context.Context, workerAddr, id, name string) (string, error) {
	userBaggage, err := baggage.Parse(fmt.Sprintf("user.id=%s,user.name=%s", id, name))
	if err != nil {
		otel.Handle(err)
	}

	req, err := http.NewRequestWithContext(baggage.ContextWithBaggage(ctx, userBaggage), http.MethodGet, workerAddr+"/hello", nil)
	if err != nil {
		return "", err
	}
	otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}
	bts, _ := io.ReadAll(resp.Body)
	return string(bts), nil
}

// get user name by user id
func getUser(ctx context.Context, id string) (string, error) {
	// start a new span from context.
	newCtx, span := tracer.Start(ctx, "getUser", trace.WithAttributes(attribute.String("user.id", id)))
	defer span.End()
	// add start event
	span.AddEvent("start to get user",
		trace.WithTimestamp(time.Now()),
	)
	var username string
	// get user name from db, if you want to trace it, `WithContext` is necessary.
	result := db.WithContext(newCtx).Raw(`select username from users where id = ?`, id).Scan(&username)
	if result.Error != nil || result.RowsAffected == 0 {
		err := fmt.Errorf("user %s not found", id)
		global.LoggerF.For(newCtx).Error(err.Error())
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	// set user info in span's attributes
	span.SetAttributes(attribute.String("user.name", username))
	// add end event
	span.AddEvent("end to get user",
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(attribute.String("user.name", username)),
	)
	span.SetStatus(codes.Ok, "")
	return username, nil
}
