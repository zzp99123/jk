// 手动在业务中打点
package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"goFoundation/webook/internal/service/sms"
)

type OpentelemetryOtel struct {
	svc    sms.Service
	tracer trace.Tracer
}

func NewopentelemetryOtel(svc sms.Service) *OpentelemetryOtel {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("webook/internal/service/sms/opentelemetry")
	return &OpentelemetryOtel{
		svc:    svc,
		tracer: tracer,
	}
}

// 还是用装饰器的方式以手动的方式打点
func (o *OpentelemetryOtel) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	ctx, span := o.tracer.Start(ctx, "sms_send"+biz,
		// 因为我是一个调用短信服务商的客户端
		trace.WithSpanKind(trace.SpanKindClient))
	defer span.End(trace.WithStackTrace(true))
	err := o.svc.Send(ctx, biz, args, numbers...)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
