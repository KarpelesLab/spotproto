[![GoDoc](https://godoc.org/github.com/KarpelesLab/spotproto?status.svg)](https://godoc.org/github.com/KarpelesLab/spotproto)

# spotproto

Protocol definitions for Spot, a real-time messaging protocol between clients and servers. Messages are encrypted end-to-end between clients, ensuring that servers only relay opaque payloads without access to message contents.

## Installation

```bash
go get github.com/KarpelesLab/spotproto
```

## Overview

The protocol uses a simple packet format where the first byte encodes both the protocol version (upper 4 bits) and the packet type (lower 4 bits).

### Packet Types

| ID | Type | Description |
|----|------|-------------|
| 0x0 | PingPong | Keep-alive packets, payload is echoed back |
| 0x1 | Handshake | Authentication handshake (request/response) |
| 0x2 | InstantMsg | Encrypted messages between clients |

### Handshake Flow

1. Server sends `HandshakeRequest` containing server info, client ID, and a nonce
2. Client responds with `HandshakeResponse` containing its public key and a signature over the request
3. Server verifies the signature and completes the handshake

### Message Format

Messages contain:
- **MessageID**: 16-byte unique identifier
- **Flags**: Control bits (response, error, encryption status)
- **Recipient/Sender**: Client IDs in format `type.server.target` or `type.target`
- **Body**: Message payload (encrypted end-to-end)

### End-to-End Encryption

Message bodies are encrypted between clients using their exchanged public keys. The server acts only as a relay and cannot decrypt message contents. The `MsgFlagNotBottle` flag indicates when a message is sent unencrypted (e.g., for system messages).

## License

See LICENSE file.
