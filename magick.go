package main 

//#cgo pkg-config: MagickWand
//#include <wand/MagickWand.h>
/*
MagickBooleanType MyMagickResizeImage(MagickWand *w,const size_t width,const size_t height,const FilterTypes filter,const double d){
       return MagickResizeImage(w,width,height,filter,d);
}
*/
import "C"

import (
    "unsafe"
    "reflect"
    "runtime"
)

type Magick struct {
    magick_wand *C.MagickWand
}

func init() {
    C.MagickWandGenesis()
}

func NewMagick() *Magick {
    m := new(Magick)
    m.magick_wand = C.NewMagickWand()
    return m
}

func (self *Magick) ReadImage( filename string ) {
    C.MagickReadImage( self.magick_wand, C.CString(filename) )
}

func (self *Magick) Resize(width, height int) {
    C.MyMagickResizeImage( self.magick_wand, C.size_t(width), C.size_t(height), C.LanczosFilter, 1.0 )
}

func (self *Magick) GetImageBlob() []byte {
    var image_size C.size_t
    ptr := C.MagickGetImageBlob( self.magick_wand, &image_size )

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
        C.ClearMagickWand(self.magick_wand)
    }
}

func (self *Magick) Destroy() {
    if self.magick_wand != nil {
        C.DestroyMagickWand(self.magick_wand)
        self.magick_wand = nil
    }
}
