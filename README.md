# Golang Words of Wisdom Proof of Work TCP Client Server

## Notes

This is just a proof of concept. The client and server are combined into a single repository for easier code reading and testing.

In a production-ready environment, of course, the code should be split into two separate repositories. 

## Why I chose POW and a Hashcash-style algorithm

1.	The configuration allows setting different difficulty levels for the algorithm (hash `length` and `difficulty` for leading zeros). Validation on the server side is very lightweight, while clients need to spend some time and computational power depending on the difficulty. This makes botnet and other attacks very resource-expensive.
2.	The implementation can be extended so that the server dynamically increases the `difficulty` and `length` based on the number of active connections.
3.	The Hashcash-style algorithm is simple and widely used, with implementations available for many programming languages.
4.	Protobuf was chosen as the transport protocol to simplify encoding and decoding for both the client and server. Protobuf implementations are available in many programming languages.
5.	A `solve_timeout` config parameter was introduced on the server side to disconnect clients that take too long to solve the challenge.

## Prepare Docker Images

```shell
docker compose build && docker compose pull
```

## Run server

```shell
docker compose up server
```

## Run client

```shell
docker compose up client
```
