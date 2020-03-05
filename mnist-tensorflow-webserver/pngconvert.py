import matplotlib.pyplot as plt
import matplotlib.image as mpimg
import numpy as np
from PIL import Image
import PIL.ImageOps   
import io


def pngToGreyscale28x28x8bit(png : bytes) -> np.ndarray:
    LUMINOSITY = 'L'

    pngio = io.BytesIO(png)
    img = Image.open(pngio).convert(LUMINOSITY)
    img = PIL.ImageOps.invert(img) # this is needed, no time to investigate the reason now
    img.thumbnail((28,28), PIL.Image.BOX)
    # imgplot = plt.imshow(img, cmap=plt.cm.binary)
    # plt.show()
    return np.asarray(img)