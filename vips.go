package main 

//#cgo pkg-config: vips
//#include <math.h>
//#include <vips/vips.h>
/*
static int
calculate_shrink( int width, int height, int thumbnail_size, double *residual )
{
    int dimension = IM_MAX( width, height );

    double factor = dimension / (double) thumbnail_size;

    double factor2 = factor < 1.0 ? 1.0 : factor;

    int shrink = floor( factor2 );

    int isize = floor( dimension / shrink );

    if( residual )
        *residual = thumbnail_size / (double) isize;

    return( shrink );
}

static int
myvips_foreign_load( const char *filename, VipsImage **out, int shrink ) 
{
    //return vips_foreign_load( filename, out, "sequential", TRUE, "shrink", shrink, NULL );
    return vips_foreign_load( filename, out, "shrink", shrink, NULL );
}

static int
myvips_jpegsave_buffer(VipsImage *in, void **buf,size_t *len)
{
    return vips_jpegsave_buffer( in, buf, len, NULL );
}

static VipsInterpolate*
myvips_interp( const char* name )
{
    return VIPS_INTERPOLATE( vips_object_new_from_string(g_type_class_ref( VIPS_TYPE_INTERPOLATE ), name ) );
}

// static int
// myvips_tilecache( VipsImage *in, VipsImage **out, int width )
// {
//     return vips_tilecache( in, out, 
//                            "tile_width", width,
//                            "tile_height", 10, 
//                            "max_tiles", 600, 
//                            "strategy", VIPS_CACHE_SEQUENTIAL, 
//                            "threaded", TRUE, NULL );
// }

static INTMASK *
sharpen_filter( void )
{

    INTMASK *mask = im_create_imaskv( "sharpen.con", 3, 3,
        -1, -1, -1,
        -1, 32, -1,
        -1, -1, -1 );
    mask->scale = 24;

    return( mask );
}

*/
import "C"

import (
    "runtime"
    "unsafe"
    "reflect"
    "math"
    "os"
)

type Magick struct {
    img_in *C.IMAGE
    img_out *C.IMAGE
    filename string
}

var sharpen_mask *C.INTMASK 
var interp *C.VipsInterpolate
func init() {
    C.im_init_world(C.CString(os.Args[0]))
    sharpen_mask = C.sharpen_filter()
    interp = C.myvips_interp(C.CString("nearest"))
}

func NewMagick() *Magick {
    m := new(Magick)
    return m
}

func (self *Magick) ReadImage( filename string ) {
    self.filename = filename
}

func (self *Magick) ResizeImage(width, height int) {
    tmp_in := C.im_open( C.CString(self.filename), C.CString("r") )
    
    shrink := C.calculate_shrink( tmp_in.Xsize, tmp_in.Ysize, C.int(math.Max(float64(width), float64(height))), nil )
    if( shrink > 8 ) {
        shrink = 8
    } else if( shrink > 4 ) {
        shrink = 4
    } else if( shrink > 2 ) {
        shrink = 2
    } else {
        shrink = 1
    }

    C.myvips_foreign_load( C.CString(self.filename), &self.img_in, shrink )

    self.img_out = C.vips_image_new()
    tmp_out := C.vips_image_new()

    C.im_conv(self.img_in, tmp_out, sharpen_mask)
    x_scale, y_scale := C.double(width)/C.double(self.img_in.Xsize), C.double(height)/C.double(self.img_in.Ysize)
    C.im_affinei_all(tmp_out, self.img_out, interp, x_scale, 0, 0, y_scale, 0, 0)

    defer func(){
        C.im_close( tmp_in )
        C.im_close( tmp_out )
    }()

}

func (self *Magick) GetImageBlob() []byte {
    var buf *C.void
    var buf_len C.size_t
    C.myvips_jpegsave_buffer( self.img_out, (*unsafe.Pointer)(unsafe.Pointer(&buf)), &buf_len )

    var theGoSlice []byte
    sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&theGoSlice)))
    sliceHeader.Cap = int(buf_len)
    sliceHeader.Len = int(buf_len)
    sliceHeader.Data = uintptr(unsafe.Pointer(buf))

    runtime.SetFinalizer( &theGoSlice, func(x *[]byte){
        h := (*reflect.SliceHeader)((unsafe.Pointer(x)))
        C.g_free(C.gpointer(unsafe.Pointer(h.Data)))
    })

    return theGoSlice
}

func (self *Magick) Clear() {
    if self.img_in != nil {
        C.im_close(self.img_in)
        self.img_in = nil
    }

    if self.img_out != nil {
        C.im_close(self.img_out)
        self.img_in = nil
    }
}
