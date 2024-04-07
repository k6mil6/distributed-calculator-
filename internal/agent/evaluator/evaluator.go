package evaluator

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"log/slog"
	"reflect"
	"time"
)

type Result struct {
	Id     int
	Result float64
}

func Evaluate(subexpression model.Subexpression, heartbeat time.Duration, ch chan int, workerId int, logger *slog.Logger) (Result, error) {
	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-ticker.C:
				ch <- workerId
			case <-done:
				return
			}
		}
	}()

	timer := time.NewTimer(time.Duration(subexpression.Timeout) * time.Second)

	<-timer.C

	ticker.Stop()

	logger.Info("expression timeout has gone, starting evaluation", slog.Int("worker_id", workerId), slog.Any("expression", subexpression.Subexpression), slog.Any("timeout", time.Duration(subexpression.Timeout)*time.Second))
	exp, err := govaluate.NewEvaluableExpression(subexpression.Subexpression)
	if err != nil {
		return Result{}, err
	}

	result, err := exp.Evaluate(nil)
	if err != nil {
		return Result{}, err
	}

	logger.Info("expression evaluated", slog.Int("worker_id", workerId), slog.Any("result", result))

	resFloat, err := getFloat(result, reflect.TypeOf(float64(0)))
	if err != nil {
		return Result{}, err
	}

	logger.Info("expression evaluated to float", slog.Int("worker_id", workerId), slog.Any("result", resFloat))

	// signal that the evaluation is complete.
	done <- struct{}{}

	return Result{Id: subexpression.ID, Result: resFloat}, nil
}

func getFloat(unk interface{}, floatType reflect.Type) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}
