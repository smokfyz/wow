# Proof of work implementation

This project implements a proof-of-work algorithm based on the same idea as Bitcoin's proof of work and Hashcash. The algorithm generates a random puzzle and then attempts to find a nonce that, when combined with the puzzle, produces a hash with a certain number of leading zeros. The number of leading zeros is determined by the difficulty level. The attacker cannot predict the nonce that will produce the desired hash, so the only way to find it is through brute force. Conversely, the verifier can easily check if the nonce is valid by hashing the puzzle and the nonce, and verifying if the hash has the required number of leading zeros. For the hash function, it employs SHA256 because it is a time-proven and secure hash function with a high computational cost, yet it is easy to verify. However, another hash function with similar properties can also be utilized. The implementation can be found here `pkg/pow/pow.go`

# Protocol

For demonstration purposes, a TCP-based request-response protocol has been implemented. The server and client communicate using encoded messages. There are five types of messages: challenge request, challenge response, verify request, verified response, and error response.

Each message is encoded as a byte slice. The first byte indicates the message type, while the remaining bytes constitute the message body. All fields in the message bodies, except for the nonce, have fixed lengths, so it is not necessary to include a separator between fields. You can find the implementation of decoding and encoding of message bodies in `pkg/protocol/bodies.go`.

# Usage

There is a client and a server implementation for the protocol. The client simulates the behavior of several clients that connect to the server and solve the proof-of-work challenge. The server listens for incoming connections, sends a challenge to the clients, and verifies the solutions.

## Tested on
* macOS 14.3.1
* docker 25.0.3
* docker-compose 2.24.5

## Run demo client and server
You can configure some parameters in `.env` file in the root of the project.
```
make compose
```

## Run tests
```
make test
```

## Run linter
```
make lint
```