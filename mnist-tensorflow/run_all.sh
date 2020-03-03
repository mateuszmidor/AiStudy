#!/usr/bin/env bash

source ./venv/bin/activate  # enable the virtualenv
python Training.py
python Inference.py
deactivate  # don't exit until you're done using TensorFlow