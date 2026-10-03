//go:build js && wasm

package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"syscall/js"

	zpl "github.com/vo1dw4lker/gozpl"
)

const idKey = "__gozplLabel"

var (
	labels = map[int]*zpl.Label{}
	nextID int

	throwIfError = js.Global().Get("Function").New("f",
		"return function(...args) { const r = f.apply(this, args); if (r instanceof Error) throw r; return r; }")
)

func main() {
	js.Global().Set("gozpl", newAPI())
	select {}
}

func newAPI() js.Value {
	object := js.Global().Get("Object")

	proto := object.New()
	for name, m := range methods {
		proto.Set(name, wrap(bind(m)))
	}

	registry := js.Global().Get("FinalizationRegistry").New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		delete(labels, args[0].Int())
		return nil
	}))

	api := object.New()
	api.Set("newLabel", wrap(func(js.Value, []js.Value) any {
		nextID++
		labels[nextID] = zpl.NewLabel()
		obj := object.Call("create", proto)
		object.Call("defineProperty", obj, idKey, map[string]any{"value": nextID})
		registry.Call("register", obj, nextID)
		return obj
	}))
	return api
}

func wrap(fn func(js.Value, []js.Value) any) js.Value {
	return throwIfError.Invoke(js.FuncOf(func(this js.Value, args []js.Value) (ret any) {
		defer func() {
			if r := recover(); r != nil {
				ret = js.Global().Get("Error").New(fmt.Sprint(r))
			}
		}()
		return fn(this, args)
	}))
}

type method func(l *zpl.Label, args []js.Value) any

func bind(m method) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		if this.Type() != js.TypeObject || this.Get(idKey).Type() != js.TypeNumber {
			panic(errors.New("gozpl: method called on a non-label value"))
		}
		l, ok := labels[this.Get(idKey).Int()]
		if !ok {
			panic(errors.New("gozpl: label no longer exists"))
		}
		if r := m(l, args); r != nil {
			return r
		}
		return this
	}
}

var methods = map[string]method{
	"toString": func(l *zpl.Label, _ []js.Value) any { return l.String() },
	"reset":    func(l *zpl.Label, _ []js.Value) any { l.Reset(); return nil },

	"setQuantity":    func(l *zpl.Label, a []js.Value) any { l.SetQuantity(intArg(a, 0, "quantity")); return nil },
	"setPrintWidth":  func(l *zpl.Label, a []js.Value) any { l.SetPrintWidth(intArg(a, 0, "width")); return nil },
	"setLabelLength": func(l *zpl.Label, a []js.Value) any { l.SetLabelLength(intArg(a, 0, "length")); return nil },
	"setPrintRate":   func(l *zpl.Label, a []js.Value) any { l.SetPrintRate(intArg(a, 0, "rate")); return nil },
	"addRaw":         func(l *zpl.Label, a []js.Value) any { l.AddRaw(strArg(a, 0, "cmd")); return nil },

	// addText(x, y, text, { font, height, width, orientation,
	//                       fieldBlock: { width, maxLines, lineSpacing, alignment } })
	"addText": func(l *zpl.Label, a []js.Value) any {
		o := optArg(a, 3)
		var opts []zpl.TextOption
		if o.has("font") || o.has("height") || o.has("width") {
			// Same defaults as zpl.AddText.
			opts = append(opts, zpl.WithFont(o.str("font", "A"), o.int("height", 15), o.int("width", 15)))
		}
		if o.has("orientation") {
			opts = append(opts, zpl.WithTextOrientation(o.orientation("orientation")))
		}
		if o.has("fieldBlock") {
			fb := options{o.v.Get("fieldBlock"), "fieldBlock"}
			if !fb.has("width") {
				panic(errors.New("fieldBlock.width is required"))
			}
			opts = append(opts, zpl.WithFieldBlock(fb.int("width", 0), fb.int("maxLines", 1), fb.int("lineSpacing", 0), fb.str("alignment", "L")))
		}
		l.AddText(intArg(a, 0, "x"), intArg(a, 1, "y"), strArg(a, 2, "text"), opts...)
		return nil
	},

	// addCode128(x, y, data, { height, orientation, printText, textAbove })
	"addCode128": func(l *zpl.Label, a []js.Value) any {
		o := optArg(a, 3)
		var opts []zpl.BarcodeOption
		if o.has("height") {
			opts = append(opts, zpl.WithBarcodeHeight(o.int("height", 0)))
		}
		if o.has("orientation") {
			opts = append(opts, zpl.WithBarcodeOrientation(o.orientation("orientation")))
		}
		if o.has("printText") || o.has("textAbove") {
			// Same defaults as zpl.AddCode128.
			opts = append(opts, zpl.WithBarcodeText(o.bool("printText", true), o.bool("textAbove", false)))
		}
		l.AddCode128(intArg(a, 0, "x"), intArg(a, 1, "y"), strArg(a, 2, "data"), opts...)
		return nil
	},

	// addQRCode(x, y, data, { magnification })
	"addQRCode": func(l *zpl.Label, a []js.Value) any {
		o := optArg(a, 3)
		var opts []zpl.QROption
		if o.has("magnification") {
			opts = append(opts, zpl.WithQRMagnification(o.int("magnification", 0)))
		}
		l.AddQRCode(intArg(a, 0, "x"), intArg(a, 1, "y"), strArg(a, 2, "data"), opts...)
		return nil
	},

	// addGraphicBox(x, y, width, height, { thickness, color, rounding })
	"addGraphicBox": func(l *zpl.Label, a []js.Value) any {
		o := optArg(a, 4)
		var opts []zpl.BoxOption
		if o.has("thickness") {
			opts = append(opts, zpl.WithBoxThickness(o.int("thickness", 0)))
		}
		if o.has("color") {
			opts = append(opts, zpl.WithBoxColor(o.str("color", "")))
		}
		if o.has("rounding") {
			opts = append(opts, zpl.WithBoxRounding(o.int("rounding", 0)))
		}
		l.AddGraphicBox(intArg(a, 0, "x"), intArg(a, 1, "y"), intArg(a, 2, "width"), intArg(a, 3, "height"), opts...)
		return nil
	},

	// addGraphicField(x, y, bytesPerRow, data: Uint8Array)
	"addGraphicField": func(l *zpl.Label, a []js.Value) any {
		l.AddGraphicField(intArg(a, 0, "x"), intArg(a, 1, "y"), intArg(a, 2, "bytesPerRow"), bytesArg(a, 3, "data"))
		return nil
	},

	// addImage(x, y, image) where image is either an ImageData-like
	// { width, height, data } with RGBA bytes, or a Uint8Array holding an
	// encoded PNG or JPEG.
	"addImage": func(l *zpl.Label, a []js.Value) any {
		l.AddImage(intArg(a, 0, "x"), intArg(a, 1, "y"), imageArg(a, 2, "image"))
		return nil
	},
}

