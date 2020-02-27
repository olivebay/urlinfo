# urlinfo-api

The urlinfo API is a service written in Go that lets your client applications check if a URL is included on any of the urlinfo lists.

The repository is configured with TravisCI to test the application code and then push the docker image to Docker Hub.


### Features

- It relies on MongoDB for data storage running in a container  
- MongoDB is initialised with sample data used to test the API
- [Gorilla Mux](https://github.com/gorilla/mux) is used for routing   


## Endpoints

#### GET /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}

```
curl localhost:9090/urlinfo/1/<hostname_and_port>/<original_path_and_query_string>
```

For example: 

```
curl  localhost:9090/urlinfo/1/domain.com
```

Upon successful match a JSON response is returned similar to the following:

```
{
  "url": "domain.com",
  "domain": "domain.com",
  "positives": true,
  "total": 2,
  "blacklists": [
    "vault",
    "spamHaus"
  ]
}
```

If no results match your request, you will get an 404 with a `url not found`

```
{
  "message": "url not found"
}
```

# Deployment

The urlinfo API can be deployed:
- In the [local machine](https://github.com/olivebay/urlinfo/blob/master/README.md#localhost-deployment) using Docker Compose   
- In a [Kubernetes K3s cluster](https://github.com/olivebay/urlinfo/blob/master/README.md#kubernetes-deployment)   


## Local machine Deployment with Docker Compose

**Requirements**

- Docker 
- Docker Compose 

**Get the repository**

```
git clone git@github.com:olivebay/urlinfo.git
cd urlinfo/
```

Starting services:

```
docker-compose up -d
```

Stoping services:

```
docker-compose stop
```

To test the API run the following:

```
curl localhost:9090/urlinfo/1/domain.com
```

## Kubernetes Deployment on K3s

**Requirements**

- A clean installation of a [Vagrant](https://www.vagrantup.com/downloads.html) box built with `ubuntu/xenial64`  
- [VirtualBox](https://www.virtualbox.org/wiki/Downloads)  


**Get the repository**

```
git clone git@github.com:olivebay/urlinfo.git
cd urlinfo/
```

**Install K3s and deploy the urlinfo API**

There is a `Makefile` that simplifies the cluster and application deployment.

To install K3s and deploy the urlinfo application on the K3s cluster:

```
make all
```

To test the API endpoint run: 

```
make test
```

To run a load test against the API:

```
make loadstart
```

To stop the load test:

```
make loadstop
```

Finally, to remove K3s and all auto-generated files:

```
make destroy
```