package main

type BL_Command interface {
	Type() uint8
}

func (command BL_FlashEraseRequest) Type() uint8 {
	return BL_FLASH_ERASE_CMD_ID
}

func (command BL_ReadVersionRequest) Type() uint8 {
	return BL_VER_CMD_ID
}

func (command BL_JumpToAppRequest) Type() uint8 {
	return BL_JUMP_TO_APP_CMD_ID
}
func (command BL_FlashReadRequest) Type() uint8 {
	return BL_MEM_READ_CMD_ID
}

func (command BL_FlashWriteRequest) Type() uint8 {
	return BL_MEM_WRITE_CMD_ID
}

var command_channel = make(chan BL_Command, 5)
