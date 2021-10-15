package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/metno/agnc/pkg/downloader"
	mmspublisher "github.com/metno/agnc/pkg/mms-publisher"
	"github.com/metno/agnc/pkg/utils"
	"github.com/metno/go-mms/pkg/mms"
	_ "gocloud.dev/blob/s3blob"
)

var timeout time.Duration = 360

func cleanFile(filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return
	}

	err = os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func processTimestep(ctx context.Context, args ...interface{}) error {
	s3URL := args[0].(string)
	url, err := url.Parse(s3URL)
	if err != nil {
		log.Fatal(err)
	}
	timestepID := filepath.Base(url.Path)

	gribDir := args[1].(string)
	ncDir := args[2].(string)
	configDir := args[3].(string)
	fimexConfigPath := fmt.Sprintf("%s/grib2nc.conf", configDir)
	inputConfigPath := fmt.Sprintf("%s/config.grib.xinclude.xml", configDir)
	inputFilePath := fmt.Sprintf("%s/%s", gribDir, timestepID)
	outputConfigPath := fmt.Sprintf("%s/cdm_writer_config.xml", configDir)
	outputFilePath := fmt.Sprintf("%s/%s.nc", ncDir, timestepID)
	outputTempFilePath := fmt.Sprintf("%s.tmp", outputFilePath)

	tooOld, err := utils.IsTooOld(timestepID, time.Now().UTC(), 1)
	if err != nil {
		log.Fatal(err)
	}

	if tooOld {
		log.Printf("%s timestep too old, skipping...", timestepID)
		return nil
	}

	defer cleanFile(inputFilePath)
	defer cleanFile(outputTempFilePath)
	downloader.DownloadTimestep(ctx, url.Host, timestepID, gribDir)

	cmd := exec.Command("fimex", "-n", "8", "-c", fimexConfigPath, "--input.config", inputConfigPath, "--input.file", inputFilePath, "--output.config", outputConfigPath, "--output.file", outputTempFilePath)
	log.Println(cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	timer := time.AfterFunc(timeout*time.Second, func() {
		log.Printf("Timeout reached, killing job %d", cmd.Process.Pid)
		cmd.Process.Kill()
	})

	err = os.Rename(outputTempFilePath, outputFilePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Subprocess %d done, deleting file %s\n", cmd.Process.Pid, inputFilePath)

	mmsProductionHub, ok := os.LookupEnv("MMS_PRODUCTION_HUB")
	if !ok {
		log.Fatal("MMS_PRODUCTION_HUB not found in the environment, can't publish to MMS")
	}

	productEvent := mms.ProductEvent{
		JobName:         "processTimestep",
		Product:         timestepID,
		ProductLocation: outputFilePath,
		ProductionHub:   mmsProductionHub,
		CreatedAt:       time.Now(),
		NextEventAt:     time.Now().Add(time.Hour * time.Duration(6)),
	}

	mmspublisher.PublishProduct(productEvent)

	if !timer.Stop() {
		log.Println("Can't stop the timer, it has been stopped before.")
	}

	return nil
}

func main() {
	mgr := worker.NewManager()
	mgr.Register("ProcessTimestep", processTimestep)
	mgr.Concurrency = 5
	mgr.Run()
}
