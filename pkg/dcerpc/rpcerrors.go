package dcerpc

// 此文件提供rpc状态信息

// status codes, references:
// https://docs.microsoft.com/windows/desktop/Rpc/rpc-return-values
// https://msdn.microsoft.com/library/default.asp?url=/library/en-us/randz/protocol/common_return_values.asp
// winerror.h
// https://www.opengroup.org/onlinepubs/9629399/apdxn.htm
// https://learn.microsoft.com/en-us/windows/win32/rpc/rpc-return-values
// https://learn.microsoft.com/en-us/previous-versions/aa505946(v=msdn.10)

const (
	EPT_S_CANT_CREATE            = 0x0000076B
	EPT_S_CANT_PERFORM_OP        = 1752
	EPT_S_INVALID_ENTRY          = 0x000006D7
	EPT_S_NOT_REGISTERED         = 0x000006D9
	RPC_S_ACCESS_DENIED          = 0x00000005
	RPC_S_ADDRESS_ERROR          = 0x000006E8
	RPC_S_ALREADY_LISTENING      = 0x000006B1
	RPC_S_ALREADY_REGISTERED     = 0x000006AF
	RPC_S_ASYNC_CALL_PENDING     = 0x000003E5
	RPC_S_BINDING_HAS_NO_AUTH    = 0x000006D2
	RPC_S_BINDING_INCOMPLETE     = 0x0000071B
	RPC_S_BUFFER_TOO_SMALL       = 0x0000007A
	RPC_S_CALL_CANCELLED         = 0x0000071A
	RPC_S_CALL_FAILED            = 0x000006BE
	RPC_S_CALL_FAILED_DNE        = 0x000006BF
	RPC_S_CALL_IN_PROGRESS       = 0x000006FF
	RPC_S_CANNOT_SUPPORT         = 0x000006E4
	RPC_S_CANT_CREATE_ENDPOINT   = 0x000006B8
	RPC_S_COMM_FAILURE           = 0x0000071C
	RPC_S_DUPLICATE_ENDPOINT     = 0x000006CC
	RPC_S_ENTRY_ALREADY_EXISTS   = 0x000006E0
	RPC_S_ENTRY_NOT_FOUND        = 0x000006E1
	RPC_S_FP_DIV_ZERO            = 0x000006E9
	RPC_S_FP_OVERFLOW            = 0x000006EA
	RPC_S_FP_UNDERFLOW           = 0x000006EB
	RPC_S_GROUP_MEMBER_NOT_FOUND = 0x0000076A
	RPC_S_INCOMPLETE_NAME        = 0x000006DB
	RPC_S_INTERFACE_NOT_FOUND    = 0x000006DF
	RPC_S_INTERNAL_ERROR         = 0x000006E6
	//RPC_S_INVALID_ARG             = 0x00000057
	RPC_S_INVALID_AUTH_IDENTITY = 0x000006D5
	RPC_S_INVALID_BINDING       = 0x000006A6
	//RPC_S_INVALID_BOUND           = 0x000006C6
	RPC_S_INVALID_ENDPOINT_FORMAT = 0x000006AA
	RPC_S_INVALID_LEVEL           = 0x00000057
	RPC_S_INVALID_NAF_ID          = 0x000006E3
	RPC_S_INVALID_NAME_SYNTAX     = 0x000006C8
	RPC_S_INVALID_NET_ADDR        = 0x000006AB
	RPC_S_INVALID_NETWORK_OPTIONS = 0x000006BC
	RPC_S_INVALID_OBJECT          = 0x0000076C
	RPC_S_INVALID_RPC_PROTSEQ     = 0x000006A8
	RPC_S_INVALID_SECURITY_DESC   = 0x0000053A
	RPC_S_INVALID_STRING_BINDING  = 0x000006A4
	RPC_S_INVALID_STRING_UUID     = 0x000006A9
	//RPC_S_INVALID_TAG              = 0x000006C5
	RPC_S_INVALID_VERS_OPTION      = 0x000006DC
	RPC_S_MAX_CALLS_TOO_SMALL      = 0x000006CE
	RPC_S_NAME_SERVICE_UNAVAILABLE = 0x000006E2
	RPC_S_NO_BINDINGS              = 0x000006B6
	RPC_S_NO_CALL_ACTIVE           = 0x000006BD
	RPC_S_NO_CONTEXT_AVAILABLE     = 0x000006E5
	RPC_S_NO_ENDPOINT_FOUND        = 0x000006AC
	RPC_S_NO_ENTRY_NAME            = 0x000006C7
	RPC_S_NO_ENV_SETUP             = 0x16c9a0c4
	RPC_S_NO_INTERFACES            = 0x00000719
	RPC_S_NO_INTERFACES_EXPORTED   = 0x16c9a0b1
	RPC_S_NO_MORE_BINDINGS         = 0x0000070E
	RPC_S_NO_MORE_ELEMENTS         = 0x16c9a0a7
	RPC_S_NO_MORE_MEMBERS          = 0x000006DD
	//RPC_S_NO_NS_PRIVILEGE          = 0x00000005
	RPC_S_NO_PRINC_NAME           = 0x0000071E
	RPC_S_NO_PROTSEQS             = 0x000006B7
	RPC_S_NO_PROTSEQS_REGISTERED  = 0x000006B2
	RPC_S_NOT_ALL_OBJS_UNEXPORTED = 0x000006DE
	RPC_S_NOT_CANCELLED           = 0x00000722
	RPC_S_NOT_LISTENING           = 0x000006B3
	RPC_S_NOT_RPC_ERROR           = 0x0000071F
	RPC_S_NOTHING_TO_EXPORT       = 0x000006DA
	RPC_S_OBJECT_NOT_FOUND        = 0x000006AE
	RPC_S_OK                      = 0x00000000
	//RPC_S_OUT_OF_MEMORY            = 0x0000000E
	RPC_S_OUT_OF_RESOURCES        = 0x000006B9
	RPC_S_OUT_OF_THREADS          = 0x000000A4
	RPC_S_PROCNUM_OUT_OF_RANGE    = 0x000006D1
	RPC_S_PROTOCOL_ERROR          = 0x000006C0
	RPC_S_PROTSEQ_NOT_FOUND       = 0x000006D0
	RPC_S_PROTSEQ_NOT_SUPPORTED   = 0x000006A7
	RPC_S_SEC_PKG_ERROR           = 0x00000721
	RPC_S_SERVER_OUT_OF_MEMORY    = 0x0000046A
	RPC_S_SERVER_TOO_BUSY         = 0x000006BB
	RPC_S_SERVER_UNAVAILABLE      = 0x000006BA
	RPC_S_STRING_TOO_LONG         = 0x000006CF
	RPC_S_TYPE_ALREADY_REGISTERED = 0x000006B0
	RPC_S_UNKNOWN_AUTHN_LEVEL     = 0x000006D4
	RPC_S_UNKNOWN_AUTHN_SERVICE   = 0x000006D3
	RPC_S_UNKNOWN_AUTHN_TYPE      = 0x000006CD
	RPC_S_UNKNOWN_AUTHZ_SERVICE   = 0x000006D6
	RPC_S_UNKNOWN_IF              = 0x000006B5
	RPC_S_UNKNOWN_MGR_TYPE        = 0x000006B4
	RPC_S_UNSUPPORTED_AUTHN_LEVEL = 0x0000071D
	RPC_S_UNKNOWN_PRINCIPAL       = 0x00000534
	RPC_S_UNSUPPORTED_NAME_SYNTAX = 0x000006C9
	RPC_S_UNSUPPORTED_TRANS_SYN   = 0x000006C2
	RPC_S_UNSUPPORTED_TYPE        = 0x000006C4
	RPC_S_UUID_LOCAL_ONLY         = 0x00000720
	RPC_S_UUID_NO_ADDRESS         = 0x000006CB
	RPC_S_WRONG_KIND_OF_BINDING   = 0x000006A5
	RPC_S_ZERO_DIVIDE             = 0x000006E7
	RPC_X_BAD_STUB_DATA           = 0x000006F7
	RPC_X_BYTE_COUNT_TOO_SMAL     = 0x000006F6
	//RPC_X_ENUM_VALUE_OUT_OF_RANGE = 0x000006F5
	RPC_X_ENUM_VALUE_TOO_LARGE = 0x000006F5
	RPC_X_INVALID_BOUND        = 0x000006C6
	//RPC_X_INVALID_BUFFER          = 0x000006F8
	RPC_X_INVALID_PIPE_OPERATION = 0x00000727
	RPC_X_INVALID_TAG            = 0x000006C5
	RPC_S_INVALID_TIMEOUT        = 0x000006AD
	RPC_X_NO_MEMORY              = 0x0000000E
	RPC_X_NO_MORE_ENTRIES        = 0x000006EC
	RPC_X_NULL_REF_POINTER       = 0x000006F4
	RPC_X_PIPE_APP_MEMORY        = 0x0000000E
	//RPC_X_SS_BAD_ES_VERSION=
	RPC_X_SS_CANNOT_GET_CALL_HANDLE = 0x000006F3
	RPC_X_SS_CHAR_TRANS_OPEN_FAIL   = 0x000006ED
	RPC_X_SS_CHAR_TRANS_SHORT_FILE  = 0x000006EE
	RPC_X_SS_CONTEXT_DAMAGED        = 0x000006F1
	RPC_X_SS_CONTEXT_MISMATCH       = 0x00000006
	RPC_X_SS_HANDLES_MISMATCH       = 0x000006F2
	RPC_X_SS_IN_NULL_CONTEXT        = 0x000006EF
	RPC_X_SS_INVALID_BUFFER         = 0x000006F8
	RPC_X_SS_WRONG_ES_VERSION       = 0x00000724
	RPC_X_SS_WRONG_STUB_VERSION     = 0x00000725
)

