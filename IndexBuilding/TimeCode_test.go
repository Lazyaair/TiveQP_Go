package indexbuilding

import (
	"fmt"
	"testing"
)

func TestTimePointEncoding(t *testing.T) {
	fmt.Println(TimePointEncoding(12, 11))
}

func TestTimeRangeEncoding(t *testing.T) {
	fmt.Println(TimeRangeEncoding(8, 0, 9, 0))
}

func TestTimeRangeEncodingComplement(t *testing.T) {
	fmt.Println(TimeRangeEncodingComplement(8, 0, 12, 0))
}
