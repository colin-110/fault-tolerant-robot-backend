package workers

import (
	"context"
	"log"
)

type TelemetryJob struct {
	RobotID string
	Payload []byte
}

type TelemetryWorkerPool struct {
	queue chan TelemetryJob
}

func NewTelemetryWorkerPool() *TelemetryWorkerPool {
	return &TelemetryWorkerPool{
		queue: make(chan TelemetryJob, TelemetryQueueSize),
	}
}

func (p *TelemetryWorkerPool) Start(ctx context.Context) {
	for i := 0; i < TelemetryWorkerCount; i++ {
		go p.worker(ctx, i)
	}
}

func (p *TelemetryWorkerPool) worker(ctx context.Context, id int) {
	log.Printf("telemetry worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("telemetry worker %d stopping", id)
			return
		case job := <-p.queue:
			p.process(job)
		}
	}
}

func (p *TelemetryWorkerPool) process(job TelemetryJob) {
	// placeholder
	log.Printf("telemetry from robot %s", job.RobotID)
}

func (p *TelemetryWorkerPool) Enqueue(job TelemetryJob) {
	select {
	case p.queue <- job:
		// accepted
	default:
		// dropped silently (metrics later)
	}
}
