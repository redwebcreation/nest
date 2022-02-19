package cloud

import "fmt"

func ValidateDsn(dsn string) error {
	if len(dsn) != 45 {
		return fmt.Errorf("invalid dsn: %s", dsn)
	}

	return nil
}

func ParseDsn(dsn string) (id, token string, err error) {
	if err = ValidateDsn(dsn); err != nil {
		return "", "", err
	}

	return dsn[:22], dsn[23:], nil
}

func FormatDsn(id, token string) string {
	return fmt.Sprintf("%s:%s", id, token)
}
