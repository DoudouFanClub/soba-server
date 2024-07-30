# soba-server

soba-server is a web server designed to handle requests from a sofa-frontend as a REST API endpoint, and it doubles as a message buffer for soba-inferer as well.

The responsibilities of this server include, storing and caching of conversations for its users, queueing messages to be sent to the inferer and streaming the responses of the inferer back to the user.

## How to run
Ensure that the mongodb and redis services has already started, soba-server depends on mongodb for persistence and redis for caching. 

Simply call the soba-server executable with `endpoints.cfg` configured properly.

If you are building from source `go tidy` will install the dependencies for you, afterwards simply use `go run .`
```
{
    "endpoints_data": [
        {
            "ip": "127.0.0.1",
            "port": "7060"
        }
    ]
}
```
The config file expects to have both an ip and a port to the inferer that it will be connecting to, as such simply rerunning the application will be sufficient to increase or decrease the amount of endpoints that the server expects to be connected to.

## Dependencies
- go v1.22.3
- go-redis v9
- go.mongodb-org driver
- Redis (Note: If you are running this server on Windows, Redis requires WSL)
- MongoDB