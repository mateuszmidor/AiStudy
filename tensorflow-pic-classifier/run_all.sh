#!/usr/bin/env bash

PYTHON=python3
PIP=pip3

function stage() {
    BLUE="\e[36m"
    RESET="\e[0m"
    msg="$1"
    
    echo
    echo -e "$BLUE$msg$RESET"
}

function checkPrerequsites() {
    stage "Checking prerequisites"

    command $PYTHON --version > /dev/null 2>&1
    [[ $? != 0 ]] && echo "You need to install $PYTHON to run this example" && exit 1

    command $PIP --version > /dev/null 2>&1
    [[ $? != 0 ]] && echo "You need to install $PIP to run this example" && exit 1

    echo "OK"
}

function setupEnv() {
	stage "Setup python virtual env"

	sudo $PIP install -U virtualenv  # system-wide install
	virtualenv --system-site-packages -p $PYTHON ./venv
	source ./venv/bin/activate # virtualenv
		$PIP install --user --upgrade pip
		$PIP install opencv-python
		$PIP install --user tensorflow 
		echo "TensorFlow test:"
		$PYTHON -c "import tensorflow as tf;print(tf.reduce_sum(tf.random.normal([1000, 1000])))"
	deactivate # virtualenv

	echo "OK"
}

function getModel() {
	stage "Get inference model"

	wget https://bit.ly/2xf0fkV
	mkdir model/
	unzip 2xf0fkV -d model/
	rm 2xf0fkV

	echo "OK"
}

function runExample() {
    stage "Running example"
	
	source ./venv/bin/activate # virtualenv
		$PYTHON Main.py
	deactivate # virtualenv
	
	echo "Done"
}

checkPrerequsites
[[ ! -d venv ]] && setupEnv
[[ ! -d model ]] && getModel
runExample