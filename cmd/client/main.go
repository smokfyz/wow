package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smokfyz/wow/pkg/pow"
	"github.com/smokfyz/wow/pkg/protocol"
)

const (
	initialDifficulty = 1
)

func simulateClient(ctx context.Context, host, port string, pause int, wg *sync.WaitGroup) {
	defer wg.Done()

	challenge, err := pow.NewChallenge(initialDifficulty)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create challenge")
	}

	client := protocol.NewClient(challenge)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to server")
	}
	defer conn.Close()

	buf := make([]byte, protocol.MaxBodySize)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Request challenge
			_, err := conn.Write(client.GetChallengeRequest())
			if err != nil {
				log.Error().Err(err).Msg("failed to write challenge request")
				return
			}

			// Handle challenge response
			n, err := conn.Read(buf)
			if err != nil {
				log.Error().Err(err).Msg("failed to read challenge response")
				return
			}
			if client.IsErrorResponse(buf[:n]) {
				err = client.HandleErrorResponse(buf[:n])
				if err != nil {
					log.Error().Err(err).Msg("failed to handle error response")
				}
				return
			}
			resp, err := client.HandleChallengeResponse(buf[:n])
			if err != nil {
				log.Error().Err(err).Msg("failed to handle challenge response")
				return
			}
			_, err = conn.Write(resp)
			if err != nil {
				log.Error().Err(err).Msg("failed to write verify request")
				return
			}

			// Handle verified response
			n, err = conn.Read(buf)
			if err != nil {
				log.Error().Err(err).Msg("failed to read verified response")
				return
			}
			if client.IsErrorResponse(buf[:n]) {
				err = client.HandleErrorResponse(buf[:n])
				if err != nil {
					log.Error().Err(err).Msg("failed to handle error response")
				}
				return
			}
			err = client.HandleVerifiedResponse(buf[:n])
			if err != nil {
				log.Error().Err(err).Msg("failed to handle verified response")
				return
			}

			time.Sleep(time.Duration(pause) * time.Second)
		}
	}
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("failed to load .env file")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse log level")
		} else {
			zerolog.SetGlobalLevel(level)
		}
	}

	timeBetweenRequests, err := strconv.Atoi(os.Getenv("TIME_BETWEEN_REQUESTS_IN_SECONDS"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse time between requests")
	}

	numberOfSimultaneousRequests, err := strconv.Atoi(os.Getenv("NUMBER_OF_SIMULTANEOUS_REQUESTS"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse number of simultaneous requests")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = ""
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}

	for range numberOfSimultaneousRequests {
		wg.Add(1)
		go simulateClient(ctx, host, port, timeBetweenRequests, &wg)
	}

	wg.Wait()
}
