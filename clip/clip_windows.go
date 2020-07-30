// +build windows

package clip

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"log"
	"syscall"
	"time"
	"unsafe"

	"github.com/gen2brain/dlgs"
	"github.com/jsummers/gobmp"
	"github.com/kbinani/screenshot"
	"github.com/pkg/errors"
)

const (
	cfBitmap      = 2
	cfDib         = 8
	cfDibV5       = 17
	cfUnicodetext = 13
	gmemMoveable  = 0x0002
	fileHeaderLen = 14
	//infoHeaderLen = 40
	//infoHeaderLen = 124
)

type infoHeader struct {
	iSize          uint32 // infoheader длина 40
	iWidth         uint32 // ширина
	iHeight        uint32 // высота
	iPLanes        uint16
	iBitCount      uint16 // бит на пиксель
	iCompression   uint32 // формат сжатия, dib - обычно от нуля до трёх
	iSizeImage     uint32 // размер изображения в пикселях?
	iXPelsPerMeter uint32
	iYPelsPerMeter uint32
	iClrUsed       uint32
	iClrImportant  uint32
}

var (
	user32                     = syscall.MustLoadDLL("user32")
	openClipboard              = user32.MustFindProc("OpenClipboard")
	closeClipboard             = user32.MustFindProc("CloseClipboard")
	emptyClipboard             = user32.MustFindProc("EmptyClipboard")
	getClipboardData           = user32.MustFindProc("GetClipboardData")
	setClipboardData           = user32.MustFindProc("SetClipboardData")
	isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable")

	kernel32       = syscall.NewLazyDLL("kernel32")
	globalAlloc    = kernel32.NewProc("GlobalAlloc")
	globalFree     = kernel32.NewProc("GlobalFree")
	globalLock     = kernel32.NewProc("GlobalLock")
	globalUnlock   = kernel32.NewProc("GlobalUnlock")
	procGlobalSize = kernel32.NewProc("GlobalSize")
	lstrcpy        = kernel32.NewProc("lstrcpyW")
	copyMemory     = kernel32.NewProc("CopyMemory")
)

func globalSize(hMem unsafe.Pointer) uint {
	ret, _, _ := procGlobalSize.Call(uintptr(hMem))

	if ret == 0 {
		panic("GlobalSize failed")
	}

	return uint(ret)
}

// waitOpenClipboard opens the clipboard, waiting for up to a second to do so.
func waitOpenClipboard() error {
	started := time.Now()
	limit := started.Add(time.Second)
	var r uintptr
	var err error
	for time.Now().Before(limit) {
		r, _, err = openClipboard.Call(0)
		if r != 0 {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return err
}

func getClipboard() ([]byte, error) {
	err := waitOpenClipboard()
	if err != nil {
		return nil, errors.Wrapf(err, "Error at waitOpenClipboard: %v")
	}
	defer closeClipboard.Call()

	r, _, err := isClipboardFormatAvailable.Call(cfDib)
	if r == 0 {
		log.Println("not Dib format: ", err)
		//dlgs.Error("Error", "Not DIB format")
		b := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureRect(b)
		if err != nil {
			return nil, err
		}
		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, img)
		if err == nil {
			return buffer.Bytes(), nil
		}
		return nil, err
	}

	h, _, err := getClipboardData.Call(cfDib)
	if h == 0 {
		log.Println("get clipbord data error: ", err)
		dlgs.Error("Error", "Can't get clipbord data")
		return nil, err
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		dlgs.Error("Error", "Global lock error")
		log.Println("globalLock error: ", err)
		return nil, err
	}

	h2 := (*infoHeader)(unsafe.Pointer(l))

	const infoHeaderLen = 124
	//ihl := new(bytes.Buffer)
	//binary.Write(ihl, binary.LittleEndian, *(*byte)(unsafe.Pointer(l + len(h2) + fileHeaderLen)))

	dataSize := h2.iSizeImage + fileHeaderLen + infoHeaderLen

	if h2.iSizeImage == 0 && h2.iCompression == 0 {
		iSizeImage := h2.iHeight * ((h2.iWidth*uint32(h2.iBitCount)/8 + 3) &^ 3)
		dataSize += iSizeImage
	}
	log.Println("datasize: ", dataSize, h2.iHeight*((h2.iWidth*uint32(h2.iBitCount)/8+3)&^3))

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.LittleEndian, uint16('B')|(uint16('M')<<8))
	binary.Write(buffer, binary.LittleEndian, uint32(dataSize))
	binary.Write(buffer, binary.LittleEndian, uint32(0))
	const sizeof_colorbar = 0
	binary.Write(buffer, binary.LittleEndian, uint32(fileHeaderLen+infoHeaderLen+sizeof_colorbar))
	log.Println("fileHeader ", buffer.Bytes(), len(buffer.Bytes()))

	j := 0
	for i := fileHeaderLen; i < int(dataSize); i++ {
		binary.Write(buffer, binary.BigEndian, *(*byte)(unsafe.Pointer(l + uintptr(j))))
		j++
	}

	img, err := gobmp.Decode(buffer)
	if err != nil {
		dlgs.Error("Error", "Can't decode BMP")
		return nil, err
	}
	err = png.Encode(buffer, img)
	if err != nil {
		dlgs.Error("Error", "Can't encode to PNG")
		return nil, err
	}

	r, _, err = globalUnlock.Call(h)
	if r == 0 {
		dlgs.Error("Error", "Global Unlock error")
		return nil, err
	}

	return buffer.Bytes(), nil
}
