# note

## build
`$ docker build -t chat:latest .`

## run
`$ docker run -it --rm -p 4000:8080 chat`

## remove image
- no container uses -> `$ docker image prune`
- no tags attached -> `$ docker rmi $(docker images -f 'dangling=true' -q)`
