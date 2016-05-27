package db_dumper

import (
	"github.com/cloudfoundry/cli/plugin"
	"errors"
	"strings"
	"fmt"
	"strconv"
	"os"
	"github.com/daviddengcn/go-colortext"
	"net/url"
	"github.com/satori/go.uuid"
	"github.com/Orange-OpenSource/db-dumper-cli-plugin/db_dumper/model"
	"github.com/olekukonko/tablewriter"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"crypto/tls"
	"io"
	"github.com/cheggaaa/pb"
	"github.com/skratchdot/open-golang/open"
	"github.com/dustin/go-humanize"
	"github.com/cloudfoundry/cli/plugin/models"
	"mime"
)

type DbDumperManager struct {
	cliConnection plugin.CliConnection
	serviceName   string
	verbose       bool
}

const (
	command_create_dump_nonexist = "cs %s %s %s -c"
	command_create_dump_exist = "update-service %s -c"
	command_restore_dump = "update-service %s -c"
	json_restore = "{\"action\": \"restore\", \"target_url\": \"%s\", \"created_at\": \"%s\", \"cf_user_token\": \"%s\", \"org\": \"%s\", \"space\": \"%s\"}"
	json_dump_nonexist = "{\"src_url\":\"%s\", \"cf_user_token\": \"%s\", \"org\": \"%s\", \"space\": \"%s\"}"
	json_dump_exist = "{\"action\": \"dump\", \"cf_user_token\": \"%s\", \"org\": \"%s\", \"space\": \"%s\"}"
	command_delete_dumps = "ds %s"
	service_name_suffix = "-dump"
)

func NewDbDumperManager(serviceName string, cliConnection plugin.CliConnection, verbose bool) *DbDumperManager {
	return &DbDumperManager{
		cliConnection: cliConnection,
		serviceName: serviceName,
		verbose: verbose,
	}
}

