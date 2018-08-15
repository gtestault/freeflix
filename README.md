### Start with Docker
clone the repo and inside the source folder build the docker image

```docker build -t freeflix .```

then run a container and bind the port 8080 exposed from the image

```docker run -p 8080:8080 freeflix```

### Screenshots

![movies](/doc/screenshots/movies.png "Movies Dashboard")

![filter](/doc/screenshots/dialog.png "Advanced Filter")