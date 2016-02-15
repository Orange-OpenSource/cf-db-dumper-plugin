package db_dumper
import (
	"github.com/cloudfoundry/cli/plugin"
	"errors"
	"strings"
	"fmt"
	"strconv"
	"bufio"
	"os"
	"github.com/daviddengcn/go-colortext"
	"net/url"
	"github.com/satori/go.uuid"
	"github.com/Orange-OpenSource/db-dumper-cli-plugin/db_dumper/model"
	"github.com/docker/docker/vendor/src/github.com/jfrazelle/go/canonical/json"
	"github.com/olekukonko/tablewriter"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"crypto/tls"
	"io"
	"github.com/cheggaaa/pb"
)
type DbDumperManager struct {
	cliConnection plugin.CliConnection
	serviceName   string
	verbose       bool
}
var command_create_dump_nonexist = "cs %s %s %s -c"
var command_create_dump_exist = "update-service %s -c"
var command_restore_dump = "update-service %s -c"
var json_restore = "{\"action\": \"restore\", \"target_url\": \"%s\", \"created_at\": \"%s\"}"
var json_dump_nonexist = "{\"src_url\":\"%s\"}"
var json_dump_exist = "{\"action\": \"dump\"}"
var command_delete_dumps = "ds %s"
var service_name_suffix = "-dump"
func NewDbDumperManager(serviceName string, cliConnection plugin.CliConnection, verbose bool) *DbDumperManager {
	return &DbDumperManager{
		cliConnection: cliConnection,
		serviceName: serviceName,
		verbose: verbose,
	}
}

