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
	procEnumDisplayDevicesW := user32dll.NewProc("EnumDisplayDevicesW")
	procEnumDisplaySettingsW := user32dll.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsExW := user32dll.NewProc("ChangeDisplaySettingsExW")

	// Find the secondary monitor
	var dd DISPLAY_DEVICE
	dd.Cb = uint32(unsafe.Sizeof(dd))
	var i uint32 = 0
	var secondaryMonitorName *uint16

	for {
		ret, _, _ := procEnumDisplayDevicesW.Call(0, uintptr(i), uintptr(unsafe.Pointer(&dd)), 0)
		if ret == 0 {
			break
		}
		if dd.StateFlags&DISPLAY_DEVICE_ATTACHED_TO_DESKTOP != 0 && dd.StateFlags&DISPLAY_DEVICE_PRIMARY_DEVICE == 0 {
			secondaryMonitorName = &dd.DeviceName[0]
			break
		}
		i++
	}

	if secondaryMonitorName == nil {
		fmt.Println("No secondary monitor found.")
		return
	}

	// Get the display information for the secondary monitor
	devMode := new(DEVMODE)
	devMode.DmSize = uint16(unsafe.Sizeof(*devMode))
	ret, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(secondaryMonitorName)),
		uintptr(ENUM_CURRENT_SETTINGS), uintptr(unsafe.Pointer(devMode)))
	if ret == 0 {
		fmt.Println("Couldn't extract display settings for the secondary monitor.")
		return
	}

	// Show the current display information for the secondary monitor
	fmt.Printf("Secondary monitor resolution: %dx%d\n", devMode.DmPelsWidth, devMode.DmPelsHeight)
	fmt.Printf("Bits per pixel: %d\n", devMode.DmBitsPerPel)
	fmt.Printf("Display frequency: %d\n", devMode.DmDisplayFrequency)

	// Change the display resolution for the secondary monitor
	newMode := *devMode
	newMode.DmPelsWidth = uint32(argWidth)
	newMode.DmPelsHeight = uint32(argHeight)
	newMode.DmDisplayFrequency = uint32(argFrequency)
	newMode.DmFields = DM_PELSWIDTH | DM_PELSHEIGHT | DM_DISPLAYFREQUENCY

	ret, _, _ = procChangeDisplaySettingsExW.Call(
		uintptr(unsafe.Pointer(secondaryMonitorName)),
		uintptr(unsafe.Pointer(&newMode)),
		0,
		CDS_UPDATEREGISTRY|CDS_NORESET,
		0)

	switch ret {
	case uintptr(DISP_CHANGE_SUCCESSFUL):
		fmt.Println("Successfully changed the secondary monitor resolution and frequency.")
	case uintptr(DISP_CHANGE_RESTART):
		fmt.Println("Restart required to apply the resolution changes.")
	case uintptr(DISP_CHANGE_BADMODE):
		fmt.Println("The resolution or frequency are not supported by the secondary monitor.")
	case uintptr(DISP_CHANGE_FAILED):
		fmt.Println("Failed to change the secondary monitor resolution and frequency.")
	}

	// Apply the changes
	procChangeDisplaySettingsExW.Call(0, 0, 0, 0, 0)
}

type DISPLAY_DEVICE struct {
	Cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

const (
	DISPLAY_DEVICE_ATTACHED_TO_DESKTOP = 0x00000001
	DISPLAY_DEVICE_PRIMARY_DEVICE      = 0x00000004
	DM_PELSWIDTH                       = 0x00080000
	DM_PELSHEIGHT                      = 0x00100000
	DM_DISPLAYFREQUENCY                = 0x00400000
	CDS_UPDATEREGISTRY                 = 0x00000001
	CDS_NORESET                        = 0x10000000
)
