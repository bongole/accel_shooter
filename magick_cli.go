package main 

import (
    "os/exec"
    "fmt"
)

type MagickCLI struct {
    filename string
    resize_width int
    resize_height int
}

func NewMagickCLI() *MagickCLI {
    m := new(MagickCLI)
    return m
}

func (self *MagickCLI) ReadImage( filename string ) {
    self.filename = filename
}

func (self *MagickCLI) ResizeImage(width, height int) {
    self.resize_width = width
    self.resize_height = height
}

func (self *MagickCLI) GetImageBlob() []byte {
    var ret []byte
    if self.filename != "" {
        if 0 < self.resize_width && 0 < self.resize_height {
            ret, _ = exec.Command("convert", "-resize", fmt.Sprintf("%dx%d!", self.resize_width, self.resize_height), self.filename, "-" ).Output()
        }
    }

    return ret
}

func (self *MagickCLI) Clear() {
    self.filename = ""
    self.resize_width = 0
    self.resize_height = 0
}
