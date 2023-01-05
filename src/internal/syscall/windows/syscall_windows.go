// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows

import (
	"internal/unsafeheader"
	"sync"
	"syscall"
	"unicode/utf16"
	"unsafe"
	"errors"
)

// UTF16PtrToString is like UTF16ToString, but takes *uint16
// as a parameter instead of []uint16.
func UTF16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	// Find NUL terminator.
	end := unsafe.Pointer(p)
	n := 0
	for *(*uint16)(end) != 0 {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	// Turn *uint16 into []uint16.
	var s []uint16
	hdr := (*unsafeheader.Slice)(unsafe.Pointer(&s))
	hdr.Data = unsafe.Pointer(p)
	hdr.Cap = n
	hdr.Len = n
	// Decode []uint16 into string.
	return string(utf16.Decode(s))
}

const (
	ERROR_SHARING_VIOLATION      syscall.Errno = 32
	ERROR_LOCK_VIOLATION         syscall.Errno = 33
	ERROR_NOT_SUPPORTED          syscall.Errno = 50
	ERROR_CALL_NOT_IMPLEMENTED   syscall.Errno = 120
	ERROR_INVALID_NAME           syscall.Errno = 123
	ERROR_LOCK_FAILED            syscall.Errno = 167
	ERROR_NO_UNICODE_TRANSLATION syscall.Errno = 1113
)

const GAA_FLAG_INCLUDE_PREFIX = 0x00000010

const (
	IF_TYPE_OTHER              = 1
	IF_TYPE_ETHERNET_CSMACD    = 6
	IF_TYPE_ISO88025_TOKENRING = 9
	IF_TYPE_PPP                = 23
	IF_TYPE_SOFTWARE_LOOPBACK  = 24
	IF_TYPE_ATM                = 37
	IF_TYPE_IEEE80211          = 71
	IF_TYPE_TUNNEL             = 131
	IF_TYPE_IEEE1394           = 144
)

type SocketAddress struct {
	Sockaddr       *syscall.RawSockaddrAny
	SockaddrLength int32
}

type IpAdapterUnicastAddress struct {
	Length             uint32
	Flags              uint32
	Next               *IpAdapterUnicastAddress
	Address            SocketAddress
	PrefixOrigin       int32
	SuffixOrigin       int32
	DadState           int32
	ValidLifetime      uint32
	PreferredLifetime  uint32
	LeaseLifetime      uint32
	OnLinkPrefixLength uint8
}

type IpAdapterAnycastAddress struct {
	Length  uint32
	Flags   uint32
	Next    *IpAdapterAnycastAddress
	Address SocketAddress
}

type IpAdapterMulticastAddress struct {
	Length  uint32
	Flags   uint32
	Next    *IpAdapterMulticastAddress
	Address SocketAddress
}

type IpAdapterDnsServerAdapter struct {
	Length   uint32
	Reserved uint32
	Next     *IpAdapterDnsServerAdapter
	Address  SocketAddress
}

type IpAdapterPrefix struct {
	Length       uint32
	Flags        uint32
	Next         *IpAdapterPrefix
	Address      SocketAddress
	PrefixLength uint32
}

type IpAdapterAddresses struct {
	Length                uint32
	IfIndex               uint32
	Next                  *IpAdapterAddresses
	AdapterName           *byte
	FirstUnicastAddress   *IpAdapterUnicastAddress
	FirstAnycastAddress   *IpAdapterAnycastAddress
	FirstMulticastAddress *IpAdapterMulticastAddress
	FirstDnsServerAddress *IpAdapterDnsServerAdapter
	DnsSuffix             *uint16
	Description           *uint16
	FriendlyName          *uint16
	PhysicalAddress       [syscall.MAX_ADAPTER_ADDRESS_LENGTH]byte
	PhysicalAddressLength uint32
	Flags                 uint32
	Mtu                   uint32
	IfType                uint32
	OperStatus            uint32
	Ipv6IfIndex           uint32
	ZoneIndices           [16]uint32
	FirstPrefix           *IpAdapterPrefix
	/* more fields might be present here. */
}