func (this *DbDumperManager) CreateDump(service_name_or_url string) error {
	name, err := this.generateName(service_name_or_url)
	var command []string
	if err != nil {
		return err
	}
	if this.isServiceExist(name) {
		command = strings.Split(fmt.Sprintf(command_create_dump_exist, this.serviceName), " ")
		command = append(command, json_dump_exist)
		_, err := this.cliConnection.CliCommand(command...)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println("Service for this database doesn't exist, create it...")
	fmt.Println("")
	fmt.Println("Searching available plans...")
	fmt.Println("")
	plans, err := this.getPlanFromService()
	if err != nil {
		return err
	}
	fmt.Println("")
	plan, err := this.selectByUser("plans", "Which plans do you want ? ", plans, plans[0], plans[0])
	if err != nil {
		return err
	}
	command = strings.Split(fmt.Sprintf(command_create_dump_nonexist, this.serviceName, plan, name), " ")
	command = append(command, fmt.Sprintf(json_dump_nonexist, service_name_or_url))
	_, err = this.cliConnection.CliCommand(command...)
	return err
}

func (this *DbDumperManager) selectByUser(typeToSelect string, msg string, typeList []string, defaultValueName, defaultValue string) (string, error) {
	fmt.Println("Available " + typeToSelect + ":")

	for num, typeFlat := range typeList {
		ct.Foreground(ct.Blue, false)
		fmt.Print(strconv.Itoa(num))
		ct.ResetColor()
		fmt.Println(". " + typeFlat)
	}


	typeSelect := ""
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Println("")
		fmt.Println(msg)
		fmt.Print("Choice <")
		ct.Foreground(ct.Blue, false)
		fmt.Print(defaultValueName)
		ct.ResetColor()
		fmt.Print(">: ")
		typeBytes, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		typeNameOrId := string(typeBytes)
		if typeNameOrId == "" {
			return defaultValue, nil
		}
		typeSelect, err = this.findDatabyNameOrId(typeList, typeNameOrId)
		if err != nil {
			showError(err)
			continue
		}
		break
	}
	return typeSelect, nil
}
func (this *DbDumperManager) selectDumpDate(serviceInstance string, dumps []model.Dump, msg string) (string, error) {
	dates := make([]string, 0)
	for _, dump := range dumps {
		dates = append(dates, dump.CreatedAt)
	}
	return this.selectByUser("dump dates", msg, dates, "recent", "")
}
func (this *DbDumperManager) RestoreDump(target_service_name_or_url string, recent bool) error {
	serviceInstance, err := this.selectService("Which instance do you want to restore to '" + target_service_name_or_url + "' ?")
	if err != nil {
		return err
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

	command := strings.Split(fmt.Sprintf(command_restore_dump, serviceInstance), " ")
	command = append(command, fmt.Sprintf(json_restore, target_service_name_or_url, createdAt))
	_, err = this.cliConnection.CliCommand(command...)
	return err
}
func (this *DbDumperManager) selectService(msg string) (string, error) {

	if !this.isDbDumperServiceExist() {
		return "", errors.New("Cannot found service: " + this.serviceName)
	}

	fmt.Println("Searching available " + this.serviceName + " instances...")
	fmt.Println("")
	serviceInstance := ""
	serviceInstances, err := this.getDbDumperServiceInstance()
	if err != nil {
		return "", err
	}
	if len(serviceInstances) == 0 {
		return "", errors.New("No " + this.serviceName + " instance exist. Please create a dump first.")
	}
	prefix, err := this.GetNamePrefix()
	if err != nil {
		return "", err
	}
	fmt.Println("Available " + this.serviceName + " instances")
	firstServiceInstance := ""
	for num, serviceInstanceFlat := range serviceInstances {
		if strings.HasPrefix(serviceInstanceFlat, prefix) {
			serviceInstanceFlat = strings.TrimPrefix(serviceInstanceFlat, prefix)
		}
		if strings.HasSuffix(serviceInstanceFlat, service_name_suffix) {
			serviceInstanceFlat = strings.TrimSuffix(serviceInstanceFlat, service_name_suffix)
		}
		if num == 0 {
			firstServiceInstance = serviceInstanceFlat
		}
		ct.Foreground(ct.Blue, false)
		fmt.Print(strconv.Itoa(num))
		ct.ResetColor()
		fmt.Println(". " + serviceInstanceFlat)
	}
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Println("")
		fmt.Println(msg)
		fmt.Print("Choice <")
		ct.Foreground(ct.Blue, false)
		fmt.Print(firstServiceInstance)
		ct.ResetColor()
		fmt.Print(">: ")
		planBytes, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		serviceInstanceNameOrId := string(planBytes)
		if serviceInstanceNameOrId == "" {
			serviceInstanceNameOrId = serviceInstances[0]
		}
		serviceInstance, err = this.findDatabyNameOrId(serviceInstances, prefix + serviceInstanceNameOrId + service_name_suffix)
		if err == nil {
			break
		}
		serviceInstance, err = this.findDatabyNameOrId(serviceInstances, serviceInstanceNameOrId)
		if err == nil {
			break
		}
		showError(err)
		continue
	}
	return serviceInstance, nil
}
func (this *DbDumperManager) DownloadDump(skipInsecure bool, recent bool, inStdout bool, dumpDateOrNumber string) error {
	if inStdout {
		return errors.New("To use stdout option you need to pass a service instance.")
	}
	serviceInstance, err := this.selectService("Which instance to list ?")
	if err != nil {
		return err
	}
	return this.DownloadDumpFromInstanceName(serviceInstance, skipInsecure, recent, inStdout, dumpDateOrNumber)
}
func (this *DbDumperManager) DownloadDumpFromInstanceName(serviceInstance string, skipInsecure bool, recent bool, inStdout bool, dumpDateOrNumber string) error {
	dumps, err := this.getDumps(serviceInstance)
	if err != nil {
		return err
	}
	createdAt := dumpDateOrNumber
	if index, err := strconv.Atoi(dumpDateOrNumber); err == nil {
		if index < len(dumps) && index >= 0 {
			createdAt = dumps[index].CreatedAt
		}else {
			return errors.New("Dump number " + dumpDateOrNumber + " is not valid. (use 'db-dumper list')")
		}
	}
	if dumpDateOrNumber != "" && createdAt == "" && !recent {
		createdAt = dumpDateOrNumber
	}else if recent {
		createdAt = dumps[0].CreatedAt
	}
	if inStdout && dumpDateOrNumber == "" && !recent {
		return errors.New("stdout option can only be use with flag --dump-number or --recent")
	}
	if !recent && createdAt == "" {
		createdAt, err = this.selectDumpDate(serviceInstance, dumps, "At which date do you want your dump file ?")
		if err != nil {
			return err
		}
		if createdAt == "" {
			createdAt = dumps[0].CreatedAt
		}
	}
	var selectedDump model.Dump
	for _, dump := range dumps {
		if dump.CreatedAt == createdAt {
			selectedDump = dump
			break
		}
	}
	if selectedDump == (model.Dump{}) {
		return errors.New("The dump at the date of " + createdAt + " doesn't exist")
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

	resp, err := client.Get(selectedDump.DownloadURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Dump can't be downloaded, http status code: " + strconv.Itoa(resp.StatusCode))
	}
	fmt.Println("")
	if inStdout {
		io.Copy(os.Stdout, resp.Body)
	}else {
		bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES)
		bar.Start()
		out, err := os.Create(selectedDump.Filename)
		if err != nil {
			return err
		}
		reader := bar.NewProxyReader(resp.Body)
		io.Copy(out, reader)
		bar.Update()
	}

	fmt.Println("")
	if !inStdout {
		fmt.Println("")
		fmt.Print("File as been downloaded in ")
		ct.Foreground(ct.Blue, false)
		fmt.Print(selectedDump.Filename)
		ct.ResetColor()
		fmt.Println(" file")
	}
	return nil
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
	return this.ListFromInstanceNameWithDumps(serviceInstance, showUrl, dumps)
}
func (this *DbDumperManager) ListFromInstanceNameWithDumps(serviceInstance string, showUrl bool, dumps []model.Dump) error {
	fmt.Println("")
	headers := []string{"#", "File Name", "Created At"}

	if showUrl {
		headers = append(headers, "Download Url")
		headers = append(headers, "Dashboard Url")

	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for index, dump := range dumps {
		var data []string
		if showUrl {
			data = []string{strconv.Itoa(index), dump.Filename, dump.CreatedAt, dump.DownloadURL, dump.ShowURL}
		}else {
			data = []string{strconv.Itoa(index), dump.Filename, dump.CreatedAt}
		}
		table.Append(data)
	}
	table.Render()
	return nil
}
func (this *DbDumperManager) cliCommand(command ...string) ([]string, error) {
	if this.verbose {
		return this.cliConnection.CliCommand(command...)
	}
	return this.cliConnection.CliCommandWithoutTerminalOutput(command...)
}
func (this *DbDumperManager) getDumps(serviceInstance string) ([]model.Dump, error) {
	var err error
	command := []string{"create-service-key", serviceInstance, "plugin-key-" + serviceInstance}
	_, err = this.cliCommand(command...)
	if err != nil {
		return nil, err
	}
	command = []string{"service-key", serviceInstance, "plugin-key-" + serviceInstance}
	output, err := this.cliCommand(command...)
	if err != nil {
		return nil, err
	}
	if len(output) < 2 {
		return nil, err
	}
	var credentials model.Credentials
	datasUnparsed := output[2:]
	jsonData := ""
	for _, dataUnparsed := range datasUnparsed {
		jsonData += dataUnparsed
	}
	byt := []byte(jsonData)
	err = json.Unmarshal(byt, &credentials)
	if err != nil {
		return nil, err
	}
	command = []string{"delete-service-key", serviceInstance, "plugin-key-" + serviceInstance, "-f"}
	_, err = this.cliCommand(command...)
	if err != nil {
		return nil, err
	}
	return credentials.Dumps, nil
}
func (this *DbDumperManager) DeleteDump() error {
	serviceInstance, err := this.selectService("Which instance do you want to delete ? (dump will be really delete after a determined period) ?")
	if err != nil {
		return err
	}
	command := strings.Split(fmt.Sprintf(command_delete_dumps, serviceInstance), " ")
	_, err = this.cliConnection.CliCommand(command...)
	return err
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
	_, err := this.cliConnection.GetService(name)
	return err != nil
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
func (this *DbDumperManager) isUri(name string) bool {

	_, err := url.ParseRequestURI(name)
	return err == nil
}
func (this *DbDumperManager) isDbDumperServiceExist() bool {
	command := []string{"m"}
	output, err := this.cliCommand(command...)
	if err != nil {
		return false
	}
	if len(output) < 4 {
		return false
	}
	datasUnparsed := output[4:]
	var service string = ""
	for _, dataUnparsed := range datasUnparsed {
		dataUnparsedSplitted := strings.Split(dataUnparsed, " ")
		if dataUnparsedSplitted[0] != this.serviceName {
			continue
		}
		service = dataUnparsed
	}
	return service != ""
}
func (this *DbDumperManager) getPlanFromService() ([]string, error) {
	command := []string{"m"}
	output, err := this.cliCommand(command...)
	if err != nil {
		return nil, errors.New(strings.Join(output, "\n"))
	}
	if len(output) < 4 {
		return nil, errors.New(strings.Join(output, "\n"))
	}
	datasUnparsed := output[4:]
	var service string = ""
	for _, dataUnparsed := range datasUnparsed {
		dataUnparsedSplitted := strings.Split(dataUnparsed, " ")
		if dataUnparsedSplitted[0] != this.serviceName {
			continue
		}
		service = dataUnparsed
	}
	if service == "" {
		return nil, errors.New("Cannot found service: " + this.serviceName)
	}
	planString := strings.TrimPrefix(service, this.serviceName)
	planString = strings.TrimSpace(planString)
	plans := strings.Split(planString, ", ")


	lastPlanString := plans[len(plans) - 1]
	lastPlan := strings.Split(lastPlanString, " ")
	plans[len(plans) - 1] = lastPlan[0]


	return plans, nil
}
func showError(err error) {
	if err != nil {
		ct.Foreground(ct.Yellow, true)
		fmt.Println(fmt.Sprintf("%v", err))
		ct.ResetColor()
	}
}
func (s *DbDumperManager) findDatabyNameOrId(datas []string, dataNameOrInt string) (string, error) {
	index, err := strconv.Atoi(dataNameOrInt)
	if err == nil {
		if index < 0 || index > len(datas) - 1 {
			return dataNameOrInt, errors.New("Not valid.")
		}
		return datas[index], nil
	}
	for _, data := range datas {
		if data == dataNameOrInt {
			return data, nil
		}
	}
	return dataNameOrInt, errors.New("Not found.")
}