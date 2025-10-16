package znakdb

import "fmt"

// сделано отдельно чтобы закрывать бд
func (z *DbZnak) Example() (err error) {
	sess := z.dbSession
	defer func() {
		if err != nil {
			if errClose := sess.Close(); errClose != nil {
				err = fmt.Errorf("%w%w", errClose, err)
			}
		} else {
			err = sess.Close()
		}
	}()

	return sess.Ping()
}