type FILE_BASIC_INFO struct {
	CreationTime   syscall.Filetime
	LastAccessTime syscall.Filetime
	LastWriteTime  syscall.Filetime
	ChangedTime    syscall.Filetime
	FileAttributes uint32
}

const (
	IfOperStatusUp             = 1
	IfOperStatusDown           = 2
	IfOperStatusTesting        = 3
	IfOperStatusUnknown        = 4
	IfOperStatusDormant        = 5
	IfOperStatusNotPresent     = 6
	IfOperStatusLowerLayerDown = 7
)

//sys	GetAdaptersAddresses(family uint32, flags uint32, reserved uintptr, adapterAddresses *IpAdapterAddresses, sizePointer *uint32) (errcode error) = iphlpapi.GetAdaptersAddresses
//sys	GetComputerNameEx(nameformat uint32, buf *uint16, n *uint32) (err error) = GetComputerNameExW
//sys	MoveFileEx(from *uint16, to *uint16, flags uint32) (err error) = MoveFileExW
//sys	GetModuleFileName(module syscall.Handle, fn *uint16, len uint32) (n uint32, err error) = kernel32.GetModuleFileNameW
//sys	SetFileInformationByHandle_orig(handle syscall.Handle, fileInformationClass uint32, buf uintptr, bufsize uint32) (err error) = kernel32.SetFileInformationByHandle
//sys	VirtualQuery(address uintptr, buffer *MemoryBasicInformation, length uintptr) (err error) = kernel32.VirtualQuery

const (
	WSA_FLAG_OVERLAPPED        = 0x01
	WSA_FLAG_NO_HANDLE_INHERIT = 0x80

	WSAEMSGSIZE syscall.Errno = 10040

	MSG_PEEK   = 0x2
	MSG_TRUNC  = 0x0100
	MSG_CTRUNC = 0x0200

	socket_error = uintptr(^uint32(0))
)

var WSAID_WSASENDMSG = syscall.GUID{
	Data1: 0xa441e712,
	Data2: 0x754f,
	Data3: 0x43ca,
	Data4: [8]byte{0x84, 0xa7, 0x0d, 0xee, 0x44, 0xcf, 0x60, 0x6d},
}

var WSAID_WSARECVMSG = syscall.GUID{
	Data1: 0xf689d7c8,
	Data2: 0x6f1f,
	Data3: 0x436b,
	Data4: [8]byte{0x8a, 0x53, 0xe5, 0x4f, 0xe3, 0x51, 0xc3, 0x22},
}

var sendRecvMsgFunc struct {
	once     sync.Once
	sendAddr uintptr
	recvAddr uintptr
	err      error
}

type WSAMsg struct {
	Name        syscall.Pointer
	Namelen     int32
	Buffers     *syscall.WSABuf
	BufferCount uint32
	Control     syscall.WSABuf
	Flags       uint32
}

//sys	WSASocket(af int32, typ int32, protocol int32, protinfo *syscall.WSAProtocolInfo, group uint32, flags uint32) (handle syscall.Handle, err error) [failretval==syscall.InvalidHandle] = ws2_32.WSASocketW

func loadWSASendRecvMsg() error {
	sendRecvMsgFunc.once.Do(func() {
		var s syscall.Handle
		s, sendRecvMsgFunc.err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
		if sendRecvMsgFunc.err != nil {
			return
		}
		defer syscall.CloseHandle(s)
		var n uint32
		sendRecvMsgFunc.err = syscall.WSAIoctl(s,
			syscall.SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_WSARECVMSG)),
			uint32(unsafe.Sizeof(WSAID_WSARECVMSG)),
			(*byte)(unsafe.Pointer(&sendRecvMsgFunc.recvAddr)),
			uint32(unsafe.Sizeof(sendRecvMsgFunc.recvAddr)),
			&n, nil, 0)
		if sendRecvMsgFunc.err != nil {
			return
		}
		sendRecvMsgFunc.err = syscall.WSAIoctl(s,
			syscall.SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_WSASENDMSG)),
			uint32(unsafe.Sizeof(WSAID_WSASENDMSG)),
			(*byte)(unsafe.Pointer(&sendRecvMsgFunc.sendAddr)),
			uint32(unsafe.Sizeof(sendRecvMsgFunc.sendAddr)),
			&n, nil, 0)
	})
	return sendRecvMsgFunc.err
}

