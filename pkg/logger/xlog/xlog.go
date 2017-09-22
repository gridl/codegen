// Copyright 2017 Kubernetes Community Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logger

import (
	"os"

	"github.com/k8s-community/codegen/pkg/logger"
	"github.com/rs/xlog"
)

// newXLog creates "github.com/rs/xlog" logger
func newXLog(config *logger.Config) logger.Logger {
	var out xlog.Output
	switch config.Err {
	// We should find more matches between types of output
	case nil, os.Stderr:
		out = xlog.NewConsoleOutput()
	default:
		out = xlog.NewConsoleOutput()
	}
	return xlog.New(xlog.Config{
		Level:  xlog.Level(config.Level),
		Fields: config.Fields,
		Output: out,
	})
}
