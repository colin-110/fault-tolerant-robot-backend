package workers

import (
	"context"
	"log"
	"time"

	"github.com/colin-110/fault-tolerant-robot-backend/internal/storage"
)

func StartFailureDetector(
	ctx context.Context,
	liveness *storage.LivenessStore,
	commandStore *storage.CommandStore,
) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				dead, err := liveness.Expired(3 * time.Second)
				if err != nil {
					log.Println("failure detector error:", err)
					continue
				}

				for _, robotID := range dead {
					log.Printf("robot %s considered DEAD", robotID)

					if err := commandStore.FailInFlightCommands(robotID); err != nil {
						log.Println("failed to mark commands FAILED:", err)
					}
				}
			}
		}
	}()
}
