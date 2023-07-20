package console_interface

import (
	"acquire/internal/acquire"
	"fmt"
	"strings"
)

type ConsoleInputInterface struct {
}

func (c *ConsoleInputInterface) GetInput(request acquire.InputRequest) (acquire.InputResponse, error) {
	switch request.InputType {

	case "hotel":
		fallthrough
	case "stock":
		fallthrough
	case "number":
		fallthrough
	case "mergerAction":
		fallthrough
	case "tile":
		var val string
		fmt.Printf("%s ", request.Instruction)
		_, err := fmt.Scanf("%s", &val)
		if err != nil {
			if err.Error() == "unexpected newline" {
				return acquire.InputResponse{}, nil
			}
			return acquire.InputResponse{}, err
		}
		val = strings.ToUpper(val)
		return acquire.InputResponse{Value: val}, nil

	case "msg":
		fmt.Println(request.Instruction)
		return acquire.InputResponse{}, nil

	default:
		panic(fmt.Sprintf("can't handle input request type %s", request.InputType))
	}
}
