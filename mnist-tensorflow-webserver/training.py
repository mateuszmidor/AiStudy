from tensorflow import keras
import numpy as np
import matplotlib.pyplot as plt

# fetch the mnist numbers model
mnist_digits = keras.datasets.mnist
(train_images, train_labels), _ = mnist_digits.load_data()

class_names = np.arange(10)

# show training data
# plt.figure()
# rowCount = 10
# colCount = 30
# for i in range(rowCount * colCount):
#     # draw single digit
#     plt.subplot(rowCount, colCount, i+1) 

#     # disable axis labels
#     plt.xticks([])
#     plt.yticks([])
#     plt.grid(False)

#     # print digits and their labels
#     plt.imshow(train_images[i], cmap=plt.cm.binary)
#     plt.xlabel(class_names[train_labels[i]])

# plt.show()

# convert greyscale 0..255 to float 0.0 .. 1.0
train_images = train_images / 255.0

# build neural network
model = keras.Sequential([
    # transform image 28x28 into 784 elements vector
    keras.layers.Flatten(input_shape=(28,28)),

    # hidden neural layer with 64 neurons and activation function tanh
    # keras.layers.Dense(64, activation='tanh'),

    # hidden neural layer with 128 neurons and activation function sigmoid
    keras.layers.Dense(128, activation='sigmoid'),

    # outer layer with 10 neurons and activation function softmax
    keras.layers.Dense(len(class_names), activation='softmax')
])

# setup learing process
model.compile(optimizer='adam', loss='sparse_categorical_crossentropy', metrics=['accuracy'])

# print model info
model.summary()

# learn
model.fit(train_images, train_labels, epochs=10)

# save model
model.save('trained_model', save_format='tf')
