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

func (p *DataWriter) WriteRecords(records *[]Record) error {
	if err := p.file.Truncate(0); err != nil {
		return err
	}
	if _, err := p.file.Seek(0, 0); err != nil {
		return err
	}
	err := p.encoder.Encode(records)
	if err != nil {
		return err
	}
	return nil
}

func (p *DataWriter) Close() error {
	return p.file.Close()
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

func (p *DataReader) Close() error {
	return p.file.Close()
}

func (c *DataReader) ReadRecords() (*[]Record, error) {
	var records []Record
	err := c.decoder.Decode(&records)
	if err != nil {
		return nil, err
	}
	return &records, nil
}
