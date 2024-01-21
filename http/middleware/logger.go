package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"order-backend/http/handler"
	"order-backend/pkg/customerror"
	"order-backend/pkg/logger"
	"strconv"
	"time"
)

func Logger() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		startAt := time.Now()
		preHandleRequest(c)

		defer func() {
			postHandleRequest(c, startAt)
		}()

		return c.Next()
	}
}

func preHandleRequest(c *fiber.Ctx) {
	log := logger.Log()
	header := make(map[string]string)
	c.Request().Header.VisitAll(func(key, val []byte) {
		k := bytes.NewBuffer(key).String()
		header[k] = bytes.NewBuffer(val).String()
	})

	requestID := c.Request().Header.Peek("X-Request-ID")
	if requestID == nil {
		requestUUID, _ := uuid.NewUUID()
		requestID = []byte(requestUUID.String())
	}

	headerByte, _ := json.Marshal(header)

	log = log.With(
		zap.String("request_id", string(requestID)),
		zap.String("path", c.Path()),
		zap.Any("header", json.RawMessage(headerByte)),
		zap.String("method", c.Method()),
		zap.String("protocol", c.Protocol()),
		zap.String("remote_ip", c.IP()),
	)
	c.Locals("request_id", string(requestID))
	c.Locals("path", c.Path())

	if len(c.Body()) != 0 {
		log = log.With(zap.Any("request", json.RawMessage(c.Body())))
	}

	if (len(c.Request().URI().QueryString())) != 0 {
		log = log.With(zap.Any("query_param", c.Request().URI().QueryString()))
	}

	msg := fmt.Sprintf("[REQUEST] %v %v", c.Method(), c.Path())
	log.Info(msg)
}

func postHandleRequest(c *fiber.Ctx, startAt time.Time) {
	log := logger.Log()

	rvr := recover()
	if rvr != nil {
		var ok bool
		err, ok := rvr.(error)
		if !ok {
			err = fmt.Errorf("%v", rvr)
		}

		err = handler.ErrorResponse(c, fiber.StatusInternalServerError, customerror.Error(customerror.ErrCodeGeneric, err.Error()))
		if err != nil {
			log.Error(err.Error())
		}
	}

	var requestId string
	if val, ok := c.Locals("request_id").(string); ok {
		requestId = val
	}

	fullURL := fmt.Sprintf("%s://%s%s", c.Protocol(), c.Hostname(), c.Path())

	log = log.With(
		zap.String("request_id", requestId),
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
		zap.String("protocol", c.Protocol()),
		zap.String("remote_ip", c.IP()),
		zap.Any("status_code", c.Response().StatusCode()),
		zap.Any("response", json.RawMessage(c.Response().Body())),
		zap.Float64("latency", time.Since(startAt).Seconds()),
	)

	msg := fmt.Sprintf("[RESPONSE] %d %s %s", c.Response().StatusCode(), c.Method(), fullURL)
	switch strconv.Itoa(c.Response().StatusCode())[0] {
	case '1', '2', '3':
		log.Info(msg)
	case '4', '5':
		log.Error(msg)
	default:
		log.Panic(msg)
	}
}
