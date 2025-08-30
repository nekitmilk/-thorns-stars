package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nekitmilk/agent/internal/collector/system"
	"github.com/nekitmilk/agent/internal/config"
	"github.com/nekitmilk/agent/internal/models"
	"github.com/nekitmilk/agent/internal/sender"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		log.Println("Received shutdown signal...")
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Fatalf("Agent failed: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Валидация конфигурации
	if cfg.HostID == "" {
		return fmt.Errorf("HOST_ID environment variable is required")
	}
	if cfg.PollingInterval < 10*time.Second {
		log.Printf("Warning: Polling interval too low, setting to 10s")
		cfg.PollingInterval = 10 * time.Second
	}

	log.Printf("Starting agent for host: %s", cfg.HostID)
	log.Printf("Monitoring center URL: %s", cfg.MonitoringCenterURL)
	log.Printf("Polling interval: %v", cfg.PollingInterval)

	systemCollector := system.NewSystemCollector()
	metricSender := sender.NewHTTPSender(cfg.MonitoringCenterURL, cfg.RequestTimeout)

	ticker := time.NewTicker(cfg.PollingInterval)
	defer ticker.Stop()

	// Первый сбор
	safeCollectAndSend(systemCollector, metricSender, cfg.HostID)

	// Основной цикл
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down agent gracefully...")
			return nil
		case <-ticker.C:
			go safeCollectAndSend(systemCollector, metricSender, cfg.HostID)
		}
	}
}

func safeCollectAndSend(collector *system.SystemCollector, sender *sender.HTTPSender, hostID string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	if err := collectAndSend(collector, sender, hostID); err != nil {
		log.Printf("Collection failed: %v", err)
	}
}

func collectAndSend(collector *system.SystemCollector, sender *sender.HTTPSender, hostID string) error {
	metrics, err := collector.Collect()
	if err != nil {
		return fmt.Errorf("failed to collect metrics: %w", err)
	}

	batch := models.MetricsRequest{
		HostID:    hostID,
		Metrics:   metrics,
		Timestamp: time.Now(),
	}

	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := sender.SendMetrics(batch); err != nil {
			if attempt == maxRetries {
				return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
			}

			backoff := time.Duration(attempt*attempt) * time.Second
			log.Printf("Attempt %d failed: %v, retrying in %v", attempt, err, backoff)
			time.Sleep(backoff)
			continue
		}

		log.Printf("Successfully sent %d metrics", len(metrics))
		return nil
	}

	return nil
}
