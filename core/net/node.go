package net

import (
	"distgraphia/core/svc"
	"sync"
)

type Runner interface {
	// Run method is responsible for running the coordinator or worker.
	Run()
	GetRPCount()
}

type Coordinator struct {
	mu       sync.Mutex
	services map[string]*svc.Serviceable // Services, by names
	count    int                         // incoming RPCs
}

func (c *Coordinator) Run() {

}

func (c *Coordinator) GetRPCount() int {
	return c.count
}

type Worker struct {
	mu       sync.Mutex
	services map[string]*svc.Serviceable // Services, by names
	count    int                         // incoming RPCs
}

func (w *Worker) Run() {

}

func (w *Worker) GetRPCount() int {
	return w.count
}
