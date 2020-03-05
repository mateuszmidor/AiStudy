from tensorflow import keras
import numpy as np 

def recognizeDigit(greyscale_28x28x8bit : np.ndarray) -> str:
    """ Input is numpy array 28x28 of type byte[0..255 range] """

    # prediction mechanism expects data in range 0..1
    input_image = greyscale_28x28x8bit / 255.0

    # tensorflow expects a series of images, so called "tensor". So we convert image into array of size 1x28x28
    input_image_for_tensorflow = np.expand_dims(input_image, 0)

    # read calculated earlier model
    model =  keras.models.load_model('trained_model')

    # recognize digit. Result is array [1, 10], numbers in range 0..1
    prediction_result = model.predict(input_image_for_tensorflow)
    predicted_label = np.argmax(prediction_result)

    print("Prediction results:")
    for i, f in enumerate(prediction_result[0]):
        print("  {} - {:02d}%".format(i, int(f*100)))

    return str(predicted_label) if prediction_result[0, predicted_label] > 0.5 else '?'