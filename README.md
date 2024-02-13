# Rate limiter in Go using Redis

This is a simple rate limiter middleware implementation in Go using Redis. It is based on the token bucket algorithm. The middleware can be used to limit the number of requests per second for a specific token or IP address.


## How it works 
The middleware has token preferance over IP. If the token is present in the request, the middleware will use the token to limit the requests. If the token is not present, the middleware will use the IP address to limit the requests. If the token is present and the IP is present, the middleware will use the token to limit the requests. 

Add the header "API_KEY" to the request to use the token to limit the requests. 

## Usage 

1. You can use the default env file or create a new one with the following variables:

```env
TOKEN_TIMEOUT=5
TOKEN_RPS=5
IP_TIMEOUT=5
IP_RPS=5
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DATABASE=0
REDIS_PASSWORD=""
USE_REDIS=1
```

* `TOKEN_TIMEOUT` - The time in seconds to wait for a token to request after reaching the limit
* `TOKEN_RPS` - The number of requests per second to reach the rate limit
* `IP_TIMEOUT` - The time in seconds to wait for a IP to request after reaching the limit
* `IP_RPS` - The number of requests per second to reach the rate limit
* `REDIS_HOST` - The Redis host
* `REDIS_PORT` - The Redis port
* `REDIS_DATABASE` - The Redis database
* `REDIS_PASSWORD` - The Redis password
* `USE_REDIS` - If set to 1, the middleware will use Redis to store the tokens and IPs, otherwise it will use an in-memory map

2. Create an http server with the rate limiter middleware:
```go
func main() {
	godotenv.Load(".env")
	ratelimiter := ratelimiter.NewRateLimiterMiddleware()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(ratelimiter)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	        w.Write([]byte("Hello, world!"))
	    })

	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
```

3. Run the redis server:
```bash
docker compose up -d 
```

4. Run the server (This will run 200 requests with 10 parallel jobs)
```bash
go run cmd/server/main.go
```

5. Run the server and test the rate limiter: (This will run 200 requests with 10 parallel jobs)
```bash
seq 1 200 | xargs -Iname -P10 curl http://localhost:8080/ -H "API_KEY: 123"
```


