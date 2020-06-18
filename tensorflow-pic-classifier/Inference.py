from tensorflow import lite as tflite 
import numpy as np
from Common import opencv

class Inference(object):
    def __init__(self, model_file_path, labels_file_path):
        self.load_model_and_configure(model_file_path)
        self.load_labels_from_file(labels_file_path)

    def load_model_and_configure(self, model_path):
        self.interpreter = tflite.Interpreter(model_path)
        self.interpreter.allocate_tensors()
        self.input_details = self.interpreter.get_input_details()
        self.output_details = self.interpreter.get_output_details()

        self.input_image_height = self.input_details[0]['shape'][1]
        self.input_image_width = self.input_details[0]['shape'][2]

    def load_labels_from_file(self, file_path):
        with open(file_path, 'r') as file:
            self.labels = [line.strip() for line in file.readlines()]

    def prepare_image(self, image):
        image = opencv.cvtColor(image, opencv.COLOR_BGR2RGB)
        new_size = (self.input_image_height, self.input_image_width)
        image = opencv.resize(image, new_size, interpolation=opencv.INTER_AREA)
        return image

    def label_image(self, image):
        image = self.prepare_image(image)
        input_data = np.expand_dims(image, axis=0)
        self.interpreter.set_tensor(self.input_details[0]['index'], input_data)
        self.interpreter.invoke()

        inference_result = self.interpreter.get_tensor(self.output_details[0]['index'])
        top_one = inference_result.argmax()
        return self.labels[top_one]
        


