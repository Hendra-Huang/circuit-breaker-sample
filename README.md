# circuit-breaker-sample
Circuit breaker sample is a sample project for simulating the circuit breaker pattern. This project use prometheus for monitoring and grafana for visualization. This project is dockerized.

### Scenario
There are 3 services provided, which are Service A, Service B and Service C.
```
A -> B -> C
Service A will make a request to Service B and will wait for the response. Service B will make a request to Service C and will wait for the response. When Service C has finished processing the request, Service C will give response to Service B and Service B will give response to Service A.
```

### Get Started
1. Just run `docker-compose up -d`
2. Visit localhost:1111/test for trying a request
3. Visit localhost:3000 for the monitoring dashboard
