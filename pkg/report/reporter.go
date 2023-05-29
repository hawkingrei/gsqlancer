package report

import (
	"bufio"
	"bytes"
	"encoding/json"
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
		logging.StatusLog().Error("fail to report", zap.Error(err))
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

func (r *Reporter) FileReport(result *Reporter) error {
	var output string
	var err error
	if r.outPath == "" {
		output, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		id := r.id.Add(1)
		output = filepath.Join(r.outPath, string(id)+".txt")
	}
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString(title("Database Information"))
	writer.WriteString(title("Test Information"))
	writer.WriteString(title("Environment Information"))
	return file.Sync()
}

func title(t string) string {
	return strings.Repeat("#", 10) + t + strings.Repeat("#", 10) + "\n"
}
