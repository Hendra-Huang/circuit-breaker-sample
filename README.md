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

### Experiment
There will be 2 experiments, one is using context without circuit breaker and the other one is using circuit breaker. I use prometheus for monitoring and grafana for visualization. I start a bench test using apache bench with total 5000 requests and 100 concurrent requests. I measure the average response time of each service in one minute.

![without_circuit_breaker](https://user-images.githubusercontent.com/3291928/29350057-1657c6b0-8287-11e7-942e-39d422db6068.png)
This is the average response time without circuit breaker

![with_circuit_breaker](https://user-images.githubusercontent.com/3291928/29352829-690202d0-8292-11e7-9f33-a306c938356f.png)
This is the average response time with circuit breaker

Circuit breaker is keeping the response time stable, but the drawback is there are more requests that failed because of the timeout and the circuit breaker is open.
