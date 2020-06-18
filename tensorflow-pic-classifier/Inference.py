from tensorflow import lite as tflite 
import numpy as np
from Common import opencv

class Inference(object):
    def __init__(self, model_file_path, labels_file_path):
        # load classification model
        self.load_model_and_configure(model_file_path)

        # load labels for model
        self.load_labels_from_file(labels_file_path)

    def load_model_and_configure(self, model_path):
        # create image interpreter
        self.interpreter = tflite.Interpreter(model_path)
        self.interpreter.allocate_tensors()

        # store input and output details
        self.input_details = self.interpreter.get_input_details()
        self.output_details = self.interpreter.get_output_details()

        # store input image dimmensions
        self.input_image_height = self.input_details[0]['shape'][1]
        self.input_image_width = self.input_details[0]['shape'][2]

    def load_labels_from_file(self, file_path):
        # load labels into array
        with open(file_path, 'r') as file:
            self.labels = [line.strip() for line in file.readlines()]

    def prepare_image(self, image):
        # convert image to BGR
        image = opencv.cvtColor(image, opencv.COLOR_BGR2RGB)

        # resize image to desired size
        new_size = (self.input_image_height, self.input_image_width)
        image = opencv.resize(image, new_size, interpolation=opencv.INTER_AREA)
        return image

    def label_image(self, image):
        # prepare image
        image = self.prepare_image(image)

        # add dummy dimmension
        input_data = np.expand_dims(image, axis=0)

        # set interpreter input
        self.interpreter.set_tensor(self.input_details[0]['index'], input_data)

        # interpret image
        self.interpreter.invoke()

        # get probability array
        inference_result = self.interpreter.get_tensor(self.output_details[0]['index'])

        # get index of highest probability
        top_one = inference_result.argmax()

        # return label for highest probability 
        return self.labels[top_one]



