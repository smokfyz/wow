package protocol

import "github.com/rs/zerolog/log"

type Verifier interface {
	SetDifficulty(difficulty byte) error
	Solve(puzzle []byte) []byte
}

type Client struct {
	verifier Verifier
}

func NewClient(verifier Verifier) Client {
	return Client{verifier}
}

func (c Client) GetChallengeRequest() []byte {
	return ChallengeRequestBody{}.Encode()
}

func (c Client) HandleChallengeResponse(rawBody []byte) ([]byte, error) {
	body, err := DecodeChallengeResponseBody(rawBody)
	if err != nil {
		log.Error().Msgf("failed to decode challenge response: %s", err)
		return nil, err
	}

	err = c.verifier.SetDifficulty(body.difficulty)
	if err != nil {
		log.Error().Msgf("failed to set difficulty: %s", err)
		return nil, err
	}

	nonce := c.verifier.Solve(body.puzzle)
	log.Debug().Interface("puzzle", body.puzzle).Msgf("solved puzzle: %x", nonce)
	return VerifyRequestBody{body.puzzle, nonce}.Encode(), nil
}

func (c Client) HandleVerifiedResponse(rawBody []byte) error {
	body, err := DecodeVerifiedResponseBody(rawBody)
	if err != nil {
		log.Error().Msgf("failed to decode verified response: %s", err)
		return err
	}
	log.Info().Msgf("server verified nonce: %s", body.result)
	return nil
}

func (c Client) HandleErrorResponse(rawBody []byte) error {
	body, err := DecodeErrorResponseBody(rawBody)
	if err != nil {
		log.Error().Msgf("failed to decode error response: %s", err)
		return err
	}
	log.Info().Msgf("server responded with an error: %s", body.err)
	return nil
}

func (c Client) IsErrorResponse(rawBody []byte) bool {
	return rawBody[0] == errorResponse
}
