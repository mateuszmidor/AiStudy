#!/usr/bin/env bash

# use python 3.7 and pip3.7 to be compatible with latest tensorflow 2.1.0 (as of 03.03.2020)
# it is safe to install python3.7 with: pamac build python37
PYTHON=python3.7
PIP=pip3.7

# install virtualenv, tensorflow, matplotlib
function setupEnv() {
  sudo $PIP install -U virtualenv  # system-wide install
  virtualenv --system-site-packages -p $PYTHON ./venv
  source ./venv/bin/activate # virtualenv
      $PIP install --user --upgrade pip
      $PIP install --user tensorflow 
      $PIP install --user matplotlib
      echo "TensorFlow test:"
      $PYTHON -c "import tensorflow as tf;print(tf.reduce_sum(tf.random.normal([1000, 1000])))"
  deactivate # virtualenv
}

# train the model to recognize hand-written digits 0..9
function trainModel() {
  source ./venv/bin/activate # virtualenv
    $PYTHON Training.py
  deactivate  # virtualenv
}

# use the trained model to recognize a digit from random digit image
function runRecognize() {
  source ./venv/bin/activate # virtualenv
    $PYTHON Inference.py
  deactivate  # virtualenv
}


[[ ! -d venv ]] && setupEnv
[[ ! -d trained_model ]] && trainModel
[[ -d trained_model ]] && runRecognize