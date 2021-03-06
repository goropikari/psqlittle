package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/goropikari/psqlittle/backend"
	"github.com/goropikari/psqlittle/core"
	trans "github.com/goropikari/psqlittle/translator"
)

const (
	payloadBytesLength = 4
	tagLength          = 1
	bufferSize         = 1024
)

var dbmsPORT = getEnvWithDefault("DBMS_PORT", "5432")
var dbmsHOST = getEnvWithDefault("DBMS_HOST", "127.0.0.1")
var dataPath = getEnvWithDefault("DBMS_DATA_PATH", "data.db")
var queryReady []byte = []byte{0x5a, 0x00, 0x00, 0x00, 0x05, 0x49}
var acceptMsg []byte = []byte{0x43, 0x00, 0x00, 0x00, 0x7, 0x4f, 0x4b, 0x00}

// Run starts DBMS server
func Run() {
	db, path := setupDB()
	ln, err := net.Listen("tcp", dbmsHOST+":"+dbmsPORT)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn, db, path)
	}
}

func handleConnection(c net.Conn, db backend.DB, path string) {

	startup(c)
	defer c.Close()
	for {
		tag, query, err := readQuery(c)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				os.Exit(1)
			}
			break
		}
		if tag == 0x58 {
			// 0x58 -> X: terminate
			return
		}
		res, err := handleQuery(db, query)
		if err != nil {
			fmt.Println(err)
			// Ideally, error msg should be sent if errors occur
			c.Write(makeCommandCompleteMsg(err.Error()))
			c.Write(queryReady)
			continue
		}
		if res == nil {
			// Query except for SELECT
			c.Write(acceptMsg)
			c.Write(queryReady)
			writeLog(path, query)
		} else {
			sendResult(c, res)
		}
	}
}

func startup(c net.Conn) error {
	// https://www.pgcon.org/2014/schedule/attachments/330_postgres-for-the-wire.pdf
	// https://www.postgresql.org/docs/12/protocol-message-formats.html
	sizeByte, err := read(c, payloadBytesLength)
	if err != nil {
		return err
	}
	c.Write([]byte{0x4e})

	size := int(binary.BigEndian.Uint32(sizeByte))
	if _, err := read(c, size-payloadBytesLength); err != nil {
		return err
	}
	// AuthenticationOk
	// 0x52 -> Z: ReadyForQuery
	c.Write([]byte{0x52, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00})
	// fake client encoding for python PostgreSQL connector
	c.Write([]byte{0x53, 0x00, 0x00, 0x00, 0x19, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x00, 0x55, 0x54, 0x46, 0x38, 0x00})
	// // fake server version
	// c.Write([]byte{0x53, 0x00, 0x00, 0x00, 0x18, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x00, 0x31, 0x32, 0x2e, 0x36, 0x00})

	// ReadyForQuery
	c.Write(queryReady)

	return nil
}

func makeCommandCompleteMsg(s string) []byte {
	body := make([]byte, 0)
	body = append(body, []byte(s)...)
	body = append(body, 0x00)
	l := len(body)
	lb := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(lb, uint32(l+payloadBytesLength))
	payload := make([]byte, 0)
	payload = append(payload, 0x43) // 0x43 -> C: CommandComplete
	payload = append(payload, lb...)
	payload = append(payload, body...)

	return payload
}

func sendResult(c net.Conn, res trans.Result) {
	cols := res.GetColumns()
	header := makeColDesc(cols)
	c.Write(header)
	recs := res.GetRecords()
	if len(recs) != 0 {
		rowByte := makeDataRows(recs)
		c.Write(rowByte)
	}

	c.Write(selectFooter(len(recs)))
	c.Write(queryReady)
}

func selectFooter(n int) []byte {
	body := []byte("SELECT ")
	s := fmt.Sprintf("%v", n)
	body = append(body, []byte(s)...)
	body = append(body, 0x00)

	payload := make([]byte, 0)
	payload = append(payload, 0x43)
	lenBytes := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(body)+payloadBytesLength))
	payload = append(payload, lenBytes...)
	payload = append(payload, body...)

	return payload
}

