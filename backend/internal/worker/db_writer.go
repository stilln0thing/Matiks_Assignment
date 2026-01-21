package worker

import (
	"context"
	"log"
	"time"

	"github.com/stilln0thing/matiks_leaderboard/internal/models"
	"github.com/stilln0thing/matiks_leaderboard/internal/repository"
)

type DBWriter struct {
	queue         chan models.RatingUpdate
	repo          *repository.UserRepository
	batchSize     int
	flushInterval time.Duration
}

func NewDBWriter(repo *repository.UserRepository, queueSize, batchSize int, flushInterval time.Duration) *DBWriter {
	return &DBWriter{
		queue:         make(chan models.RatingUpdate, queueSize),
		repo:          repo,
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}
}
func (w *DBWriter) Queue() chan<- models.RatingUpdate {
	return w.queue
}
func (w *DBWriter) Start(ctx context.Context) {
	batch := make([]models.RatingUpdate, 0, w.batchSize)
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()
	log.Printf("[DBWriter] Started - batch: %d, interval: %v", w.batchSize, w.flushInterval)
	for {
		select {
		case update := <-w.queue:
			batch = append(batch, update)
			if len(batch) >= w.batchSize {
				w.flush(ctx, batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.flush(ctx, batch)
				batch = batch[:0]
			}
		case <-ctx.Done():
			// Flush remaining on shutdown
			if len(batch) > 0 {
				w.flush(context.Background(), batch)
			}
			log.Println("[DBWriter] Stopped")
			return
		}
	}
}
func (w *DBWriter) flush(ctx context.Context, batch []models.RatingUpdate) {
	if len(batch) == 0 {
		return
	}
	start := time.Now()
	err := w.repo.BatchUpdateRatings(ctx, batch)
	if err != nil {
		log.Printf("[DBWriter] Error: %v", err)
		return
	}
	log.Printf("[DBWriter] Flushed %d updates in %v", len(batch), time.Since(start))
}
