package protocol

import (
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/rs/zerolog/log"
)

type Challenger interface {
	GetDifficulty() byte
	GeneratePuzzle() []byte
	Verify(puzzle []byte, nonce []byte) bool
}

type Encoder interface {
	Encode() []byte
}

// Server is a simple server that handles challenge and verify requests.
// It uses a cache to store puzzles with a TTL. Clients must solve the puzzle within the TTL.
// If the puzzle is not solved within the TTL, the server will remove the puzzle from the cache and
// the client will have to request a new puzzle. In case of correct verification, the server will
// return a random wisdom.
type Server struct {
	challenge Challenger
	cache     *ttlcache.Cache
}

func NewServer(challenge Challenger, ttl time.Duration) Server {
	cache := ttlcache.NewCache()
	cache.SetTTL(ttl)
	return Server{challenge, cache}
}

func (s Server) HandleRequest(rawBody []byte) []byte {
	var resp Encoder

	if len(rawBody) == 0 {
		panic("empty body")
	}

	switch rawBody[0] {
	case challengeRequest:
		log.Debug().Msg("received challenge request")
		resp = s.handleChallengeRequest()
	case verifyRequest:
		body, err := DecodeVerifyRequestBody(rawBody)
		if err == nil {
			log.Debug().Interface("puzzle", body.puzzle).Interface("nonce", body.nonce).Msg("received verify request")
			resp = s.handleVerifyRequest(body)
		} else {
			log.Debug().Err(err).Msg("failed to decode verify request")
			resp = ErrorResponseBody{err.Error()}
		}
	default:
		log.Debug().Msg("unexpected body type")
		resp = ErrorResponseBody{"unexpected body type"}
	}
	return resp.Encode()
}

func (s Server) handleChallengeRequest() Encoder {
	puzzle := s.challenge.GeneratePuzzle()
	s.cache.Set(string(puzzle), true)
	log.Debug().Interface("puzzle", puzzle).Msg("generated puzzle")
	return ChallengeResponseBody{puzzle, s.challenge.GetDifficulty()}
}

func (s Server) handleVerifyRequest(body VerifyRequestBody) Encoder {
	if _, ok := s.cache.Get(string(body.puzzle)); !ok {
		log.Debug().Interface("puzzle", body.puzzle).Msg("puzzle not found")
		return ErrorResponseBody{"puzzle not found"}
	}

	if !s.challenge.Verify(body.puzzle, body.nonce) {
		log.Debug().Interface("puzzle", body.puzzle).Interface("nonce", body.nonce).Msg("verification failed")
		return ErrorResponseBody{"verification failed"}
	}

	log.Debug().Interface("puzzle", body.puzzle).Interface("nonce", body.nonce).Msg("verification succeeded")
	s.cache.Remove(string(body.puzzle))
	return VerifiedResponseBody{getRandomWisdom()}
}
