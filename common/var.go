package common

import (
	"io"

	"go.uber.org/zap"
)

// Log object
var Log *zap.SugaredLogger

// log output target object
var Outer io.Writer

type AlgoVideoPipelineCallbackResponse struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

const (
	API_OK                  = "OK"
	API_FAILED              = "FAILED"
	API_ALGO_INTERNAL_FAILE = "ALGO_INTERNAL_FAILED"
	API_CALLBACK_OK         = "OK"
	API_CALLBACK_FAILED     = "FAILED"
)
