package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	CCHDEVICENAME                 = 32
	CCHFORMNAME                   = 32
	ENUM_CURRENT_SETTINGS  uint32 = 0xFFFFFFFF
	ENUM_REGISTRY_SETTINGS uint32 = 0xFFFFFFFE
	DISP_CHANGE_SUCCESSFUL uint32 = 0
	DISP_CHANGE_RESTART    uint32 = 1
	DISP_CHANGE_FAILED     uint32 = 0xFFFFFFFF
	DISP_CHANGE_BADMODE    uint32 = 0xFFFFFFFE
)

// DEVMODE is a structure used to
// specify characteristics of display
// and print devices.
type DEVMODE struct {
	DmDeviceName       [CCHDEVICENAME]uint16
	DmSpecVersion      uint16
	DmDriverVersion    uint16
	DmSize             uint16
	DmDriverExtra      uint16
	DmFields           uint32
	DmOrientation      int16
	DmPaperSize        int16
	DmPaperLength      int16
	DmPaperWidth       int16
	DmScale            int16
	DmCopies           int16
	DmDefaultSource    int16
	DmPrintQuality     int16
	DmColor            int16
	DmDuplex           int16
	DmYResolution      int16
	DmTTOption         int16
	DmCollate          int16
	DmFormName         [CCHFORMNAME]uint16
	DmLogPixels        uint16
	DmBitsPerPel       uint32
	DmPelsWidth        uint32
	DmPelsHeight       uint32
	DmDisplayFlags     uint32
	DmDisplayFrequency uint32
	DmICMMethod        uint32
	DmICMIntent        uint32
	DmMediaType        uint32
	DmDitherType       uint32
	DmReserved1        uint32
	DmReserved2        uint32
	DmPanningWidth     uint32
	DmPanningHeight    uint32
}

func main() {
	// CLI arguments
	argWidth, _ := strconv.ParseUint(os.Args[1], 10, 32)
	argHeight, _ := strconv.ParseUint(os.Args[2], 10, 32)
	argFrequency, _ := strconv.ParseUint(os.Args[3], 10, 32)

	// Load the DLL and the procedures.
	user32dll := syscall.NewLazyDLL("user32.dll")
	procEnumDisplaySettingsW := user32dll.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsW := user32dll.NewProc("ChangeDisplaySettingsW")

	// Get the display information.
	devMode := new(DEVMODE)
	ret, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(nil)),
		uintptr(ENUM_CURRENT_SETTINGS), uintptr(unsafe.Pointer(devMode)))

	if ret == 0 {
		fmt.Println("Couldn't extract display settings.")
		os.Exit(1)
	}

	// Show the display information.
	fmt.Printf("Screen resolution: %dx%d\n", devMode.DmPelsWidth, devMode.DmPelsHeight)
	fmt.Printf("Bits per pixel: %d\n", devMode.DmBitsPerPel)
	fmt.Printf("Display orientation: %d\n", devMode.DmOrientation)
	fmt.Printf("Display flags: %x\n", devMode.DmDisplayFlags)
	fmt.Printf("Display frequency: %d\n", devMode.DmDisplayFrequency)

	// Change the display resolution.
	newMode := *devMode
	newMode.DmPelsWidth = uint32(argWidth)
	newMode.DmPelsHeight = uint32(argHeight)
	newMode.DmDisplayFrequency = uint32(argFrequency)
	ret, _, _ = procChangeDisplaySettingsW.Call(uintptr(unsafe.Pointer(&newMode)),
		uintptr(0))

	switch ret {
	case uintptr(DISP_CHANGE_SUCCESSFUL):
		fmt.Println("Successfuly changed the display resolution and frequency.")

	case uintptr(DISP_CHANGE_RESTART):
		fmt.Println("Restart required to apply the resolution changes.")

	case uintptr(DISP_CHANGE_BADMODE):
		fmt.Println("The resolution or frequency are not supported by the display.")

	case uintptr(DISP_CHANGE_FAILED):
		fmt.Println("Failed to change the display resolution and frequency.")
	}
}
