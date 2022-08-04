package memcached

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func (m *Memcached) Get(key string) ([]byte, error) {
	const end = "END"

	conn := m.GetConn()
	defer func() {
		m.connsChan <- conn
	}()

	_, err := conn.Write([]byte("get " + key + "\r\n"))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(conn)
	var data []byte
	for scanner.Scan() {
		if scanner.Text() == end {
			return data, scanner.Err()
		}

		data = scanner.Bytes()
	}

	return data, scanner.Err()
}

func (m *Memcached) Set(key string, value []byte, exp int) error {
	const stored = "STORED"

	lenStr := strconv.Itoa(len(value))
	expStr := strconv.Itoa(exp)
	preparedData := []byte("set " + key + " 0 " + expStr + " " + lenStr + " \r\n")
	preparedData = append(preparedData, value...)
	preparedData = append(preparedData, []byte("\r\n")...)

	conn := m.GetConn()
	defer func() {
		m.connsChan <- conn
	}()

	_, err := conn.Write(preparedData)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), stored) {
			return nil
		}

		return fmt.Errorf("memcached: %v", scanner.Text())
	}

	return scanner.Err()
}

func (m *Memcached) Delete(key string) error {
	const deleted = "DELETED"

	conn := m.GetConn()
	defer func() {
		m.connsChan <- conn
	}()

	_, err := conn.Write([]byte("delete " + key + "\r\n"))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), deleted) {
			return nil
		}

		return fmt.Errorf("memcached: %v", scanner.Text())
	}

	return scanner.Err()
}
