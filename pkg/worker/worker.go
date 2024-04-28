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

package worker

import (
	"net/http"

	"github.com/astaxie/beego"
	bcontext "github.com/astaxie/beego/context"
	"go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"

	"github.com/lizw91/otel-demo/global"
)

var tracer = otel.Tracer("otel-demo")

const serviceName = "worker"

func Init() error {
	global.InitLoggerFactory(serviceName)
	return nil
}

type helloController struct {
	beego.Controller
}

func (c *helloController) Hello() {
	ctx, span := tracer.Start(c.Ctx.Request.Context(), "sayHello")
	defer span.End()
	reqBaggage := baggage.FromContext(ctx)
	global.LoggerF.For(ctx).Info("get user name from baggage", zap.String("user.name", reqBaggage.Member("user.name").Value()))
	c.Ctx.WriteString("hello " + reqBaggage.Member("user.name").Value())
}

func Run() error {
	beego.Router("/hello", &helloController{}, "*:Hello")
	beego.InsertFilter("*", beego.BeforeRouter, func(c *bcontext.Context) {
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		newctx, span := tracer.Start(ctx, "getUserFromBaggage")
		defer span.End()

		reqBaggage := baggage.FromContext(newctx)
		span.SetAttributes(
			attribute.String("user.id", reqBaggage.Member("user.id").Value()),
			attribute.String("user.name", reqBaggage.Member("user.name").Value()),
		)
		c.Request = c.Request.WithContext(newctx)
	})
	beego.RunWithMiddleWares(":8081", otelbeego.NewOTelBeegoMiddleWare("worker", otelbeego.WithSpanNameFormatter(
		func(operation string, req *http.Request) string {
			return req.URL.Path
		},
	)))
	return nil
}