func makeDataRow(rec core.Values) []byte {
	dataRow := make([]byte, 0)
	nc := len(rec)
	ncb := make([]byte, 2)
	binary.BigEndian.PutUint16(ncb, uint16(nc))
	dataRow = append(dataRow, ncb...)
	for _, val := range rec {
		if val == nil {
			dataRow = append(dataRow, []byte{0xff, 0xff, 0xff, 0xff}...)
		} else {
			s := fmt.Sprintf("%v", val)
			sb := []byte(s)
			slen := len(sb)
			lenByte := make([]byte, payloadBytesLength)
			binary.BigEndian.PutUint32(lenByte, uint32(slen))
			dataRow = append(dataRow, lenByte[:]...)
			dataRow = append(dataRow, sb[:]...)
		}
	}

	payload := make([]byte, 0)
	payload = append(payload, 0x44) // 0x44 -> D: DataRow
	lenByte := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(lenByte, uint32(len(dataRow)+payloadBytesLength))
	payload = append(payload, lenByte...)
	payload = append(payload, dataRow...)

	return payload
}

func makeColDesc(cols []string) []byte {
	payload := make([]byte, 0)
	n := len(cols)
	numCols := make([]byte, 2)
	binary.BigEndian.PutUint16(numCols, uint16(n))
	payload = append(payload, numCols[:]...)

	for k, col := range cols {
		payload = append(payload, []byte(col)...)
		payload = append(payload, 0x00)
		payload = append(payload, []byte{0x00, 0x00, 0x40, 0x06}...) // object id
		idx := make([]byte, 2)
		binary.BigEndian.PutUint16(idx, uint16(k+1))
		payload = append(payload, idx[:]...)                         // col id
		payload = append(payload, []byte{0x00, 0x00, 0x04, 0x13}...) // data type
		payload = append(payload, []byte{0xff, 0xff}...)             // data type size
		payload = append(payload, []byte{0xff, 0xff, 0xff, 0xff}...) // type modifier
		payload = append(payload, []byte{0x00, 0x00}...)             // format code
	}

	length := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(length, uint32(len(payload)+payloadBytesLength))
	packet := make([]byte, 0)
	packet = append(packet, 0x54) // 0x54 -> T: RowDescription
	packet = append(packet, length[:]...)
	packet = append(packet, payload[:]...)

	return packet
}

func makeDataRows(recs core.ValuesList) []byte {
	dataRows := make([]byte, 0)
	for _, rec := range recs {
		dataRows = append(dataRows, makeDataRow(rec)...)
	}

	return dataRows
}

func handleQuery(db backend.DB, query string) (trans.Result, error) {
	raNode, err := trans.NewPGTranslator(query).Translate()
	if err != nil {
		return nil, err
	}
	res, err := raNode.Eval(db)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readQuery(c net.Conn) (byte, string, error) {
	data := make([]byte, 0)
	buf := make([]byte, bufferSize)
	for {
		n, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				os.Exit(1)
			}
			break
		}
		if n < bufferSize {
			data = append(data, buf[:n]...)
			break
		}
		data = append(data, buf[:]...)
	}
	tag := data[0]
	size := parseSize(data[1:5])
	var query string
	if size >= 5 {

		query = string(data[5:size][:])
	}

	return tag, query, nil
}

func parseSize(bs []byte) int {
	return int(binary.BigEndian.Uint32(bs))
}

func read(c net.Conn, n int) ([]byte, error) {
	reader := bufio.NewReader(c)
	data := make([]byte, 0)
	for i := 0; i < n; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		data = append(data, b)
	}

	return data, nil
}

func writeLog(path, query string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(query); err != nil {
		log.Println(err)
	}
}

func setupDB() (backend.DB, string) {
	path := dataPath

	db := backend.NewDatabase()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return db, path
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ss := strings.Split(string(bytes), ";")
	for _, s := range ss {
		if strings.Trim(s, " \n") == "" {
			continue
		}
		raNode, _ := trans.NewPGTranslator(s).Translate()
		_, err := raNode.Eval(db)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return db, path
}

func getEnvWithDefault(key string, d string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return d
}
