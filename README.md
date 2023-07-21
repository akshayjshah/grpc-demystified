gRPC Demystified
================

[![Watch the video](https://img.youtube.com/vi/rNI_pCa9slQ/0.jpg)](https://youtu.be/rNI_pCa9slQ)

This repository contains the slides and code for a lightning talk I gave at
Gophercon 2022. In the talk, we build a gRPC server &mdash; from scratch
&mdash; using just the Go standard library.

If you'd like a copy of the slides, they're available in
[Keynote](grpc-demystified.key) or [PDF](grpc-demystified.pdf) format. You can
watch the talk on [YouTube](https://youtu.be/rNI_pCa9slQ) (or click the image
above).

The code includes a [REST handler](rest.go) and a from-scratch [gRPC
handler](grpc.go), both implementing the same logic. There's also a client for
each, along with a `grpc-go` client to show that our handler is speaking the
wire protocol correctly. To start the HTTP server and make a request with each
client, `go run .`.

If this talk appeals to you, the [Connect](https://connect.build) RPC framework
may be right up your alley.
