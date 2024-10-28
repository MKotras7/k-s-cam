import cv2
import requests
import numpy as np
from datetime import datetime

# Replace with the URL of your ESP32-CAM MJPEG stream
url = 'http://192.168.1.201/stream'

try:
    # Open a persistent connection to the MJPEG stream
    stream = requests.get(url, stream=True)

    if stream.status_code == 200:
        bytes_data = bytes()

        for chunk in stream.iter_content(chunk_size=1024):
            bytes_data += chunk
            a = bytes_data.find(b'\xff\xd8')
            b = bytes_data.find(b'\xff\xd9')

            if a != -1 and b != -1:
                jpg = bytes_data[a:b+2]
                bytes_data = bytes_data[b+2:]

                # Decode JPEG to a NumPy array
                image = cv2.imdecode(np.frombuffer(jpg, dtype=np.uint8), cv2.IMREAD_COLOR)

                # Display the image in a window
                cv2.imshow('MJPEG Stream', image)

                # Save the image with a timestamp
                timestamp = datetime.now().strftime("%Y%m%d_%H%M%S%f")
                filename = f"frame_{timestamp}.jpg"
                cv2.imwrite(filename, image)

                # Break the loop if 'q' is pressed
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    break
    else:
        print("Failed to connect to the stream.")

except KeyboardInterrupt:
    print("Stream interrupted by user.")

finally:
    # Ensure the OpenCV window is closed properly
    cv2.destroyAllWindows()