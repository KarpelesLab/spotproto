# Spot protocol

Each client establish a connection to one or more spot servers. Each connection allows sending & receiving packets that are encapsulated and as such already have a known length.

As such, each packet only starts with a single byte identifying the packet type & version, and then the actual data depending on the packet type/version.

4 bits version (0), 4 bits type (0~15)

Packet types:

* 0x0 (C→S) Ping / (S→C) pong
* 0x1 (S→C) Handshake request, including to-be connection ID & server info in a cbor map: `['srv':'srvcode','cid':'srvcode:clientid','rnd':'<random data>']`
* 0x1 (C→S) Handshake response, cbor map of: `['id':'<signed id card>','sig':'<signature of handshake initial packet>']`. Upon receiving this packet, if the signature is valid, the ID card is registered and connection established.
* 0x2 (C→S & S→C) Instant message in instant message format (see below)

Each connection has an anonymous name (srvcode:clientid, where clientid is a random printable string), and can also be identified by the sha256 of idcard.Self. If multiple connections are made, messages sent to idcard.Self will be randomly distributed.

## Initial flow

* Upon connection, the server sends Handshake Start
* The client responds with Handshake Response
* If the provided IDCard isn't up to date in terms of groups, the server may send a new handshare request with `['grp':[...]]` set. The client must update its ID and try again.
* The server sends HandshakeRequest with Ready=true

## Instant message

Messages can be sent host to host. These are instant single messages (A→B), and can be dropped if the host isn't found or the message is lost in transit (rare). Each sent message can have a return address. Instant messages can have a body size of up to 65535 bytes.

An address has the form:

    target type:target/endpoint

For example:

    k:j0NDRmSPa5bfid2pAcUXaxCm2Dlh3TwayItZstwyeqQ/eth
    k:srvcode:j0NDRmSPa5bfid2pAcUXaxCm2Dlh3TwayItZstwyeqQ/eth (srvcode can optionally be added)
    c:srvcode:clientid/eth (srvcode is required when using clientid)
    g:j0NDRmSPa5bfid2pAcUXaxCm2Dlh3TwayItZstwyeqQ/api:json (send to a random nearby member of this group)

Endpoint types:

* 00: named endpoint
* 01: response endpoint

Message structure: the message is a structure of the following format:

* message id (16 bytes)
* flags (varint)
* len+recipient address
* len+sender address
* body (byte array)

The following flags are defined:

* 1: `MSG_NOTBOTTLE`: body is not an encrypted bottle. Normally messages must be encrypted for recipient and signed by sender using cryptutil.Bottle, however some protocols may skip this for improved efficiency, such as "eth".

### Well known instant message endpoints

#### eth

TODO update this.

body is a network frame, typically ipv4 or ipv6, typically up to 1500 bytes but allowed up to 65k. If a host receives something on the eth endpoint it doesn't need to respond to it, but to forward the frame to the local tuntap device if any is in use. If not, the packet can be ignored.

The `eth` protocol uses `MSG_NOTBOTTLE` flag. Instead the body is encrypted using the following method: a ECDSA key is generated and can be used to encrypt multiple packets. Each packet has the following format: public key, IV,

