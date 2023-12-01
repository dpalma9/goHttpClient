# Overview
Golang HTTP TSL Client.

# Run

Build:

```bash
$ docker build -t gohttpclientibmcloud -f Dockerfile .
```
And run:

```bash
$ docker run -d --rm --name goclient gohttpclientibmcloud
```

# Dev with VSCode

.devcontainer.json is a file that VSCode will read if you install the Dev Containers extensions.
This file is used to pre-configure the VSCode with all the plugins you need on your dev container.

# Develop (docker)

In order to build your docker test image in go, you need to run it with the following command the first time:

```Dockerfile
# Run just the first time; then comment the line:
RUN go mod init myclient && go mod tidy
```

Run your test container:

```bash
$ docker run -d --rm --name goclient gohttpclientibmcloud
```

If you want to test your container using a localhost service, in order to the container was able to connect with localhost, you'll need to run the container on the host network:

```bash
$ docker run -d --rm --name goclient --network host gohttpclientibmcloud
```