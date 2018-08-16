### Project

Streaming server with an integrated BitTorrent client. The server can stream over http while downloading a torrent. 

### Start with Docker
clone the repo and inside the source folder build the docker image.

```
git clone git@github.com:ninjaintrouble/freeflix.git
cd freeflix
docker build -t freeflix .
```

then run a container and bind the port 8080 exposed from the image

```docker run -p 8080:8080 freeflix```

### VPN 

Please configure your VPN before deploying. The BitTorrent client has no blocklist

### Screenshots

![movies](/doc/screenshots/movies.png "Movies Dashboard")

![filter](/doc/screenshots/dialog.png "Advanced Filter")