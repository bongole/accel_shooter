package main

import (
    "net/http"
    "fmt"
    "time"
    "strconv"
    "runtime"
    "runtime/debug"
)

type RequestContext struct {
    magick *Magick
    request_count int
}

type CtxList []*RequestContext
func (self CtxList) RequestTotal() int {
    sum := 0
    for _, v := range self {
        sum += v.request_count
    }
    return sum
}

type Chan chan *RequestContext

func (ch Chan) ServeImage( w http.ResponseWriter, r *http.Request ) {
    ctx := <- ch
    m := ctx.magick

    filename := "test.jpg"

    m.ReadImage(filename)
    m.Resize( 100, 100 )
    buf := m.GetImageBlob()

    w.Header().Add("Content-Type", http.DetectContentType(buf))
    w.Header().Add("Content-Length", strconv.Itoa(len(buf)))

    w.Write(buf)

    defer func() {
        m.Clear()

        ctx.request_count++

        ch <- ctx
    }()
}

func disableGC() {
    debug.SetGCPercent(-1)
}

func enableGC() {
    debug.SetGCPercent(100)
}

const REQUEST_GC_THRESHOLD = 100
func MyGC( ctxs CtxList ){
    prev_request_total := 0
    for {
        current_request_total := ctxs.RequestTotal()
        if d := current_request_total - prev_request_total; REQUEST_GC_THRESHOLD < d {
            enableGC()
            runtime.GC()
            fmt.Println("GC!!")
            disableGC()
            prev_request_total = current_request_total
        }

        time.Sleep( 1 * time.Second )
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    disableGC()

    go func(){
        fmt.Println("NumCPU = ", runtime.NumCPU())
        for {
            fmt.Println("NumGoroutine = ", runtime.NumGoroutine())
            time.Sleep( 1 * time.Second )
        }
    }()

    num_of_ch := runtime.NumCPU()
    ctxs := make( CtxList, num_of_ch)
    c := make( Chan, num_of_ch )

    for i := 0; i < num_of_ch; i++ {
        ctx := new(RequestContext)
        ctx.magick = NewMagick()
        ctx.request_count = 0

        ctxs[i] = ctx
        c <- ctx
    }

    go MyGC(ctxs)

    http.HandleFunc("/", c.ServeImage)
    http.ListenAndServe(":8080", nil)
}
