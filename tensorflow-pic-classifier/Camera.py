import Common as common
from Common import opencv 

class Camera(object):
    def __init__(self):
        # capture from default camera (0 is default)
        self.camera_capture = opencv.VideoCapture(0)

    def capture_frame(self, ignore_first_frame):
        # skip first frame if needed
        if (ignore_first_frame):
            self.camera_capture.read()

        # capture usable camera frame
        (capture_status, self.current_camera_frame) = self.camera_capture.read()

        # handle error if any
        if (capture_status):
            return self.current_camera_frame
        else:
            print(common.CAPTURE_FAILED)

    def display_image_with_label(self, image, label):
        # put label on image
        image_with_label = opencv.putText(
            image, label,
            common.TEXT_ORIGIN,
            common.FONT_FACE,
            common.FONT_SCALE,
            common.GREEN,
            common.FONT_THICKNESS,
            common.FONT_LINE
        )

        # display image
        opencv.imshow(common.PREVIEW_WINDOW_NAME, image_with_label)

        # wait key pressed
        opencv.waitKey()

    def display_current_frame_with_label(self, label):
        self.display_image_with_label(self.current_camera_frame, label)
