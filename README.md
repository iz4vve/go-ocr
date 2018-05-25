# go-ocr
OCR client for recognizing and transcribing files.


## Requirements and installation
The module requires a local installation of tesseract to work, and a working installation of the library and Go.

###### on Arch Linux
```bash
$ sudo pacman -S tesseract tesseract-data-eng
...

$ go get -u -v github.com/iz4vve/go-ocr
...

$./go-ocr /path/to/png
```