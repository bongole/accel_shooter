package main 

import (
    "os/exec"
    "fmt"
)

type GMagickCLI struct {
    filename string
    resize_width int
    resize_height int
}

func NewGMagickCLI() *GMagickCLI {
    m := new(GMagickCLI)
    return m
}

func (self *GMagickCLI) ReadImage( filename string ) {
    self.filename = filename
}

func (self *GMagickCLI) ResizeImage(width, height int) {
    self.resize_width = width
    self.resize_height = height
}

func (self *GMagickCLI) GetImageBlob() []byte {
    var ret []byte
    if self.filename != "" {
        if 0 < self.resize_width && 0 < self.resize_height {
            ret, _ = exec.Command("gm", "convert", "-resize", fmt.Sprintf("%dx%d!", self.resize_width, self.resize_height), self.filename, "-" ).Output()
        }
    }

    return ret
}

func (self *GMagickCLI) Clear() {
    self.filename = ""
    self.resize_width = 0
    self.resize_height = 0
}
