# llm-go

llm-go is a web server designed to handle requests from a client as a REST API endpoint.

This server handles load balancing between multiple llm instances, message queueing and data storing.

This server is meant to be used in the following architecture.

## Dependencies
- go v1.22.3
- go-redis v9
- go.mongodb-org driver
- rabbitmq amqp091-go driver
- Docker (if you choose to run it through the offical Docker image)

## Debugging Tools
- MongoDB Compass
- Postman

## Launch Options (In CLI)
- redis
  - linux-distribution-of-your-choice
  - redis-cli
- mongodb
  - mongod
- 