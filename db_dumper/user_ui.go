package db_dumper

import (
	"errors"
	"strings"
	"fmt"
	"strconv"
	"bufio"
	"os"
	"github.com/daviddengcn/go-colortext"
	"github.com/Orange-OpenSource/db-dumper-cli-plugin/db_dumper/model"
)

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
func (this *DbDumperManager) selectDump(serviceInstance string, recent bool, dumpDateOrNumber string) (model.Dump, error) {
	dumps, err := this.getDumps(serviceInstance)
	if err != nil {
		return model.Dump{}, err
	}
	if len(dumps) == 0 {
		return model.Dump{}, errors.New("There is no dumps available")
	}
	createdAt := dumpDateOrNumber
	if index, err := strconv.Atoi(dumpDateOrNumber); err == nil {
		if index < len(dumps) && index >= 0 {
			createdAt = dumps[index].CreatedAt
		} else {
			return model.Dump{}, errors.New("Dump number " + dumpDateOrNumber + " is not valid. (use 'db-dumper list')")
		}
	}
	if dumpDateOrNumber != "" && createdAt == "" && !recent {
		createdAt = dumpDateOrNumber
	} else if recent {
		createdAt = dumps[0].CreatedAt
	}
	if !recent && createdAt == "" {
		createdAt, err = this.selectDumpDate(serviceInstance, dumps, "At which date do you want your dump file ?")
		if err != nil {
			return model.Dump{}, err
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
		return model.Dump{}, errors.New("The dump at the date of " + createdAt + " doesn't exist")
	}
	return selectedDump, nil
}
