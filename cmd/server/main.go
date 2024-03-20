package main

import (
	"context"
	"errors"
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

func handleConnection(ctx context.Context, conn net.Conn, proto protocol.Server, wg *sync.WaitGroup, timeout time.Duration) {
	defer conn.Close()
	defer wg.Done()

	buf := make([]byte, protocol.MaxBodySize)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
				log.Warn().Err(err).Msg("failed to set read deadline")
				return
			}
			n, err := conn.Read(buf)
			if err != nil {
				log.Warn().Err(err).Msg("failed to read request")
				return
			}
			resp := proto.HandleRequest(buf[:n])
			if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
				log.Warn().Err(err).Msg("failed to set write deadline")
				return
			}
			if _, err := conn.Write(resp); err != nil {
				log.Warn().Err(err).Msg("failed to write response")
				return
			}
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

	difficulty, err := strconv.Atoi(os.Getenv("PUZZLE_DIFFICULTY"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse difficulty")
	}

	challenge, err := pow.NewChallenge(byte(difficulty))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create challenge")
	}

	ttl, err := strconv.Atoi(os.Getenv("PUZZLE_TTL_IN_SECONDS"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse TTL")
	}

	server := protocol.NewServer(challenge, time.Duration(ttl)*time.Second)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		log.Info().Msg("shutting down server")
		ln.Close()
	}()

	log.Info().Str("port", port).Msg("starting server")

	timeout, err := strconv.Atoi(os.Getenv("CONNECTION_TIMEOUT_IN_SECONDS"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse timeout")
	}

	wg := sync.WaitGroup{}
	for {
		conn, err := ln.Accept()
		if errors.Is(err, net.ErrClosed) {
			wg.Wait()
			return
		}
		if err != nil {
			log.Warn().Err(err).Msg("failed to accept connection")
			continue
		}
		wg.Add(1)
		go handleConnection(ctx, conn, server, &wg, time.Duration(timeout)*time.Second)
	}
}
