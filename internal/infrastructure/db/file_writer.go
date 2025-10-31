package db

import (
	"encoding/json"
	"os"
)

type DataWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewDataWriter(filename string) (*DataWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	e := json.NewEncoder(file)
	return &DataWriter{
		file:    file,
		encoder: e,
	}, nil
}

func (dw *DataWriter) WriteRecords(records *[]Record) error {
	if err := dw.file.Truncate(0); err != nil {
		return err
	}
	if _, err := dw.file.Seek(0, 0); err != nil {
		return err
	}
	err := dw.encoder.Encode(records)
	if err != nil {
		return err
	}
	return nil
}

func (dw *DataWriter) Close() error {
	return dw.file.Close()
}

type DataReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewDataReader(filename string) (*DataReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(file)
	return &DataReader{
		file:    file,
		decoder: d,
	}, nil
}

func (dr *DataReader) Close() error {
	return dr.file.Close()
}

func (dr *DataReader) ReadRecords() (*[]Record, error) {
	var records []Record
	err := dr.decoder.Decode(&records)
	if err != nil {
		return nil, err
	}
	return &records, nil
}
