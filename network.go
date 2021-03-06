// +build windows

package winapi

import (
	"fmt"
	"unsafe"

	so "github.com/Codehardt/go-win64api/shared"
)

var shNetShareEnum = modNetapi32.NewProc("NetShareEnum")

const SHARE_MAX_PREFERRED_LENGTH = 0xFFFFFFFF

type SHARE_INFO_2 struct {
	Shi2_netname      *uint16
	Shi2_type         uint32
	Shi2_remark       *uint16
	Shi2_permissions  uint32
	Shi2_max_uses     uint32
	Shi2_current_uses uint32
	Shi2_path         *uint16
	Shi2_passwd       *uint16
}

func ListNetworkShares() ([]so.NetworkShare, error) {
	var (
		dataPointer  uintptr
		resumeHandle uintptr
		entriesRead  uint32
		entriesTotal uint32
		sizeTest     SHARE_INFO_2
		retVal       = make([]so.NetworkShare, 0)
	)
	ret, _, _ := shNetShareEnum.Call(
		uintptr(0),
		uintptr(uint32(2)), // SHARE_INFO_2
		uintptr(unsafe.Pointer(&dataPointer)),
		uintptr(uint32(SHARE_MAX_PREFERRED_LENGTH)),
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&entriesTotal)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)
	if ret != NET_API_STATUS_NERR_Success {
		return nil, fmt.Errorf("error fetching network shares")
	} else if dataPointer == uintptr(0) {
		return nil, fmt.Errorf("null poinnter while fetching entry")
	}
	var iter = dataPointer
	for i := uint32(0); i < entriesRead; i++ {
		var data = (*SHARE_INFO_2)(unsafe.Pointer(iter))
		sd := so.NetworkShare{
			Name:        UTF16toString(data.Shi2_netname),
			Comment:     UTF16toString(data.Shi2_remark),
			Permissions: data.Shi2_permissions,
			MaxUses:     data.Shi2_max_uses,
			CurrentUses: data.Shi2_current_uses,
			Path:        UTF16toString(data.Shi2_path),
		}
		retVal = append(retVal, sd)
		iter = uintptr(unsafe.Pointer(iter + unsafe.Sizeof(sizeTest)))
	}
	usrNetApiBufferFree.Call(dataPointer)
	return retVal, nil
}
