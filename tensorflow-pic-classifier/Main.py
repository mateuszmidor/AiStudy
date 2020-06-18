import Common as common 
import Camera as cameraModule 
import Inference as inferenceModule

if __name__ == "__main__":
    model_file_path = 'model/mobilenet_v1_1.0_224_quant.tflite'
    labels_file_path = 'model/labels_mobilenet_quant_v1_224.txt'

    model = inferenceModule.Inference(model_file_path, labels_file_path)
    camera = cameraModule.Camera()
    
    image = camera.capture_frame(True)
    label = model.label_image(image)

    camera.display_current_frame_with_label(label)
