import http.server
import socketserver
import json
import base64
import typing


def extractBase64BytesFromHttpPostData(base64_post_data : bytes) -> bytes:
    comma_pos = base64_post_data.find(b',')
    return base64_post_data[comma_pos+1:]

class CheckDigitHandler(http.server.BaseHTTPRequestHandler):
    checkDigit : typing.Callable[[bytes], str] = None # function; input: PNG bytes, output: recognized digit as string

    def do_GET(self):
        if self.path.lower() == "/":
            self.insertHeaders200OK()  
            self.handleMain()
            return

        self.insertHeaders404NotFound()
        self.handle404NotFound()

    def do_POST(self):
        if self.path.lower() == "/checkdigit":
            self.insertHeaders200OK()
            self.handleCheckDigit()
            return 

        self.insertHeaders404NotFound()

    def insertHeaders200OK(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()

    def insertHeaders404NotFound(self):
        self.send_response(404)
        self.send_header("Content-type", "text/html")
        self.end_headers()

    def handleMain(self):
        with open("main.html") as f:
            self.writeString(f.read())

    def handleCheckDigit(self):
        content_length = int(self.headers['Content-Length']) # <--- Gets the size of data
        post_data = self.rfile.read(content_length) # <--- Gets the data itself
        base64_png_bytes = extractBase64BytesFromHttpPostData(post_data)
        raw_png_bytes = base64.b64decode(base64_png_bytes)
        recognized_digit = str(CheckDigitHandler.checkDigit(raw_png_bytes))
        result = {"value" : recognized_digit}
        result_json_string = json.dumps(result)    
        self.writeString(result_json_string)

    def handle404NotFound(self):
        self.writeString("<html><head><title>404 Not Found</title></head>")
        self.writeString("<body><p>404 Not Found</p>")
        self.writeString("</body></html>")

    def writeString(self, str: str):
        bytes = bytearray(str, "UTF-8")
        self.wfile.write(bytes)


def runServer(port : int, checkDigitHandler : typing.Callable[[bytes], str]):
    CheckDigitHandler.checkDigit = checkDigitHandler

    with socketserver.TCPServer(("", port), CheckDigitHandler) as httpd:
        print("serving at port", port)
        httpd.serve_forever()