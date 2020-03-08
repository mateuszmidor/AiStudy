from tensorflow import keras
import numpy as np 
import matplotlib.pyplot as plt 
import random

# fetch test data from MNIST
mnist_digits = keras.datasets.mnist
_, (test_images, test_labels) = mnist_digits.load_data()

# select an image and relevant label. Scale image to range 0..1 by dividing by 255.0
testImageIndex = random.randint(0, len(test_images))

test_image = test_images[testImageIndex]
plt.imshow(test_image, cmap=plt.cm.binary)
print(test_image.shape)
print(np.amin(test_image))
print(np.amax(test_image))
test_image = test_image / 255.0
print(np.amin(test_image))
print(np.amax(test_image))

actual_label = test_labels[testImageIndex]

# tensorflow expects a series of images, so called "tensor". So we convert image into array of size 1x28x28
test_image_for_tensorflow = np.expand_dims(test_image, 0)

# read calculated earlier model
input_path = 'trained_model'
model =  keras.models.load_model(input_path)

# recognize digit. Result is 10 numbers in range 0..1
prediction_result = model.predict(test_image_for_tensorflow)
predicted_label = np.argmax(prediction_result)

# display results
plt.figure()
plt.xticks([])
plt.yticks([])
plt.grid(False)
plt.imshow(test_image, cmap=plt.cm.binary)
plt.xlabel('Correct: {}, recognized: {}, img index: {}'.format(actual_label, predicted_label, testImageIndex), fontsize=20)
plt.show()

# plt.savefig('img_' + str(testImageIndex) + '.png')