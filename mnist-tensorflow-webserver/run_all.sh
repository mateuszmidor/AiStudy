#!/usr/bin/env bash

DockerImage="python_tensor_matplotlib"
HttpPort=8000
trap stopRecognitionServer SIGINT

# build docker image with all necessary libs
function buildDockerImage() {
  docker build -t $DockerImage . 
}

# train the model to recognize hand-written digits 0..9
function trainRecognitionModel() {
  docker run -v `pwd`:/home -w /home $DockerImage python training.py
}

# run recognition server at localhost:$HttpPort
function runRecognitionServer() {
  docker run -v `pwd`:/home -w /home -p $HttpPort:$HttpPort $DockerImage python main.py $HttpPort &
  sleep 3
  firefox localhost:$HttpPort
  echo "Digit recognition server is running at localhost:$HttpPort. CTRL+C to exit..."
  while true; do sleep 1; done # wait for CTRL + C in a loop
}

# stop the server
function stopRecognitionServer() {
  container_id=`docker ps | grep $DockerImage | cut -d ' ' -f1`
  docker stop $container_id
  exit 0
}


# "main()"  :)
[[ `docker images | grep $DockerImage` == "" ]]  && buildDockerImage
[[ ! -d trained_model ]] && trainRecognitionModel
[[ -d trained_model ]] && runRecognitionServer