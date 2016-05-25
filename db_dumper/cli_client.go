package db_dumper

import (
	"errors"
	"strings"
	"strconv"
	"github.com/Orange-OpenSource/db-dumper-cli-plugin/db_dumper/model"
	"encoding/json"
	"fmt"
	"github.com/Orange-OpenSource/db-dumper-cli-plugin/db_dumper/progress_bar"
	"time"
)

func (this *DbDumperManager) waitServiceAction(serviceName string, action string) error {

	ipb := progress_bar.NewIndeterminateProgressBar(action)

	for ; true; {
		service, err := this.cliConnection.GetService(serviceName)
		if (err != nil) {
			time.Sleep(2 * time.Second)
			service, err = this.cliConnection.GetService(serviceName)
			if (err != nil) {
				return err;
			}
		}
		state := service.LastOperation.State
		switch (state) {
		case "succeeded":
			fmt.Println(action + " finished.")
			return nil;
		case "in progress":
			break;
		case "failed":
		case "internal error":
			return errors.New(service.LastOperation.Description);
		}
		ipb.Next()
	}
	return nil;
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

func (this *DbDumperManager) cliCommand(command ...string) ([]string, error) {
	var output []string
	var err error
	if this.verbose {
		output, err = this.cliConnection.CliCommand(command...)
	} else {
		output, err = this.cliConnection.CliCommandWithoutTerminalOutput(command...)
	}
	if (err != nil) {
		return output, err
	}
	output = strings.Split(strings.Join(output, "\n"), "\n")
	return output, err
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
func (this *DbDumperManager) isDbDumperInstance(instance string) bool {
	service, err := this.cliConnection.GetService(instance)
	if err != nil {
		return false
	}
	return service.ServiceOffering.Name == this.serviceName
}
func (this *DbDumperManager) checkIsDbDumperInstance(instance string) error {
	if (!this.isDbDumperInstance(instance)) {
		return errors.New("Instance " + instance + " is not an instance of db-dumper-service")
	}
	return nil
}
func (this *DbDumperManager) generateJsonFrom(template string, values ...interface{}) (string, error) {
	token, err := this.cliConnection.AccessToken()
	if err != nil {
		return "", nil
	}
	if strings.HasPrefix(token, "bearer ") {
		token = strings.TrimPrefix(token, "bearer ")
	}
	org, err := this.cliConnection.GetCurrentOrg()
	if err != nil {
		return "", nil
	}
	space, err := this.cliConnection.GetCurrentSpace()
	if err != nil {
		return "", nil
	}
	values = append(values, token, org.Name, space.Name)
	return fmt.Sprintf(template, values...), nil
}