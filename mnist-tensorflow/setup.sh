#!/usr/bin/env bash

sudo pacman -Syy
sudo pacman -Syu
sudo pacman -S python3-dev
sudo pacman -S python-pip

sudo pip3 install -U virtualenv  # system-wide install
virtualenv --system-site-packages -p python3 ./venv # create local virtualenv
source ./venv/bin/activate  # enable the virtualenv
pip install --upgrade pip
pip install tensorflow
pip install matplotlib
pip list  # show packages installed within the virtual environment
python -c "import tensorflow as tf;print(tf.reduce_sum(tf.random.normal([1000, 1000])))"
# deactivate  # don't exit until you're done using TensorFlow