import matplotlib.pyplot as plt
import matplotlib.image as mpimg
import numpy as np
from PIL import Image
import PIL.ImageOps   
import io


def pngToGreyscale28x28x8bit(png : bytes) -> np.ndarray:
    LUMINOSITY = 'L'

    pngio = io.BytesIO(png)
    img = Image.open(pngio)
    img = img.convert(LUMINOSITY)
    img.thumbnail((28,28), PIL.Image.LINEAR)
    return np.asarray(img)