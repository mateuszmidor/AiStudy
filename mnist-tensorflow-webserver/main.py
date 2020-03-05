import server
import pngconvert
import inference
import sys

def checkDigit(raw_png_bytes : bytes) -> str:
    greyscale_28x28x8bit = pngconvert.pngToGreyscale28x28x8bit(raw_png_bytes)
    digit = inference.recognizeDigit(greyscale_28x28x8bit)
    return digit

if __name__ == '__main__':
    port = int(sys.argv[1])
    
    # stopping docker results in listening socket exception
    try:
        server.runServer(port, checkDigit)
    except:
        print("Digit recognition server killed. Bye!")