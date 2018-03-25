# go-easytransfer
-----------------------------------------------------------------
golang port for arduino easytransfer library https://github.com/madsci1016/Arduino-EasyTransfer


example
```go
import (
	"fmt"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/voyager3m/go-easytransfer"
)

/*
structures should be the same as in arduino
all fields should be public (start with capital letter)
*/
type RecTx struct {
	Cmd      uint8
	Options  uint16
}

type RecRx struct {
	Result  uint8
}

func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/ttyACM0",
		BaudRate:        19200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
		InterCharacterTimeout: 100,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()

	var tx = RecTx{1, 0xffff}
	var rx BulletRx

	/* send data */
	easytransfer.SendData(&tx, port)

	/* receive data */
	time.Sleep(100 * time.Millisecond)
	if easytransfer.ReceiveData(&rx, port) {
	  fmt.Println(rx)
	}
}
```
