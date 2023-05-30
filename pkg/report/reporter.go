package report

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"go.uber.org/zap"
)

type Reporter struct {
	URL     string
	outPath string
	client  http.Client
	id      atomic.Int64
}

func NewReporter(url, outPath string) *Reporter {
	var err error
	if outPath == "" {
		outPath, err = os.Getwd()
		if err != nil {
			logging.StatusLog().Fatal("fail to get pwd for outpath", zap.Error(err))
		}
	}
	fileInfo, err := os.Stat(outPath)
	if err != nil {
		logging.StatusLog().Fatal("fail to check outpath", zap.Error(err))
	}
	if !fileInfo.IsDir() {
		logging.StatusLog().Fatal("report output path must be a directory")
	}
	return &Reporter{
		client:  http.Client{Timeout: time.Duration(10) * time.Second},
		URL:     url,
		outPath: outPath,
	}
}

func (r *Reporter) Report(result *ReportResult) {
	if err := r.HttpReport(result); err != nil {
		logging.StatusLog().Error("fail to http report", zap.Error(err))
	}
	if err := r.StdoutReport(result); err != nil {
		logging.StatusLog().Error("fail to stdout report", zap.Error(err))
	}
	if err := r.FileReport(result); err != nil {
		logging.StatusLog().Error("fail to file report", zap.Error(err))
	}
}

func (r *Reporter) HttpReport(result *ReportResult) error {
	if r.URL != "" {
		data, err := json.Marshal(result)
		if err != nil {
			return err
		}
		resp, err := r.client.Post(r.URL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	return nil
}

const defaultIndentation = "    "

func (r *Reporter) FileReport(result *ReportResult) error {
	var output string
	var err error
	id := r.id.Add(1)
	output = filepath.Join(r.outPath, string(id)+".txt")
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	r.WriteBuffer(result, writer)
	return file.Sync()
}

func (r *Reporter) StdoutReport(result *ReportResult) error {
	var buffer bytes.Buffer
	r.WriteBuffer(result, &buffer)
	_, err := fmt.Fprintln(os.Stdout, buffer)
	return err
}

type write interface {
	WriteString(s string) (int, error)
}

func (r *Reporter) WriteBuffer(result *ReportResult, buffer write) {
	// Database Information
	buffer.WriteString(title("Database Information"))
	buffer.WriteString("database version: " + result.DatabaseVersion + "\n")

	// Test Information
	buffer.WriteString(title("Test Information"))
	buffer.WriteString("test time: " + result.Timestamp.String() + "\n")
	buffer.WriteString("test method: " + result.Method + "\n")
	buffer.WriteString("test prepare: " + "\n")
	for _, p := range result.Process {
		buffer.WriteString(defaultIndentation)
		buffer.WriteString(p)
		buffer.WriteString("\n")
	}
	if result.Stack != "" {
		buffer.WriteString("stack: \n")
		buffer.WriteString(result.Stack)
		buffer.WriteString("\n")
	}
	buffer.WriteString("error sql: \n")
	buffer.WriteString(result.ErrorSql)
	buffer.WriteString("\n")

	// Environment Information
	buffer.WriteString(title("Environment Information"))
	for k, v := range result.EnvironmentVariables {
		buffer.WriteString(k + ": " + v + "\n")
	}
}

func title(t string) string {
	return strings.Repeat("#", 10) + t + strings.Repeat("#", 10) + "\n"
}
