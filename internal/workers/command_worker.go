package workers

import (
	"context"
	"log"

	"github.com/colin-110/fault-tolerant-robot-backend/internal/storage"
)

type CommandJob struct {
	CommandID string
}

type CommandWorkerPool struct {
	queue chan CommandJob
	store *storage.CommandStore
}

func NewCommandWorkerPool(store *storage.CommandStore) *CommandWorkerPool {
	return &CommandWorkerPool{
		queue: make(chan CommandJob, CommandQueueSize),
		store: store,
	}
}

func (p *CommandWorkerPool) Start(ctx context.Context) {
	for i := 0; i < CommandWorkerCount; i++ {
		go p.worker(ctx, i)
	}
}

func (p *CommandWorkerPool) worker(ctx context.Context, id int) {
	log.Printf("command worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("command worker %d stopping", id)
			return
		case job := <-p.queue:
			p.process(job)
		}
	}
}

func (p *CommandWorkerPool) process(job CommandJob) {
	// placeholder logic for now
	log.Printf("processing command %s", job.CommandID)
}

func (p *CommandWorkerPool) Enqueue(job CommandJob) bool {
	select {
	case p.queue <- job:
		return true
	default:
		// queue full → backpressure
		return false
	}
}
