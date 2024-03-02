# Micronet
A microservice framework based on RPC protocol and Gob serializer that implements the observer pattern for real time notifiations.

## Client
A Client can send requests but cannot recieve any.
By default, every Client can ping a server.

## Server
A server can recieve requests but cannot send any.
By default, every Server registers a ping handler.

## ClientServer
A ClientServer can send and recieve requests from and to any other Client, Server or ClientServer.

## Pub/Sub
This framework implements the observer pattern, allowing you to configure a publish/subscribe communication between two microservices.

