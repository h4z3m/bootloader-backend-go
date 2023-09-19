package main

import (
	"encoding/json"
)

type CommonResponseData struct {
	Status    bool     `json:"status"`
	ErrorCode uint8    `json:"error"`
	Errors    []string `json:"errors"`
}

type BL_ReadVersionRequest struct {
	CommandID uint8 `json:"commandId"`
}

type BL_ReadVersionResponse struct {
	CommonResponseData
	Version uint8 `json:"version"`
}

type BL_FlashEraseRequest struct {
	CommandID    uint8  `json:"commandId"`
	StartAddress uint32 `json:"address"`
	PageCount    uint32 `json:"count"`
}

type BL_FlashEraseResponse struct {
	CommonResponseData
}

type BL_FlashWriteRequest struct {
	CommandID    uint8  `json:"commandId"`
	StartAddress uint32 `json:"address"`
	Data         []byte `json:"binaryData"`
	Size         uint32 `json:"size"`
}

type BL_FlashWriteResponse struct {
	CommonResponseData
}

type BL_JumpToAppRequest struct {
	CommandID uint8 `json:"commandId"`
}

type BL_JumpToAppResponse struct {
	CommonResponseData
}

type BL_FlashReadRequest struct {
	CommandID    uint8  `json:"commandId"`
	StartAddress uint32 `json:"address"`
	Length       uint32 `json:"length"`
}

type BL_FlashReadResponse struct {
	Data []byte `json:"binaryData"`
	CommonResponseData
}

const (
	BL_MEM_WRITE_CMD_ID      = 2
	BL_MEM_READ_CMD_ID       = 3
	BL_VER_CMD_ID            = 4
	BL_FLASH_ERASE_CMD_ID    = 5
	BL_ENTER_CMD_MODE_CMD_ID = 7
	BL_JUMP_TO_APP_CMD_ID    = 8
)

type BL_NACK_t uint8

const (
	BL_NACK_SUCCESS           BL_NACK_t = 0
	BL_NACK_INVALID_CMD       BL_NACK_t = 1 << 0
	BL_NACK_INVALID_KEY       BL_NACK_t = 1 << 1
	BL_NACK_INVALID_ADDRESS   BL_NACK_t = 1 << 2
	BL_NACK_INVALID_LENGTH    BL_NACK_t = 1 << 3
	BL_NACK_INVALID_DATA      BL_NACK_t = 1 << 4
	BL_NACK_INVALID_CRC       BL_NACK_t = 1 << 5
	BL_NACK_OPERATION_FAILURE BL_NACK_t = 1 << 6
)

// GetErrorStringData generates an array of error strings based on the given BL_NACK_t fields.
//
// The function takes a BL_NACK_t type as the parameter named "fields". It loops through the BL_NACK_t values from 1 to 8
// and checks if the corresponding bit is set in the "fields" parameter. If the bit is set, it appends the corresponding
// error string to the "strArray". The function returns the generated error string array and nil as the error.
func GetErrorStringData(fields BL_NACK_t) []string {
	if fields == BL_NACK_SUCCESS {
		// No strings to generate, return nil and error
		return nil
	}

	strArray := make([]string, 0)

	for field := uint8(1); field <= 8; field *= 2 {
		if field&uint8(fields) != 0 {
			switch BL_NACK_t(field) {
			case BL_NACK_INVALID_CMD:
				str := "Invalid command"
				strArray = append(strArray, str)
			case BL_NACK_INVALID_KEY:
				str := "Invalid key"
				strArray = append(strArray, str)
			case BL_NACK_INVALID_ADDRESS:
				str := "Invalid address"
				strArray = append(strArray, str)
			case BL_NACK_INVALID_LENGTH:
				str := "Invalid length"
				strArray = append(strArray, str)
			case BL_NACK_INVALID_DATA:
				str := "Invalid data"
				strArray = append(strArray, str)
			case BL_NACK_INVALID_CRC:
				str := "Invalid CRC"
				strArray = append(strArray, str)
			case BL_NACK_OPERATION_FAILURE:
				str := "Operation failure"
				strArray = append(strArray, str)
			default:
				break
			}
		}
	}

	return strArray
}

// Map of response types
var responseMap = map[uint8]func([]byte) (interface{}, error){
	BL_VER_CMD_ID: func(data []byte) (interface{}, error) {
		var response BL_ReadVersionResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	},
	BL_FLASH_ERASE_CMD_ID: func(data []byte) (interface{}, error) {
		var response BL_FlashEraseResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	},
	BL_MEM_WRITE_CMD_ID: func(data []byte) (interface{}, error) {
		var response BL_FlashWriteResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	},
	BL_MEM_READ_CMD_ID: func(data []byte) (interface{}, error) {
		var response BL_FlashReadResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	},
	BL_JUMP_TO_APP_CMD_ID: func(data []byte) (interface{}, error) {
		var response BL_JumpToAppResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	},
}
