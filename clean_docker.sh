#!/usr/bin/env sh

docker stop $(docker container ls --all -q)
docker rm $(docker container ls --all -q)
docker network rm $(docker network ls -q -f name=1)

echo "\e[32mCleaned up.\e[0m"
echo "\e[32mErrors in the output are completely fine.\e[0m"