func intArg(args []js.Value, i int, name string) int {
	if i >= len(args) || args[i].Type() != js.TypeNumber {
		panic(fmt.Errorf("%s must be a number", name))
	}
	return args[i].Int()
}

func strArg(args []js.Value, i int, name string) string {
	if i >= len(args) || args[i].Type() != js.TypeString {
		panic(fmt.Errorf("%s must be a string", name))
	}
	return args[i].String()
}

func isBytes(v js.Value) bool {
	return v.InstanceOf(js.Global().Get("Uint8Array")) || v.InstanceOf(js.Global().Get("Uint8ClampedArray"))
}

func bytesArg(args []js.Value, i int, name string) []byte {
	if i >= len(args) || !isBytes(args[i]) {
		panic(fmt.Errorf("%s must be a Uint8Array", name))
	}
	b := make([]byte, args[i].Length())
	js.CopyBytesToGo(b, args[i])
	return b
}

func imageArg(args []js.Value, i int, name string) image.Image {
	if i >= len(args) || args[i].Type() != js.TypeObject {
		panic(fmt.Errorf("%s must be ImageData or a Uint8Array", name))
	}
	v := args[i]
	if isBytes(v) {
		img, _, err := image.Decode(bytes.NewReader(bytesArg(args, i, name)))
		if err != nil {
			panic(fmt.Errorf("%s: %w", name, err))
		}
		return img
	}

	o := options{v, name}
	w, h := o.int("width", 0), o.int("height", 0)
	if !isBytes(v.Get("data")) {
		panic(fmt.Errorf("%s.data must be a Uint8Array or Uint8ClampedArray", name))
	}
	if w <= 0 || h <= 0 || v.Get("data").Length() != 4*w*h {
		panic(fmt.Errorf("%s: data length must be width*height*4 RGBA bytes", name))
	}
	pix := make([]byte, 4*w*h)
	js.CopyBytesToGo(pix, v.Get("data"))
	// ImageData is not alpha-premultiplied, so it maps to NRGBA.
	return &image.NRGBA{Pix: pix, Stride: 4 * w, Rect: image.Rect(0, 0, w, h)}
}

// options reads fields of an optional JavaScript options object.
type options struct {
	v    js.Value
	name string
}

func optArg(args []js.Value, i int) options {
	if i >= len(args) || args[i].IsUndefined() || args[i].IsNull() {
		return options{js.Undefined(), "options"}
	}
	if args[i].Type() != js.TypeObject {
		panic(errors.New("options must be an object"))
	}
	return options{args[i], "options"}
}

func (o options) has(key string) bool {
	return o.v.Type() == js.TypeObject && !o.v.Get(key).IsUndefined()
}

func (o options) get(key string, t js.Type) (js.Value, bool) {
	if !o.has(key) {
		return js.Value{}, false
	}
	v := o.v.Get(key)
	if v.Type() != t {
		panic(fmt.Errorf("%s.%s must be a %s", o.name, key, t))
	}
	return v, true
}

func (o options) int(key string, def int) int {
	if v, ok := o.get(key, js.TypeNumber); ok {
		return v.Int()
	}
	return def
}

func (o options) str(key, def string) string {
	if v, ok := o.get(key, js.TypeString); ok {
		return v.String()
	}
	return def
}

func (o options) bool(key string, def bool) bool {
	if v, ok := o.get(key, js.TypeBoolean); ok {
		return v.Bool()
	}
	return def
}

func (o options) orientation(key string) zpl.Orientation {
	switch s := zpl.Orientation(o.str(key, "")); s {
	case zpl.OrientationNormal, zpl.OrientationRotated, zpl.OrientationInverted, zpl.OrientationBottomUp:
		return s
	default:
		panic(fmt.Errorf("%s.%s must be one of N, R, I, B", o.name, key))
	}
}
