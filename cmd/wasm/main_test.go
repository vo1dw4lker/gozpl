//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
	"testing"
)

var api = newAPI()

func TestChaining(t *testing.T) {
	got := api.Call("newLabel").
		Call("setPrintWidth", 800).
		Call("addText", 10, 20, "Hello", map[string]any{"font": "0", "height": 30, "width": 25, "orientation": "R"}).
		Call("addCode128", 10, 60, "12345", map[string]any{"height": 80, "printText": false}).
		Call("addQRCode", 10, 160, "https://example.com", map[string]any{"magnification": 4}).
		Call("addGraphicBox", 0, 0, 100, 50, map[string]any{"thickness": 3, "color": "W", "rounding": 2}).
		Call("toString").String()

	want := "^XA\n" +
		"^PW800\n" +
		"^FO10,20^A0R,30,25^FDHello^FS\n" +
		"^FO10,60^BCN,80,N,N,N,N^FD12345^FS\n" +
		"^FO10,160^BQN,2,4^FDQA,https://example.com^FS\n" +
		"^FO0,0^GB100,50,3,W,2^FS\n" +
		"^XZ\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDefaults(t *testing.T) {
	got := api.Call("newLabel").
		Call("addText", 1, 2, "a").
		Call("addText", 1, 2, "b", map[string]any{"height": 40}).
		Call("addText", 1, 2, "c", map[string]any{"fieldBlock": map[string]any{"width": 300}}).
		Call("toString").String()

	want := "^XA\n" +
		"^FO1,2^AAN,15,15^FDa^FS\n" +
		"^FO1,2^AAN,40,15^FDb^FS\n" +
		"^FO1,2^AAN,15,15^FB300,1,0,L^FDc^FS\n" +
		"^XZ\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGraphics(t *testing.T) {
	data := js.Global().Get("Uint8Array").New(2)
	js.CopyBytesToJS(data, []byte{0xC0, 0xC0})

	// 8x1 RGBA image, first 4 pixels black, rest white.
	pix := js.Global().Get("Uint8ClampedArray").New(32)
	b := make([]byte, 32)
	for i := 16; i < 32; i++ {
		b[i] = 0xFF
	}
	for i := 3; i < 32; i += 4 {
		b[i] = 0xFF
	}
	js.CopyBytesToJS(pix, b)
	img := map[string]any{"width": 8, "height": 1, "data": pix}

	got := api.Call("newLabel").
		Call("addGraphicField", 10, 20, 1, data).
		Call("addImage", 0, 0, img).
		Call("toString").String()

	want := "^XA\n^FO10,20^GFA,2,2,1,C0C0^FS\n^FO0,0^GFA,1,1,1,F0^FS\n^XZ\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReset(t *testing.T) {
	got := api.Call("newLabel").Call("addRaw", "^FX").Call("reset").Call("toString").String()
	if got != "^XA\n^XZ\n" {
		t.Errorf("got %q", got)
	}
}

func TestErrorsThrow(t *testing.T) {
	cases := []struct {
		name string
		call func()
		msg  string
	}{
		{"bad int", func() { api.Call("newLabel").Call("setQuantity", "5") }, "quantity must be a number"},
		{"missing arg", func() { api.Call("newLabel").Call("addText", 1, 2) }, "text must be a string"},
		{"bad option type", func() { api.Call("newLabel").Call("addQRCode", 1, 2, "x", map[string]any{"magnification": "4"}) }, "options.magnification must be a number"},
		{"bad orientation", func() { api.Call("newLabel").Call("addText", 1, 2, "x", map[string]any{"orientation": "X"}) }, "one of N, R, I, B"},
		{"bad image", func() {
			api.Call("newLabel").Call("addImage", 0, 0, map[string]any{"width": 2, "height": 2, "data": js.Global().Get("Uint8Array").New(3)})
		}, "width*height*4"},
		{"detached method", func() { api.Call("newLabel").Get("toString").Invoke() }, "non-label value"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				r := recover()
				err, ok := r.(js.Error)
				if !ok {
					t.Fatalf("expected a thrown js.Error, got %v", r)
				}
				if !strings.Contains(err.Error(), c.msg) {
					t.Errorf("error %q does not contain %q", err.Error(), c.msg)
				}
			}()
			c.call()
		})
	}
}

func TestStillUsableAfterError(t *testing.T) {
	func() {
		defer func() { recover() }()
		api.Call("newLabel").Call("setQuantity", nil)
	}()
	if got := api.Call("newLabel").Call("setQuantity", 2).Call("toString").String(); got != "^XA\n^PQ2\n^XZ\n" {
		t.Errorf("got %q", got)
	}
}