func (this *DbDumperManager) CreateDump(service_name_or_url string, plan string) error {
	name, err := this.generateName(service_name_or_url)
	var command []string
	if err != nil {
		return err
	}
	if this.isServiceExist(name) {
		command = strings.Split(fmt.Sprintf(command_create_dump_exist, name), " ")
		commandJson, err := this.generateJsonFrom(json_dump_exist)
		if err != nil {
			return err
		}
		command = append(command, commandJson)
		_, err = this.cliCommand(command...)
		if err != nil {
			return err
		}
		err = this.waitServiceAction(name, "Creating dump")
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println("Service for this database doesn't exist, creating it...")
	fmt.Println("")
	if plan == "" {
		plan, err = this.selectPlan()
		if err != nil {
			return err
		}
	}
	command = strings.Split(fmt.Sprintf(command_create_dump_nonexist, this.serviceName, plan, name), " ")
	commandJson, err := this.generateJsonFrom(json_dump_nonexist, service_name_or_url)
	if err != nil {
		return err
	}
	command = append(command, commandJson)
	_, err = this.cliCommand(command...)
	if err != nil {
		return err
	}
	return this.waitServiceAction(name, "Creating dump")
}
func (this *DbDumperManager) RestoreDump(target_service_name_or_url string, recent bool, sourceInstance string) error {
	var serviceInstance string
	var err error
	if sourceInstance != "" {
		serviceInstance = sourceInstance
	} else {
		serviceInstance, err = this.selectService("Which instance do you want to restore to '" + target_service_name_or_url + "' ?")
		if err != nil {
			return err
		}
	}
	createdAt := ""
	if recent == false {
		dumps, err := this.getDumps(serviceInstance)
		if err != nil {
			return err
		}
		createdAt, err = this.selectDumpDate(serviceInstance, dumps, "At which date do you want to restore ?")
		if err != nil {
			return err
		}
	}
	this.cliConnection.AccessToken()
	command := strings.Split(fmt.Sprintf(command_restore_dump, serviceInstance), " ")
	commandJson, err := this.generateJsonFrom(json_restore, target_service_name_or_url, createdAt)
	if err != nil {
		return err
	}
	command = append(command, commandJson)
	_, err = this.cliCommand(command...)
	if err != nil {
		return err
	}
	return this.waitServiceAction(serviceInstance, "Restoring dump")
}
func (this *DbDumperManager) DownloadDump(skipInsecure bool, recent bool, inStdout bool, original bool, dumpDateOrNumber string) error {
	if inStdout {
		return errors.New("To use stdout option you need to pass a service instance.")
	}
	serviceInstance, err := this.selectService("Which instance to list ?")
	if err != nil {
		return err
	}
	return this.DownloadDumpFromInstanceName(serviceInstance, skipInsecure, recent, inStdout, original, dumpDateOrNumber)
}
func (this *DbDumperManager) DownloadDumpFromInstanceName(serviceInstance string, skipInsecure bool, recent bool, inStdout bool, original bool, dumpDateOrNumber string) error {
	if inStdout && dumpDateOrNumber == "" && !recent {
		return errors.New("stdout option can only be use with flag --dump-number or --recent")
	}
	selectedDump, err := this.selectDump(serviceInstance, recent, dumpDateOrNumber)
	if err != nil {
		return err
	}
	var tlsConfig *tls.Config
	if skipInsecure {
		tlsConfig = &tls.Config{InsecureSkipVerify: skipInsecure}
	}
	var transport http.RoundTripper

	transport = &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy: http.ProxyFromEnvironment,
	}
	client := &http.Client{
		Transport: transport,
	}

	downloadUrl := selectedDump.DownloadURL
	if (original) {
		downloadUrl += "?original=1"
	}
	resp, err := client.Get(downloadUrl)
	if err != nil {
		return err
	}

	fileName, err := getFileNameFromHttpResponse(resp)
	if err != nil {
		return err
	}
	if fileName == "" {
		fileName = selectedDump.Filename
	}
	if resp.StatusCode != 200 {
		return errors.New("Dump can't be downloaded, http status code: " + strconv.Itoa(resp.StatusCode))
	}
	fmt.Println("")

	err = downloadFile(resp, fileName, inStdout)
	if err != nil {
		return err
	}
	fmt.Println("")
	if !inStdout {
		fmt.Println("")
		fmt.Print("File as been downloaded in ")
		ct.Foreground(ct.Blue, false)
		fmt.Print(fileName)
		ct.ResetColor()
		fmt.Println(" file")
	}
	return nil
}
func getFileNameFromHttpResponse(resp *http.Response) (string, error) {
	_, dispositionParams, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return "", err
	}
	fileName := dispositionParams["filename"]
	fileName = strings.Replace(fileName, "/", "_", -1)
	return dispositionParams["filename"], nil
}
func downloadFile(resp *http.Response, fileName string, inStdout bool) error {
	if inStdout {
		io.Copy(os.Stdout, resp.Body)
		return nil
	}
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	if resp.ContentLength != -1 {
		bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES)
		bar.Start()
		reader := bar.NewProxyReader(resp.Body)
		io.Copy(out, reader)
		bar.Update()
		return nil
	}
	io.Copy(out, resp.Body)
	return nil
}
func (this *DbDumperManager) ShowDump(recent bool, dumpDateOrNumber string) error {
	serviceInstance, err := this.selectService("Which instance to list ?")
	if err != nil {
		return err
	}
	return this.ShowDumpFromInstanceName(serviceInstance, recent, dumpDateOrNumber)
}
func (this *DbDumperManager) ShowDumpFromInstanceName(serviceInstance string, recent bool, dumpDateOrNumber string) error {
	selectedDump, err := this.selectDump(serviceInstance, recent, dumpDateOrNumber)
	if err != nil {
		return err
	}
	if selectedDump.ShowURL == "" {
		return errors.New("This dump cannot be showed, generally this mean that the file is only in binary.")
	}
	return open.Run(selectedDump.ShowURL)
}
func (this *DbDumperManager) List(showUrl bool) error {
	serviceInstance, err := this.selectService("Which instance to list ?")
	if err != nil {
		return err
	}
	return this.ListFromInstanceName(serviceInstance, showUrl)
}
func (this *DbDumperManager) ListFromInstanceName(serviceInstance string, showUrl bool) error {

	dumps, err := this.getDumps(serviceInstance)
	if err != nil {
		return err
	}
	if len(dumps) == 0 {
		return errors.New("There is no dumps available")
	}
	return this.ListFromInstanceNameWithDumps(serviceInstance, showUrl, dumps)
}
func (this *DbDumperManager) ListFromInstanceNameWithDumps(serviceInstance string, showUrl bool, dumps []model.Dump) error {
	fmt.Println("")
	headers := []string{"#", "File Name", "Created At", "Size", "Is Deleted ?"}

	if showUrl {
		headers = append(headers, "Download Url")
		headers = append(headers, "Dashboard Url")

	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for index, dump := range dumps {
		var data []string
		if showUrl {
			data = []string{strconv.Itoa(index), dump.Filename, dump.CreatedAt, humanize.Bytes(dump.Size), strconv.FormatBool(dump.Deleted), dump.DownloadURL, dump.ShowURL}
		} else {
			data = []string{strconv.Itoa(index), dump.Filename, dump.CreatedAt, humanize.Bytes(dump.Size), strconv.FormatBool(dump.Deleted)}
		}
		table.Append(data)
	}
	table.Render()
	return nil
}
func (this *DbDumperManager) DeleteDump(serviceInstance string, force bool) error {
	var err error
	if serviceInstance == "" {
		serviceInstance, err = this.selectService("Which instance do you want to delete ? (dump will be really delete after a determined period)")
		if err != nil {
			return err
		}
	}
	this.deleteServiceKey(serviceInstance)
	command := strings.Split(fmt.Sprintf(command_delete_dumps, serviceInstance), " ")
	if (force) {
		command = append(command, "-f")
	}
	_, err = this.cliConnection.CliCommand(command...)
	return err
}
func (this *DbDumperManager) GetNamePrefix() (string, error) {
	org, err := this.cliConnection.GetCurrentOrg()
	if err != nil {
		return "", err
	}
	space, err := this.cliConnection.GetCurrentSpace()
	if err != nil {
		return "", err
	}
	hash := md5.Sum([]byte(org.Name + "-" + space.Name + "-"))
	md5String := hex.EncodeToString(hash[:])
	return md5String[:8], nil
}
func (this *DbDumperManager) GetNameSuffix() (string, error) {
	return service_name_suffix, nil
}

func (this *DbDumperManager) getDbDumperServiceInstance() ([]string, error) {
	services, err := this.cliConnection.GetServices()
	if err != nil {
		return nil, err
	}
	dbDumperServices := make([]string, 0)
	for _, service := range services {
		if service.Service.Name == this.serviceName {
			dbDumperServices = append(dbDumperServices, service.Name)
		}

	}
	return dbDumperServices, nil
}
func (this *DbDumperManager) isServiceExist(name string) bool {
	service, err := this.cliConnection.GetService(name)
	if (err != nil) {
		return false
	}
	return service != (plugin_models.GetService_Model{})
}
func (this *DbDumperManager) generateName(name string) (string, error) {

	if this.isUri(name) {
		nameUUID := uuid.NewV5(uuid.NamespaceDNS, name)
		name = nameUUID.String()

	}
	prefix, err := this.GetNamePrefix()
	if err != nil {
		return "", err
	}
	return prefix + name + service_name_suffix, nil
}

func (this *DbDumperManager) isUri(name string) bool {

	_, err := url.ParseRequestURI(name)
	return err == nil
}

func showError(err error) {
	if err != nil {
		ct.Foreground(ct.Yellow, true)
		fmt.Println(fmt.Sprintf("%v", err))
		ct.ResetColor()
	}
}
