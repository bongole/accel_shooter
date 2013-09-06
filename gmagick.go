package main 

//#cgo pkg-config: GraphicsMagickWand
//#include <wand/magick_wand.h>
/*
unsigned int MagickResizeImage(MagickWand *,const unsigned long,const unsigned long,const FilterTypes,const double);
unsigned char *MagickWriteImageBlob(MagickWand *,size_t *);
*/
import "C"

import (
    "unsafe"
    "reflect"
    "runtime"
    "os"
)

type Magick struct {
    magick_wand *C.MagickWand
}

func init() {
    C.InitializeMagick(C.CString(os.Args[0]))
}

func NewMagick() *Magick {
    m := new(Magick)
    m.magick_wand = C.NewMagickWand()
    return m
}

func (self *Magick) ReadImage( filename string ) {
    C.MagickReadImage( self.magick_wand, C.CString(filename) )
}

func (self *Magick) ResizeImage(width, height int) {
    C.MagickResizeImage( self.magick_wand, C.ulong(width), C.ulong(height), C.LanczosFilter, 1.0 )
}

func (self *Magick) GetImageBlob() []byte {
    var image_size C.size_t
    ptr := C.MagickWriteImageBlob( self.magick_wand, &image_size )

    var theGoSlice []byte
    sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&theGoSlice)))
    sliceHeader.Cap = int(image_size)
    sliceHeader.Len = int(image_size)
    sliceHeader.Data = uintptr(unsafe.Pointer(ptr))

    runtime.SetFinalizer( &theGoSlice, func(x *[]byte){
        h := (*reflect.SliceHeader)((unsafe.Pointer(x)))
        C.MagickRelinquishMemory(unsafe.Pointer(h.Data))
    })

    return theGoSlice
}

func (self *Magick) Clear() {
    if self.magick_wand != nil {
        self.Destroy()
        self.magick_wand = C.NewMagickWand()
    }
}

func (self *Magick) Destroy() {
    if self.magick_wand != nil {
        C.DestroyMagickWand(self.magick_wand)
        self.magick_wand = nil
    }
}
