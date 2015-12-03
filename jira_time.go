package gojira
import "time"

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05.999-0700"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	ct.Time, err = time.Parse(ctLayout, string(b))
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format(ctLayout)), nil
}

var nilTime = (time.Time{}).UnixNano()
func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}