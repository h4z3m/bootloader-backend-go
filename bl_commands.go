package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// readVersion is a function that reads the version information from the request body and sends it through a channel to be processed.
//
// Parameters:
//   - c: a pointer to a gin.Context object representing the HTTP request and response.
//
// Returns:
//
//	None.
func readVersion(c *gin.Context) {
	log.Println("readVersion")

	command := BL_ReadVersionRequest{
		CommandID: BL_VER_CMD_ID,
	}

	c.Bind(&command)
	log.Println(command)

	select {
	case command_channel <- command:
	default:
		log.Println("Channel full")
	}

	log.Println("Sent through channel")
	var response = <-response_channel

	switch response.Type() {
	case BL_VER_CMD_ID:
		response := response.(BL_ReadVersionResponse)
		error_code := response.ErrorCode
		response.Errors = GetErrorStringData(BL_NACK_t(error_code))
		c.IndentedJSON(http.StatusOK, response)
	default:
		c.IndentedJSON(http.StatusInternalServerError, response)
	}
}

// flashErase erases the flash memory.
//
// It takes a *gin.Context as a parameter and does the following:
//  1. Reads the request body and parses it as a BL_FlashEraseRequest struct.
//  2. Sends the command through the command_channel.
//  3. Waits for a response from the response_channel.
//  4. Based on the response type, either returns the response with a HTTP status code of 200 OK
//     or returns a response with a HTTP status code of 500 Internal Server Error.
func flashErase(c *gin.Context) {

	log.Println("flashErase")

	// Read body
	command := BL_FlashEraseRequest{
		CommandID: BL_FLASH_ERASE_CMD_ID,
	}
	c.Bind(&command)
	log.Println(command)

	select {
	case command_channel <- command:
	default:
		log.Println("Channel full")
	}

	log.Println("Sent through channel")

	var response = <-response_channel

	switch response.Type() {
	case BL_FLASH_ERASE_CMD_ID:
		response := response.(BL_FlashEraseResponse)
		error_code := response.ErrorCode
		response.Errors = GetErrorStringData(BL_NACK_t(error_code))
		c.IndentedJSON(http.StatusOK, response)
	default:
		c.IndentedJSON(http.StatusInternalServerError, response)
	}
}

// flashWrite writes the contents of a file to flash memory.
//
// Parameters:
// - c: a pointer to the gin.Context object representing the HTTP request and response.
//
// Returns:
// None.
func flashWrite(c *gin.Context) {
	log.Println("flashWrite")
	command := BL_FlashWriteRequest{
		CommandID: BL_MEM_WRITE_CMD_ID,
	}
	file, err := c.FormFile("file")

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	}
	fileData, err := file.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error)
	}
	// Convert file to binary
	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error)
	}

	command.Data = fileBytes
	size := uint32(file.Size)
	address := c.Request.FormValue("address")
	i, err := strconv.Atoi(address)
	if err != nil {
		log.Println(err)
	}
	command.StartAddress = uint32(i)
	command.Size = size
	log.Println(command)

	select {
	case command_channel <- command:
	default:
		log.Println("Channel full")
	}
	log.Println("Sent through channel")

	var response = <-response_channel

	switch response.Type() {
	case BL_MEM_WRITE_CMD_ID:
		response := response.(BL_FlashWriteResponse)
		error_code := response.ErrorCode
		response.Errors = GetErrorStringData(BL_NACK_t(error_code))
		c.IndentedJSON(http.StatusOK, response)
	default:
		c.IndentedJSON(http.StatusInternalServerError, response)
	}
}

// flashRead receives a request to read data from the flash memory.
//
// It expects a *gin.Context object as a parameter. This context object is used to handle the request and response.
// The function reads the request body and creates a BL_FlashReadRequest command from it.
// The command is then sent through the command_channel.
// If the channel is full, a "Channel full" log message is printed.
// After sending the command, the function waits for a response on the response_channel.
// The response is then processed based on its type.
// If the type is BL_MEM_READ_CMD_ID, the response is converted to a BL_FlashReadResponse.
// The error code is extracted from the response and used to get the corresponding error string data.
// Finally, the response is returned as an indented JSON with a status of http.StatusOK.
// If the response type is not BL_MEM_READ_CMD_ID, the response is returned as an indented JSON with a status of http.StatusInternalServerError.
func flashRead(c *gin.Context) {
	log.Println("flashErase")

	// Read body
	command := BL_FlashReadRequest{
		CommandID: BL_MEM_READ_CMD_ID,
	}

	c.Bind(&command)
	log.Println(command)

	select {
	case command_channel <- command:
	default:
		log.Println("Channel full")
	}

	log.Println("Sent through channel")

	var response = <-response_channel

	switch response.Type() {
	case BL_MEM_READ_CMD_ID:
		response := response.(BL_FlashReadResponse)
		error_code := response.ErrorCode
		response.Errors = GetErrorStringData(BL_NACK_t(error_code))
		c.IndentedJSON(http.StatusOK, response)
	default:
		c.IndentedJSON(http.StatusInternalServerError, response)
	}
}

// jumpToApp is a Go function that handles the request to jump to the app.
//
// The function takes in a *gin.Context parameter.
// It does not return anything.
func jumpToApp(c *gin.Context) {
	log.Println("jumpToApp")

	// Read body
	command := BL_JumpToAppRequest{
		CommandID: BL_JUMP_TO_APP_CMD_ID,
	}
	log.Println(command)

	select {
	case command_channel <- command:
	default:
		log.Println("Channel full")
	}

	log.Println("Sent through channel")

	var response = <-response_channel

	switch response.Type() {
	case BL_JUMP_TO_APP_CMD_ID:
		response := response.(BL_JumpToAppResponse)
		error_code := response.ErrorCode
		response.Errors = GetErrorStringData(BL_NACK_t(error_code))
		c.IndentedJSON(http.StatusOK, response)
	default:
		c.IndentedJSON(http.StatusInternalServerError, response)
	}

}

// getClientEventStatus returns the status of the client event.
func getClientEventStatus() string {
	if connected {
		return "Connected"
	}
	return "Disconnected"

}

// getStatus returns the client event status as an indented JSON response.
//
// Parameter:
// - c: a pointer to the gin.Context object.
//
// Returns:
// - None.
func getStatus(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, getClientEventStatus())
}
