# go-docker-cloud
========

Go library for Dockercloud's API. Full documentation available at [https://docs.docker.com/apidocs/docker-cloud/](https://docs.docker.com/apidocs/docker-cloud/)

##Set up

**Installation:**

In order to install the Dockercloud Go library, you can use :

	go get github.com/docker/go-dockercloud/dockercloud


**Auth:**

In order to be able to make requests to the API, you should first obtain an ApiKey for your account. For this, log into Dockercloud, click on the menu on the upper right corner of the screen and select **Get Api Key**.

You can use your ApiKey with the Go library in any of the following ways:

Manually set in your Go code

	dockercloud.User = "yourUsernameHere"
	dockercloud.ApiKey = "yourApiKeyHere"

Use the docker.cfg file

Set the environment variables DOCKERCLOUD_USER and DOCKERCLOUD_APIKEY

**Namespace:**

In order to access the objects of a specific organization, you need to first set the **Namespace**. As for the authentication there are 2 ways of doing this:

- Manually setting the Namespace in the Go code:

	dockercloud.Namespace = "yourOrganizationNamespace"

- Set the environment variable DOCKERCLOUD_NAMESPACE

##Examples


**Note**

Each of the methods that require a uuid number as argument can also use a resource_uri as argument.

**Creating and deploying a NodeCluster**

```
nodecluster, err := dockercloud.CreateNodeCluster(dockercloud.NodeCreateRequest{Name: "Go-SDK-test", Region: "/api/v1/region/digitalocean/lon1/", NodeType: "/api/v1/nodetype/digitalocean/1gb/", Target_num_nodes: 2})

if err != nil {
  log.Println(err)
}

if err = nodecluster.Deploy(); err != nil {
   log.Println(err)
}
```

**Creating and starting a Stack**

```
stack, err := dockercloud.CreateStack(dockercloud.StackCreateRequest{Name: "new-stack", Services: []dockercloud.ServiceCreateRequest{{Image: "dockercloud/hello-world", Name: "test", Target_num_containers: 2}}})

if err != nil {
  log.Println(err)
}

if err = stack.Start(); err != nil {
   log.Println(err)
}
```

**Listing running containers**

```
containers, err := dockercloud.ListContainers()

if err != nil {
	log.Println(err)
}

log.Println(containers)
```

**Stopping a running service**

```
service, err := dockercloud.GetService("7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")
// or service, err := dockercloud.GetService("/api/v1/service/7eaf7fff-882c-4f3d-9a8f-a22317ac00ce")


if err != nil {
	log.Println(err)
}

if err = service.StopService(); err != nil {
   log.Println(err)
}
```


**Events**

Initialize event stream with namespace and Filters:
```
    myNamespace := dockercloud.NewNamespace("mynamespace")
    myFilter := dockercloud.NewStreamFilter(&dockercloud.EventFilter{Type: "container"})

    stream := dockercloud.NewStream(myNamespace)
    stream := dockercloud.NewStream(myFilter)
    stream := dockercloud.NewStream(myNamespace, myFilter)
```

Usage:
```
func OnMessage(event *dockercloud.Event) {
    log.Printf("On Message: %+v: ", event)
}

func OnError(err error) {
    log.Printf("On Error: %+v: ", err)
}

func OnConnect(namespace string) {
    log.Printf("On Connect: Stream %s", namespace)
}

func OnClose(namespace string) {
    log.Printf("On Close: Stream %s", namespace)
}

func main() {
    stream := dockercloud.NewStream()
    stream.OnError(OnError)
    stream.OnMessage(OnMessage)
    // stream.OnConnect(OnConnect)
    // stream.OnClose(OnClose)

    go func() {
        time.Sleep(10 * time.Second)
        stream.Close()
    }()

    if err := stream.Connect(); err == nil {
        stream.RunForever()
    } else {
        log.Print("Connect err: " + err.Error())
    }
```

Alternatively, you can use channels to handle messages and errors
```
    stream := dockercloud.NewStream()
    err := stream.Connect()

    if err := stream.Connect(); err == nil {
    	go stream.RunForever()
    } else {
        log.Print("Connect err: " + err.Error())
    }

	for {
		select {
			case msg := <- stream.MessageChan:
                log.Printf("%+v", msg)
			case err := <- stream.ErrorChan:
			    log.Printf("%+v", err)
		}
	}
}
```


The complete API Documentation is available [here](https://docs.docker.com/apidocs/docker-cloud/) with additional examples written in Go.
