# Proof of work implementation

This project implement a proof of work algorithm based on same idea as Bitcoin's proof of work and Hashcash.
The algorithm generates a random puzzle and then tries to find a nonce that when combined with the puzzle
produces a hash with a certain number of leading zeros. The number of leading zeros is determined by the difficulty
level. The attacker can't predict the nonce that will produce the desired hash, so the only way to find it is by
brute force. On the other hand, the verifier can easily check if the nonce is valid by hashing the puzzle and the nonce
and checking if the hash has the required number of leading zeros. For hash function, it uses SHA256 because it is time proven and secure hash function which has high computational cost but at the same time easy to check. But another hash function with similar properties can be used as well.

# Protocol

For demonstration purposes was implemented TCP based request-response protocol.
The server and client communicate using encoded messages.
There are 5 types of messages: challenge request, challenge response, verify request, verified response, and error response.

Each message encoded as a byte slice. The first byte is the message type. The rest of the bytes are the message body.
All fields except nonce in the message bodies has fixed length so it's not necessary to include separator between fields.
You can see implementation of decoding and encoding of message bodies in `pkg/protocol/bodies.go`.

# Usage

## Tested on
* macOS 14.3.1
* docker 25.0.3
* docker-compose 2.24.5

## Run demo
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
make list
```