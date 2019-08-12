package semver

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements the database/sql.Scanner interface.
func (v *Version) Scan(src interface{}) (err error) {
	var strVal string
	switch src.(type) {
	case string:
		strVal = src.(string)
	case []byte:
		strVal = string(src.([]byte))
	default:
		return fmt.Errorf("Version.Scan: cannot convert %T to string.", src)
	}

	tmpv, err := Parse(strVal)
	if err != nil {
		return
	}
	*v = *tmpv
	return
}

// Value implements the database/sql/driver.Valuer interface.
func (s Version) Value() (driver.Value, error) {
	return s.String(), nil
}