var RpcStatusCodes = map[uint32]string{
	EPT_S_CANT_CREATE:            "An entry into the endpoint mapper database cannot be created.",
	EPT_S_CANT_PERFORM_OP:        "General failure when trying to perform an operation on the endpoint mapper database.",
	EPT_S_INVALID_ENTRY:          "The specified endpoint mapper database entry is invalid.",
	EPT_S_NOT_REGISTERED:         "There are no more endpoints available from the endpoint-map database.",
	RPC_S_ACCESS_DENIED:          "Access for making the remote procedure call was denied.",
	RPC_S_ADDRESS_ERROR:          "An addressing error has occurred on the server.",
	RPC_S_ALREADY_LISTENING:      "The server is already listening.",
	RPC_S_ALREADY_REGISTERED:     "The object UUID has already been registered.",
	RPC_S_ASYNC_CALL_PENDING:     "The asynchronous remote procedure call has not yet completed.",
	RPC_S_BINDING_HAS_NO_AUTH:    "The binding does not contain any authentication information.",
	RPC_S_BINDING_INCOMPLETE:     "Not all required elements from the binding handle were supplied.",
	RPC_S_BUFFER_TOO_SMALL:       "The buffer given to RPC by the caller is too small.",
	RPC_S_CALL_CANCELLED:         "The remote procedure call was canceled, or if a call time out was specified, the call timed out.",
	RPC_S_CALL_FAILED:            "The remote procedure call failed. Implies the server was reachable at a certain point in time, and execution of the remote procedure call on the server may have started.",
	RPC_S_CALL_FAILED_DNE:        "The remote procedure call failed, and execution on the server did not start. Implies the server was reachable at a certain point in time.",
	RPC_S_CALL_IN_PROGRESS:       "A remote procedure call is still in progress.",
	RPC_S_CANNOT_SUPPORT:         "The requested operation is not supported.",
	RPC_S_CANT_CREATE_ENDPOINT:   "The endpoint cannot be created.",
	RPC_S_COMM_FAILURE:           "Unable to communicate with the server.",
	RPC_S_DUPLICATE_ENDPOINT:     "The endpoint is a duplicate.",
	RPC_S_ENTRY_ALREADY_EXISTS:   "The entry already exists.",
	RPC_S_ENTRY_NOT_FOUND:        "The entry is not found.",
	RPC_S_FP_DIV_ZERO:            "A floating-point operation at the server has caused a divide by zero.",
	RPC_S_FP_OVERFLOW:            "A floating-point overflow has occurred at the server.",
	RPC_S_FP_UNDERFLOW:           "A floating-point underflow has occurred at the server.",
	RPC_S_GROUP_MEMBER_NOT_FOUND: "The group member has not been found.",
	RPC_S_INCOMPLETE_NAME:        "The entry name is incomplete.",
	RPC_S_INTERFACE_NOT_FOUND:    "The interface has not been found.",
	RPC_S_INTERNAL_ERROR:         "An internal error has occurred in a remote procedure call.",
	//RPC_S_INVALID_ARG:               "The specified argument is not valid.",
	RPC_S_INVALID_AUTH_IDENTITY: "The specified authentication identity could not be used. For example an LRPC client stopped functioning in the middle of an RPC and the server could not impersonate it. Or, credentials for a client could not be acquired by the security provider.",
	RPC_S_INVALID_BINDING:       "The binding handle is invalid.",
	//RPC_S_INVALID_BOUND:             "The array bounds are invalid.",
	RPC_S_INVALID_ENDPOINT_FORMAT: "The endpoint format is invalid.",
	RPC_S_INVALID_LEVEL:           "The version, level, or flags parameter is invalid.",
	RPC_S_INVALID_NAF_ID:          "The network-address family is invalid.",
	RPC_S_INVALID_NAME_SYNTAX:     "The name syntax is invalid.",
	RPC_S_INVALID_NET_ADDR:        "The network address is invalid.",
	RPC_S_INVALID_NETWORK_OPTIONS: "The network options are invalid.",
	RPC_S_INVALID_OBJECT:          "The object is invalid.",
	RPC_S_INVALID_RPC_PROTSEQ:     "The RPC protocol sequence is invalid.",
	RPC_S_INVALID_SECURITY_DESC:   "The security descriptor is not in the valid format.",
	RPC_S_INVALID_STRING_BINDING:  "The string binding is invalid.",
	RPC_S_INVALID_STRING_UUID:     "The string UUID is invalid.",
	//RPC_S_INVALID_TAG:              "The discriminant value does not match any of the case values. There is no default case.",
	RPC_S_INVALID_TIMEOUT:          "The time-out value is invalid.",
	RPC_S_INVALID_VERS_OPTION:      "The version option is invalid.",
	RPC_S_MAX_CALLS_TOO_SMALL:      "The maximum number of calls is too small.",
	RPC_S_NAME_SERVICE_UNAVAILABLE: "The name service is unavailable.",
	RPC_S_NO_BINDINGS:              "There are no bindings.",
	RPC_S_NO_CALL_ACTIVE:           "There is no remote procedure call active in this thread.",
	RPC_S_NO_CONTEXT_AVAILABLE:     "No security context is available to allow impersonation.",
	RPC_S_NO_ENDPOINT_FOUND:        "No endpoint has been found.",
	RPC_S_NO_ENTRY_NAME:            "The binding does not contain an entry name.",
	RPC_S_NO_ENV_SETUP:             "No environment variable is set up.",
	RPC_S_NO_INTERFACES:            "No interfaces are registered.",
	RPC_S_NO_INTERFACES_EXPORTED:   "No interfaces have been exported.",
	RPC_S_NO_MORE_BINDINGS:         "There are no more bindings.",
	RPC_S_NO_MORE_ELEMENTS:         "There are no more elements.",
	RPC_S_NO_MORE_MEMBERS:          "There are no more members.",
	//RPC_S_NO_NS_PRIVILEGE:           "There is no privilege for a name-service operation.",
	RPC_S_NO_PRINC_NAME:           "No principal name is registered.",
	RPC_S_NO_PROTSEQS:             "There are no protocol sequences.",
	RPC_S_NO_PROTSEQS_REGISTERED:  "No protocol sequences have been registered.",
	RPC_S_NOT_ALL_OBJS_UNEXPORTED: "Not all objects are unexported.",
	RPC_S_NOT_CANCELLED:           "The thread is not canceled.",
	RPC_S_NOT_LISTENING:           "The server is not listening.",
	RPC_S_NOT_RPC_ERROR:           "The status code requested is not valid.",
	RPC_S_NOTHING_TO_EXPORT:       "There is nothing to export.",
	RPC_S_OBJECT_NOT_FOUND:        "The object UUID has not been found.",
	RPC_S_OK:                      "The requested operation completed successfully.",
	//RPC_S_OUT_OF_MEMORY:           "The needed memory is not available.",
	RPC_S_OUT_OF_RESOURCES:        "Not enough resources are available to complete this operation.",
	RPC_S_OUT_OF_THREADS:          "The RPC run-time library was not able to create another thread.",
	RPC_S_PROCNUM_OUT_OF_RANGE:    "The procedure number is out of range.",
	RPC_S_PROTOCOL_ERROR:          "An RPC protocol error has occurred.",
	RPC_S_PROTSEQ_NOT_FOUND:       "The RPC protocol sequence has not been found.",
	RPC_S_PROTSEQ_NOT_SUPPORTED:   "The RPC protocol sequence is not supported.",
	RPC_S_SEC_PKG_ERROR:           "An error that has no RPC mapping was returned by the security package. Retrieve the security provider error using the RPC Extended Error Mechanism.",
	RPC_S_SERVER_OUT_OF_MEMORY:    "The server has insufficient memory to complete this operation.",
	RPC_S_SERVER_TOO_BUSY:         "The server is too busy to complete this operation.",
	RPC_S_SERVER_UNAVAILABLE:      "The server is unavailable.",
	RPC_S_STRING_TOO_LONG:         "The string is too long.",
	RPC_S_TYPE_ALREADY_REGISTERED: "The type UUID has already been registered.",
	RPC_S_UNKNOWN_AUTHN_LEVEL:     "The authentication level is unknown.",
	RPC_S_UNKNOWN_AUTHN_SERVICE:   "The authentication service is unknown.",
	RPC_S_UNKNOWN_AUTHN_TYPE:      "The authentication type is unknown.",
	RPC_S_UNKNOWN_AUTHZ_SERVICE:   "The authorization service is unknown.",
	RPC_S_UNKNOWN_IF:              "The interface is unknown.",
	RPC_S_UNKNOWN_MGR_TYPE:        "The manager type is unknown.",
	RPC_S_UNSUPPORTED_AUTHN_LEVEL: "The authentication level is not supported.",
	RPC_S_UNKNOWN_PRINCIPAL:       "The principal name is not recognized.",
	RPC_S_UNSUPPORTED_NAME_SYNTAX: "The name syntax is not supported.",
	RPC_S_UNSUPPORTED_TRANS_SYN:   "The transfer syntax is not supported by the server.",
	RPC_S_UNSUPPORTED_TYPE:        "The type UUID is not supported.",
	RPC_S_UUID_LOCAL_ONLY:         "A UUID valid only for the local computer has been allocated.",
	RPC_S_UUID_NO_ADDRESS:         "No network address is available for constructing a UUID.",
	RPC_S_WRONG_KIND_OF_BINDING:   "The binding handle is not the correct type.",
	RPC_S_ZERO_DIVIDE:             "The server has attempted an integer divide by zero.",
	RPC_X_BAD_STUB_DATA:           "The stub has received bad data.",
	RPC_X_BYTE_COUNT_TOO_SMAL:     "The byte count is too small.",
	//RPC_X_ENUM_VALUE_OUT_OF_RANGE: "The enumeration value is out of range.",
	RPC_X_ENUM_VALUE_TOO_LARGE: "The enumeration constant must be less than 65535.",
	RPC_X_INVALID_BOUND:        "The specified bounds of an array are inconsistent.",
	//RPC_X_INVALID_BUFFER:          "The pointer does not contain the address of a valid data buffer.",
	RPC_X_INVALID_PIPE_OPERATION: "The requested pipe operation is not supported.",
	RPC_X_INVALID_TAG:            "The discriminant value does not match any of the case values. There is no default case.",
	RPC_X_NO_MEMORY:              "Insufficient memory is available.",
	RPC_X_NO_MORE_ENTRIES:        "The list of servers available for the [auto_handle] binding has been exhausted.",
	RPC_X_NULL_REF_POINTER:       "A null reference pointer has been passed to the stub.",
	//RPC_X_PIPE_APP_MEMORY:         "Insufficient memory is available for pipe data.",
	//RPC_X_SS_BAD_ES_VERSION:         "The operation for the serializing handle is not valid.",
	RPC_X_SS_CANNOT_GET_CALL_HANDLE: "The stub is unable to get the call handle.",
	RPC_X_SS_CHAR_TRANS_OPEN_FAIL:   "The file designated by DCERPCCHARTRANS cannot be opened.",
	RPC_X_SS_CHAR_TRANS_SHORT_FILE:  "The file containing the character-translation table has fewer than 512 bytes.",
	RPC_X_SS_CONTEXT_DAMAGED:        "The context handle changed during a call. Only raised on the client side.",
	RPC_X_SS_CONTEXT_MISMATCH:       "The context handle does not match any known context handles.",
	RPC_X_SS_HANDLES_MISMATCH:       "The binding handles passed to a remote procedure call do not match.",
	RPC_X_SS_IN_NULL_CONTEXT:        "A null context handle is passed in an in parameter position.",
	RPC_X_SS_INVALID_BUFFER:         "The buffer is not valid for the operation.",
	RPC_X_SS_WRONG_ES_VERSION:       "The software version is incorrect.",
	RPC_X_SS_WRONG_STUB_VERSION:     "The stub version is incorrect.",
}
