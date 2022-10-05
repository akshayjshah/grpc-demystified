gRPC Demystified
================

https://user-images.githubusercontent.com/972790/194163802-3cb7cdfe-144f-4646-a630-87ae595a518b.mp4

This repository contains the slides and code for a lightning talk I hope to
give at Gophercon 2022. In the talk, we build a gRPC server &mdash; from
scratch &mdash; using just the Go standard library.

The slides are available in [Keynote](grpc-demystified.key) or
[PDF](grpc-demystified.pdf) format. There's also a five-minute
[recording](grpc-demystified.mp4) of me practicing the talk.

The code includes a [REST handler](rest.go) and a from-scratch [gRPC
handler](grpc.go), both implementing the same logic. There's also a client for
each, along with a `grpc-go` client to show that our handler is speaking the
wire protocol correctly. To start the HTTP server and make a request with each
client, `go run .`.

If this talk appeals to you, the [Connect](https://connect.build) RPC framework
may be right up your alley.
