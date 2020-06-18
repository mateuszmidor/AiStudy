import cv2 as opencv

# Strings
CAPTURE_FAILED = "Capture failed"

# Image labelling params
GREEN = (0, 255, 0)
PREVIEW_WINDOW_NAME = "Labeled image"
TEXT_ORIGIN = (5, 20)
FONT_FACE = opencv.FONT_HERSHEY_PLAIN
FONT_SCALE = 2
FONT_THICKNESS = 2
FONT_LINE = opencv.LINE_AA