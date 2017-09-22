// Copyright 2017 Kubernetes Community Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"html/template"
	"log"
	"time"

	"github.com/k8s-community/codegen/pkg/config"
	"github.com/k8s-community/codegen/pkg/logger"
	"github.com/k8s-community/codegen/pkg/router"
)

// Handler defines common part for all handlers
type Handler struct {
	logger      logger.Logger
	config      *config.Config
	maintenance bool
	stats       *stats
}

type stats struct {
	requests        *Requests
	averageDuration time.Duration
	maxDuration     time.Duration
	totalDuration   time.Duration
	requestsCount   time.Duration
	startTime       time.Time
}

// New returns new instance of the Handler
func New(logger logger.Logger, config *config.Config) *Handler {
	return &Handler{
		logger: logger,
		config: config,
		stats: &stats{
			requests:  new(Requests),
			startTime: time.Now(),
		},
	}
}

// Base handler implements middleware logic
func (h *Handler) Base(handle func(router.Control)) func(router.Control) {
	return func(c router.Control) {
		timer := time.Now()
		handle(c)
		h.countDuration(timer)
		h.collectCodes(c)
	}
}

// Root handler shows version
func (h *Handler) Root(c router.Control) {
	t, err := template.ParseFiles("templates/layout.html", "templates/generate.html")
	if err != nil {
		log.Fatalf("Cannot parse `generate.html` or `layout.html`: %s", err)
	}

	t.ExecuteTemplate(c, "layout", nil)
}

func (h *Handler) countDuration(timer time.Time) {
	if !timer.IsZero() {
		h.stats.requestsCount++
		took := time.Now()
		duration := took.Sub(timer)
		h.stats.totalDuration += duration
		if duration > h.stats.maxDuration {
			h.stats.maxDuration = duration
		}
		h.stats.averageDuration = h.stats.totalDuration / h.stats.requestsCount
		h.stats.requests.Duration.Max = h.stats.maxDuration.String()
		h.stats.requests.Duration.Average = h.stats.averageDuration.String()
	}
}

func (h *Handler) collectCodes(c router.Control) {
	if c.GetCode() >= 500 {
		h.stats.requests.Codes.C5xx++
	} else {
		if c.GetCode() >= 400 {
			h.stats.requests.Codes.C4xx++
		} else {
			if c.GetCode() >= 200 && c.GetCode() < 300 {
				h.stats.requests.Codes.C2xx++
			}
		}
	}
}
