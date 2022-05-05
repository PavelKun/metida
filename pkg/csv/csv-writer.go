package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"sync"
)

// CsvWriter represent encoding/csv.Writer with mutex.
type CsvWriter struct {
	mutex     *sync.Mutex
	csvWriter *csv.Writer
}

// NewCsvWriter return new instance of CsvWriter.
func NewCsvWriter(fileName string) (*CsvWriter, error) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(csvFile)

	return &CsvWriter{csvWriter: w, mutex: &sync.Mutex{}}, nil
}

// Write writes a single CSV record with mutex.
func (w *CsvWriter) Write(row []string) error {
	w.mutex.Lock()
	err := w.csvWriter.Write(row)
	w.mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}

// Flush writes any buffered data with mutex.
func (w *CsvWriter) Flush() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.csvWriter.Flush()
}

func (w *CsvWriter) worker(wg *sync.WaitGroup, ports chan interface{}, once *sync.Once) {
	defer wg.Done()

	for p := range ports {
		once.Do(func() {
			w.writeNamesStruct(p)
		})
		w.writeValueStruct(p)

		w.Flush() // Делает запись в файл после каждого значения
	}
}

// writeNameStruct writes a names structure fields in file.
func (w *CsvWriter) writeNamesStruct(val interface{}) {
	var sliceWriter []string

	elem := reflect.ValueOf(val).Elem()

	for i := 0; i < elem.NumField(); i++ {
		name := elem.Type().Field(i).Name
		sliceWriter = append(sliceWriter, fmt.Sprint(name))
	}

	err := w.Write(sliceWriter)
	if err != nil {
		//ToDo logging
	}
}

// writeNameStruct writes a value structure fields in file.
func (w *CsvWriter) writeValueStruct(val interface{}) {
	var sliceWriter []string

	elem := reflect.ValueOf(val).Elem()

	for i := 0; i < elem.NumField(); i++ {
		value := elem.Field(i)
		sliceWriter = append(sliceWriter, fmt.Sprint(value))
	}

	err := w.Write(sliceWriter)
	if err != nil {
		//ToDo logging
	}
}

// WriteCsvFromChan launches workers to write in file.
func WriteCsvFromChan(wg *sync.WaitGroup, ports chan interface{}, numberWorkers int, nameFile string) error {
	w, err := NewCsvWriter(nameFile)
	if err != nil {
		return err
	}

	var once sync.Once

	for i := 0; i < numberWorkers; i++ {
		wg.Add(1)

		go w.worker(wg, ports, &once)
	}

	return nil
}
