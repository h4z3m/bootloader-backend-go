package main

type BL_Response interface {
	Type() uint8
}

var response_channel = make(chan BL_Response, 1)

func (response BL_ReadVersionResponse) Type() uint8 {
	return BL_VER_CMD_ID
}

func (response BL_FlashEraseResponse) Type() uint8 {
	return BL_FLASH_ERASE_CMD_ID

}
func (response BL_FlashWriteResponse) Type() uint8 {
	return BL_MEM_WRITE_CMD_ID
}
func (response BL_JumpToAppResponse) Type() uint8 {
	return BL_JUMP_TO_APP_CMD_ID
}

func (response BL_FlashReadResponse) Type() uint8 {
	return BL_MEM_READ_CMD_ID
}
