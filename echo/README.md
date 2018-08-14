# note

## build
`$ docker build -t echo-server:latest .`

## run
`$ docker run -it --rm -p 4000:18888 echo-server`

## remove image
- no container uses -> `$ docker image prune`
- no tags attached -> `$ docker rmi $(docker images -f 'dangling=true' -q)`
