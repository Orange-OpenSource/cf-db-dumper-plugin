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
)
type DbDumperManager struct {
	cliConnection plugin.CliConnection
	serviceName   string
}
var command_create_dump_nonexist = "cs db-dumper-service %s %s -c"
var command_create_dump_exist = "update-service %s -c"
var command_restore_dump = "update-service %s -c"
var json_restore = "{\"action\": \"restore\", \"target_url\": \"%s\"}"
var json_dump_nonexist = "{\"src_url\":\"%s\"}"
var json_dump_exist = "{\"action\": \"dump\"}"
var command_delete_dumps = "ds %s"
var service_name_suffix = "-dump"
func NewDbDumperManager(serviceName string, cliConnection plugin.CliConnection) *DbDumperManager {
	return &DbDumperManager{
		cliConnection: cliConnection,
		serviceName: serviceName,
	}
}

func (this *DbDumperManager) CreateDump(service_name_or_url string) error {
	name, err := this.generateName(service_name_or_url)
	var command []string
	if err != nil {
		return err
	}
	if this.isServiceExist(name) {
		command = strings.Split(fmt.Sprintf(command_create_dump_exist), " ")
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
	plans, err := this.getPlanFromService()
	if err != nil {
		return err
	}
	fmt.Println("Available plans:")
	for num, planFlat := range plans {
		fmt.Println(strconv.Itoa(num) + ". " + planFlat)
	}
	plan := ""
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("Which plan do you want <" + plans[0] + "> ? ")
		planBytes, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		planNameOrId := string(planBytes)
		if planNameOrId == "" {
			planNameOrId = plans[0]
		}
		plan, err = this.findDatabyNameOrId(plans, planNameOrId)
		if err != nil {
			showError(err)
			continue
		}
		break
	}
	command = strings.Split(fmt.Sprintf(command_create_dump_nonexist, plan, name), " ")
	command = append(command, fmt.Sprintf(json_dump_nonexist, service_name_or_url))
	_, err = this.cliConnection.CliCommand(command...)
	return err
}
func (this *DbDumperManager) RestoreDump(target_service_name_or_url string) error {
	serviceInstance, err := this.selectService("Which instance do you want to restore to '" + target_service_name_or_url + "'")
	if err != nil {
		return err
	}
	command := strings.Split(fmt.Sprintf(command_restore_dump, serviceInstance), " ")
	command = append(command, fmt.Sprintf(json_restore, target_service_name_or_url))
	_, err = this.cliConnection.CliCommand(command...)
	return err
}
func (this *DbDumperManager) selectService(msg string) (string, error) {

	if !this.isDbDumperServiceExist() {
		return "", errors.New("Cannot found service: " + this.serviceName)
	}

	fmt.Println("Searching available " + this.serviceName + " instances...")
	serviceInstance := ""
	serviceInstances, err := this.getDbDumperServiceInstance()
	if err != nil {
		return "", err
	}
	if len(serviceInstances) == 0 {
		return "", errors.New("No " + this.serviceName + " instance exist. Please create a dump first.")
	}
	fmt.Println("Available " + this.serviceName + " instances:")
	prefix, err := this.getNamePrefix()
	if err != nil {
		return "", err
	}
	for num, serviceInstanceFlat := range serviceInstances {
		if strings.HasPrefix(serviceInstanceFlat, prefix) {
			serviceInstanceFlat = strings.TrimPrefix(serviceInstanceFlat, prefix)
		}
		if strings.HasSuffix(serviceInstanceFlat, service_name_suffix) {
			serviceInstanceFlat = strings.TrimSuffix(serviceInstanceFlat, service_name_suffix)
		}
		fmt.Println(strconv.Itoa(num) + ". " + serviceInstanceFlat)
	}
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print(msg + " <" + serviceInstances[0] + "> ? ")
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
func (this *DbDumperManager) DeleteDump() error {
	serviceInstance, err := this.selectService("Which instance do you want to delete (dump will be really delete after a determined period)")
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
	prefix, err := this.getNamePrefix()
	if err != nil {
		return "", err
	}
	return prefix + name + service_name_suffix, nil
}
func (this *DbDumperManager) getNamePrefix() (string, error) {
	org, err := this.cliConnection.GetCurrentOrg()
	if err != nil {
		return "", err
	}
	space, err := this.cliConnection.GetCurrentSpace()
	if err != nil {
		return "", err
	}
	return org.Name + "-" + space.Name + "-", nil
}
func (this *DbDumperManager) isUri(name string) bool {

	_, err := url.ParseRequestURI(name)
	return err == nil
}
func (this *DbDumperManager) isDbDumperServiceExist() bool {
	command := []string{"m"}
	output, err := this.cliConnection.CliCommandWithoutTerminalOutput(command...)
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
	output, err := this.cliConnection.CliCommandWithoutTerminalOutput(command...)
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