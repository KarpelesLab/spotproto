# Spot protocol

Each client establish a connection to one or more spot servers. Each connection allows sending & receiving packets that are encapsulated and as such already have a known length.

As such, each packet only starts with a single byte identifying the packet type & version, and then the actual data depending on the packet type/version.

4 bits version (0), 4 bits type (0~15)

Packet types:

* 0x0 (C→S) Ping / (S→C) pong
* 0x1 (S→C) Handshake request, CBOR map with fields: `srv` (server code), `cid` (client ID as `srvcode.clientid`), `rnd` (random nonce), `grp` (optional group list), `rdy` (optional, true when handshake is complete)
* 0x1 (C→S) Handshake response, CBOR map with fields: `id` (optional client identifier), `key` (PKIX-encoded public key), `sig` (signature over the handshake request). Upon receiving this packet, if the signature is valid, the key is registered and connection established.
* 0x2 (C→S & S→C) Instant message in instant message format (see below)

Each connection has an anonymous name (`srvcode.clientid`, where clientid is a random printable string), and can also be identified by the SHA-256 hash of the public key. If multiple connections are made, messages sent to the key-based ID will be randomly distributed.

## Initial flow

* Upon connection, the server sends Handshake Start
* The client responds with Handshake Response
* If the provided key isn't up to date in terms of groups, the server may send a new handshake request with the `grp` field set. The client must update its ID and try again.
* The server sends HandshakeRequest with Ready=true

## Instant message

Messages can be sent host to host. These are instant single messages (A→B), and can be dropped if the host isn't found or the message is lost in transit (rare). Each sent message can have a return address. Instant messages can have a body size of up to 65535 bytes.

An address has the form:

    type.target/endpoint
    type.srvcode.target/endpoint

For example:

    k.j0NDRmSPa5bfid2pAcUXaxCm2Dlh3TwayItZstwyeqQ/eth
    k.srvcode.j0NDRmSPa5bfid2pAcUXaxCm2Dlh3TwayItZstwyeqQ/eth (srvcode can optionally be added)
    c.srvcode.clientid/eth (srvcode is required when using connection-based ID)
    g.j0NDRmSPa5bfid2pAcUXaxCm2Dlh3TwayItZstwyeqQ/api (send to a random nearby member of this group)

Endpoint types:

* 00: named endpoint
* 01: response endpoint

Message structure: the message is a structure of the following format:

* message id (16 bytes)
* flags (varint)
* len+recipient address
* len+sender address
* body (byte array)

The following flags are defined (bit positions):

* 1 (bit 0): `MsgFlagResponse` - This is a response message that must not trigger further responses
* 2 (bit 1): `MsgFlagError` - The message body contains an error string
* 4 (bit 2): `MsgFlagNotBottle` - Body is not an encrypted bottle. Normally messages must be encrypted for recipient and signed by sender using cryptutil.Bottle, however some protocols may skip this for improved efficiency or when already encrypted by another mechanism

### Well known instant message endpoints

#### eth

The `eth` endpoint is used for network frame tunneling (virtual ethernet). The body is a network frame, typically IPv4 or IPv6, up to 1500 bytes (MTU) but allowed up to 65535 bytes. If a host receives something on the eth endpoint, it should forward the frame to the local tuntap device if one is in use. If no tuntap device is configured, the packet can be ignored.

The `eth` protocol uses the `MsgFlagNotBottle` flag since the payload is encrypted using a separate key exchange mechanism rather than the standard Bottle encryption.
