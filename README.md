# shipping-api
```bash
# Start the server
make run

# Book shipping
curl -i -X POST http://localhost:8080/api/shipping --data '{"origin":"SE","destination":"SE","weight":400}'
> HTTP/1.1 201 Created
> Content-Type: application/json; charset=utf-8
> Location: /api/shipping/11f713e5-f826-4264-8481-19fb69331cde
> Date: Tue, 24 Jan 2023 21:22:43 GMT
> Content-Length: 45

> {"id":"11f713e5-f826-4264-8481-19fb69331cde"}

# Get info about booking
curl -i http://localhost:8080/api/shipping/11f713e5-f826-4264-8481-19fb69331cde
> HTTP/1.1 200 OK
> Content-Type: application/json; charset=utf-8
> Date: Tue, 24 Jan 2023 21:24:22 GMT
> Content-Length: 121

> {"id":"11f713e5-f826-4264-8481-19fb69331cde","origin":"SE","destination":"SE","weight":400,"price":2000,"currency":"SEK"}
```

# API
`[POST] /api/shipping/` - book shipping  
`[GET] /api/shipping/:id` - get booking information by id

# Deployment
## 1000 monthly users
I would deploy this as simply as possible in the very early stages as we don't yet know the full extent of the service, odds are a lot of stuff will change rapidly.  
E.g. as a container with Firebase, EC2 instance, Heroku, or similar, focus being on rapid development and iteration.  
With 1000 monthly active users we shouldn't see scaling issues, so 2 instances for redundancy with a load balancer in front should be sufficient.  

## 10.000.000 monthly users
With that many users we hopefully have a stable API and we know the service should not change too much as we already have seen success.  
There are multiple decisions to make:

### Microservice vs monolith
Microservices _can_ be easier to scale as we can spot-scale single bottleneck services.  
It can also make it easier to work with a higher number of collaborators than a single codebase/unit of deployment.  
However microservices come with a cost, it's more difficult to monitor, trace, and maintenance may increase since every service needs it's own pipeline, packaging, tests etc.  
If microservices are chosen I would most likely deploy them in a Kubernetes cluster, as it's incredibly good at orchestrating multiple services, and a single unified API can be used.  

### Caches
At this time caching for read requests can really help speed up response times and offload the services.  

### Adding read replicas for the database
Data is more often read than written, so we can separate into read and write replicas of the DB if we are using e.g. Postgres.  
If using a non-relational DB such as MongoDB we can scale it horizontally by adding more replicas.  