func WSASendMsg(fd syscall.Handle, msg *WSAMsg, flags uint32, bytesSent *uint32, overlapped *syscall.Overlapped, croutine *byte) error {
	err := loadWSASendRecvMsg()
	if err != nil {
		return err
	}
	r1, _, e1 := syscall.Syscall6(sendRecvMsgFunc.sendAddr, 6, uintptr(fd), uintptr(unsafe.Pointer(msg)), uintptr(flags), uintptr(unsafe.Pointer(bytesSent)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)))
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func WSARecvMsg(fd syscall.Handle, msg *WSAMsg, bytesReceived *uint32, overlapped *syscall.Overlapped, croutine *byte) error {
	err := loadWSASendRecvMsg()
	if err != nil {
		return err
	}
	r1, _, e1 := syscall.Syscall6(sendRecvMsgFunc.recvAddr, 5, uintptr(fd), uintptr(unsafe.Pointer(msg)), uintptr(unsafe.Pointer(bytesReceived)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

const (
	ComputerNameNetBIOS                   = 0
	ComputerNameDnsHostname               = 1
	ComputerNameDnsDomain                 = 2
	ComputerNameDnsFullyQualified         = 3
	ComputerNamePhysicalNetBIOS           = 4
	ComputerNamePhysicalDnsHostname       = 5
	ComputerNamePhysicalDnsDomain         = 6
	ComputerNamePhysicalDnsFullyQualified = 7
	ComputerNameMax                       = 8

	MOVEFILE_REPLACE_EXISTING      = 0x1
	MOVEFILE_COPY_ALLOWED          = 0x2
	MOVEFILE_DELAY_UNTIL_REBOOT    = 0x4
	MOVEFILE_WRITE_THROUGH         = 0x8
	MOVEFILE_CREATE_HARDLINK       = 0x10
	MOVEFILE_FAIL_IF_NOT_TRACKABLE = 0x20
)

func Rename(oldpath, newpath string) error {
	from, err := syscall.UTF16PtrFromString(oldpath)
	if err != nil {
		return err
	}
	to, err := syscall.UTF16PtrFromString(newpath)
	if err != nil {
		return err
	}
	return MoveFileEx(from, to, MOVEFILE_REPLACE_EXISTING)
}

//sys LockFileEx(file syscall.Handle, flags uint32, reserved uint32, bytesLow uint32, bytesHigh uint32, overlapped *syscall.Overlapped) (err error) = kernel32.LockFileEx
//sys UnlockFileEx(file syscall.Handle, reserved uint32, bytesLow uint32, bytesHigh uint32, overlapped *syscall.Overlapped) (err error) = kernel32.UnlockFileEx

const (
	LOCKFILE_FAIL_IMMEDIATELY = 0x00000001
	LOCKFILE_EXCLUSIVE_LOCK   = 0x00000002
)

const MB_ERR_INVALID_CHARS = 8

//sys	GetACP() (acp uint32) = kernel32.GetACP
//sys	GetConsoleCP() (ccp uint32) = kernel32.GetConsoleCP
//sys	MultiByteToWideChar(codePage uint32, dwFlags uint32, str *byte, nstr int32, wchar *uint16, nwchar int32) (nwrite int32, err error) = kernel32.MultiByteToWideChar
//sys	GetCurrentThread() (pseudoHandle syscall.Handle, err error) = kernel32.GetCurrentThread

const STYPE_DISKTREE = 0x00

type SHARE_INFO_2 struct {
	Netname     *uint16
	Type        uint32
	Remark      *uint16
	Permissions uint32
	MaxUses     uint32
	CurrentUses uint32
	Path        *uint16
	Passwd      *uint16
}

//sys  NetShareAdd(serverName *uint16, level uint32, buf *byte, parmErr *uint16) (neterr error) = netapi32.NetShareAdd
//sys  NetShareDel(serverName *uint16, netName *uint16, reserved uint32) (neterr error) = netapi32.NetShareDel

const (
	FILE_NAME_NORMALIZED = 0x0
	FILE_NAME_OPENED     = 0x8

	VOLUME_NAME_DOS  = 0x0
	VOLUME_NAME_GUID = 0x1
	VOLUME_NAME_NONE = 0x4
	VOLUME_NAME_NT   = 0x2
)

//sys	GetFinalPathNameByHandle(file syscall.Handle, filePath *uint16, filePathSize uint32, flags uint32) (n uint32, err error) = kernel32.GetFinalPathNameByHandleW

func LoadGetFinalPathNameByHandle() error {
	return procGetFinalPathNameByHandleW.Find()
}

//sys	CreateEnvironmentBlock(block **uint16, token syscall.Token, inheritExisting bool) (err error) = userenv.CreateEnvironmentBlock
//sys	DestroyEnvironmentBlock(block *uint16) (err error) = userenv.DestroyEnvironmentBlock

//sys	RtlGenRandom(buf []byte) (err error) = advapi32.SystemFunction036


//BACKPORT(NT_51): Make SetFileInformationByHandle compatible for 
// Windows XP and eventually needs expansion once it is more utilised.
//sys	NtSetInformationFile(handle syscall.Handle, iosb *syscall.IO_STATUS_BLOCK, inBuffer *byte, inBufferLen uint32, class uint32) (ntstatus syscall.NTStatus) = ntdll.NtSetInformationFile

func LoadSetFileInformationByHandle() error {
	return procSetFileInformationByHandle.Find()
}

func SetFileInformationByHandle(handle syscall.Handle, fileInformationClass uint32, buf uintptr, bufsize uint32) (err error) {
	if LoadSetFileInformationByHandle() != nil {
		return CustomSetFileInformationByHandle(handle, fileInformationClass, buf, bufsize)
	}
	return SetFileInformationByHandle_orig(handle, fileInformationClass, buf, bufsize)
}

// FILE_INFO_BY_HANDLE_CLASS constants for SetFileInformationByHandle/GetFileInformationByHandleEx
const (
	CustomFileBasicInfo                  = 0
	CustomFileStandardInfo               = 1
	CustomFileNameInfo                   = 2
	CustomFileRenameInfo                 = 3
	CustomFileDispositionInfo            = 4
	CustomFileAllocationInfo             = 5
	CustomFileEndOfFileInfo              = 6
	CustomFileStreamInfo                 = 7
	CustomFileCompressionInfo            = 8
	CustomFileAttributeTagInfo           = 9
	CustomFileIdBothDirectoryInfo        = 10
	CustomFileIdBothDirectoryRestartInfo = 11
	CustomFileIoPriorityHintInfo         = 12
	CustomFileRemoteProtocolInfo         = 13
	CustomFileFullDirectoryInfo          = 14
	CustomFileFullDirectoryRestartInfo   = 15
	CustomFileStorageInfo                = 16
	CustomFileAlignmentInfo              = 17
	CustomFileIdInfo                     = 18
	CustomFileIdExtdDirectoryInfo        = 19
	CustomFileIdExtdDirectoryRestartInfo = 20
	CustomFileDispositionInfoEx          = 21
	CustomFileRenameInfoEx               = 22
	CustomFileCaseSensitiveInfo          = 23
	CustomFileNormalizedNameInfo         = 24
)

const (
	// FileInformationClass for NtSetInformationFile
	CustomFileBasicInformation                         = 4
	CustomFileRenameInformation                        = 10
	CustomFileDispositionInformation                   = 13
	CustomFilePositionInformation                      = 14
	CustomFileEndOfFileInformation                     = 20
	CustomFileValidDataLengthInformation               = 39
	CustomFileShortNameInformation                     = 40
	CustomFileIoPriorityHintInformation                = 43
	CustomFileReplaceCompletionInformation             = 61
	CustomFileDispositionInformationEx                 = 64
	CustomFileCaseSensitiveInformation                 = 71
	CustomFileLinkInformation                          = 72
	CustomFileCaseSensitiveInformationForceAccessCheck = 75
	CustomFileKnownFolderInformation                   = 76
/*
	// Flags for FILE_RENAME_INFORMATION
	CustomFILE_RENAME_REPLACE_IF_EXISTS                    = 0x00000001
	CustomFILE_RENAME_POSIX_SEMANTICS                      = 0x00000002
	CustomFILE_RENAME_SUPPRESS_PIN_STATE_INHERITANCE       = 0x00000004
	CustomFILE_RENAME_SUPPRESS_STORAGE_RESERVE_INHERITANCE = 0x00000008
	CustomFILE_RENAME_NO_INCREASE_AVAILABLE_SPACE          = 0x00000010
	CustomFILE_RENAME_NO_DECREASE_AVAILABLE_SPACE          = 0x00000020
	CustomFILE_RENAME_PRESERVE_AVAILABLE_SPACE             = 0x00000030
	CustomFILE_RENAME_IGNORE_READONLY_ATTRIBUTE            = 0x00000040
	CustomFILE_RENAME_FORCE_RESIZE_TARGET_SR               = 0x00000080
	CustomFILE_RENAME_FORCE_RESIZE_SOURCE_SR               = 0x00000100
	CustomFILE_RENAME_FORCE_RESIZE_SR                      = 0x00000180

	// Flags for FILE_DISPOSITION_INFORMATION_EX
	CustomFILE_DISPOSITION_DO_NOT_DELETE             = 0x00000000
	CustomFILE_DISPOSITION_DELETE                    = 0x00000001
	CustomFILE_DISPOSITION_POSIX_SEMANTICS           = 0x00000002
	CustomFILE_DISPOSITION_FORCE_IMAGE_SECTION_CHECK = 0x00000004
	CustomFILE_DISPOSITION_ON_CLOSE                  = 0x00000008
	CustomFILE_DISPOSITION_IGNORE_READONLY_ATTRIBUTE = 0x00000010

	// Flags for FILE_CASE_SENSITIVE_INFORMATION
	CustomFILE_CS_FLAG_CASE_SENSITIVE_DIR = 0x00000001

	// Flags for FILE_LINK_INFORMATION
	CustomFILE_LINK_REPLACE_IF_EXISTS                    = 0x00000001
	CustomFILE_LINK_POSIX_SEMANTICS                      = 0x00000002
	CustomFILE_LINK_SUPPRESS_STORAGE_RESERVE_INHERITANCE = 0x00000008
	CustomFILE_LINK_NO_INCREASE_AVAILABLE_SPACE          = 0x00000010
	CustomFILE_LINK_NO_DECREASE_AVAILABLE_SPACE          = 0x00000020
	CustomFILE_LINK_PRESERVE_AVAILABLE_SPACE             = 0x00000030
	CustomFILE_LINK_IGNORE_READONLY_ATTRIBUTE            = 0x00000040
	CustomFILE_LINK_FORCE_RESIZE_TARGET_SR               = 0x00000080
	CustomFILE_LINK_FORCE_RESIZE_SOURCE_SR               = 0x00000100
	CustomFILE_LINK_FORCE_RESIZE_SR                      = 0x00000180
*/
)

//BACKPORT(NT_51): Reimplement SetFileInformationByHandle with NtSetInformationFile
//Source: https://source.winehq.org/git/wine.git/blob/17e5ff74308f41ab662d46f684db2c6023a4a16b:/dlls/kernelbase/file.c#l3554
func CustomSetFileInformationByHandle(handle syscall.Handle, fileInformationClass uint32, buf uintptr, bufsize uint32) (err error) {
    var status syscall.NTStatus
	var io syscall.IO_STATUS_BLOCK

	buf_asbyteptr := (*byte)(unsafe.Pointer(buf))
	
	switch (fileInformationClass) {
		case CustomFileNameInfo: fallthrough
		case CustomFileAllocationInfo: fallthrough
		case CustomFileStreamInfo: fallthrough
		case CustomFileIdBothDirectoryInfo: fallthrough
		case CustomFileIdBothDirectoryRestartInfo: fallthrough
		case CustomFileFullDirectoryInfo: fallthrough
		case CustomFileFullDirectoryRestartInfo: fallthrough
		case CustomFileStorageInfo: fallthrough
		case CustomFileAlignmentInfo: fallthrough
		case CustomFileIdInfo: fallthrough
		case CustomFileIdExtdDirectoryInfo: fallthrough
		case CustomFileIdExtdDirectoryRestartInfo:
			println("SetFileInformationByHandle: not implemented", handle, " - ", fileInformationClass)
			return errors.New("SetFileInformationByHandle: not implemented")
		
		case CustomFileEndOfFileInfo:
			status = NtSetInformationFile( handle, &io, buf_asbyteptr, bufsize, CustomFileEndOfFileInformation )

		case CustomFileBasicInfo: //this seems to be the only called thing as of writing this
			status = NtSetInformationFile( handle, &io, buf_asbyteptr, bufsize, CustomFileBasicInformation )

		case CustomFileDispositionInfo:
			status = NtSetInformationFile( handle, &io, buf_asbyteptr, bufsize, CustomFileDispositionInformation )

		case CustomFileIoPriorityHintInfo:
			status = NtSetInformationFile( handle, &io, buf_asbyteptr, bufsize, CustomFileIoPriorityHintInformation )

		case CustomFileRenameInfo:
			println("SetFileInformationByHandle: FileRenameInfo commented out but code template exists" , handle)
			return errors.New("SetFileInformationByHandle: FileRenameInfo commented out but code template exists")
		/*
			FILE_RENAME_INFORMATION *rename_info;
			UNICODE_STRING nt_name;
			ULONG size;
		
			if ((status = RtlDosPathNameToNtPathName_U_WithStatus( ((FILE_RENAME_INFORMATION *)buf)->FileName, &nt_name, NULL, NULL ))) {
				break;
			}
		
			size = sizeof(*rename_info) + nt_name.Length;
			if ((rename_info = HeapAlloc( GetProcessHeap(), 0, size ))) {
				memcpy( rename_info, buf, sizeof(*rename_info) );
				memcpy( rename_info->FileName, nt_name.Buffer, nt_name.Length + sizeof(WCHAR) );
				rename_info->FileNameLength = nt_name.Length;
				status = NtSetInformationFile( file, &io, rename_info, size, FileRenameInformation );
				HeapFree( GetProcessHeap(), 0, rename_info );
			}
			RtlFreeUnicodeString( &nt_name );
			break;
		*/
		case CustomFileStandardInfo: fallthrough
		case CustomFileCompressionInfo: fallthrough
		case CustomFileAttributeTagInfo: fallthrough
		case CustomFileRemoteProtocolInfo: fallthrough
		default:
			println("SetFileInformationByHandle: ERROR_INVALID_PARAMETER" , handle, " - ", fileInformationClass)
			return errors.New("SetFileInformationByHandle: ERROR_INVALID_PARAMETER")
	}

	if status != 0 {
		return errors.New("SetFileInformationByHandle: Failed")
	} else {
		return nil
	}
}
