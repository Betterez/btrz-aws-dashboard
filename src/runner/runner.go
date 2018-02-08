package main

import (
	"btrzaws"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"io/ioutil"
	"os"
	"strings"
)

func loadStaticReports() string {
	results := ""
	staticFiles, err := ioutil.ReadDir("resources")
	if err != nil {
		return results
	}
	for _, file := range staticFiles {
		// only load json files!
		if strings.Index(file.Name(), ".json") == (len(file.Name()) - 5) {
			data, err := ioutil.ReadFile("./resources" + string(os.PathSeparator) + file.Name())
			if err != nil {
				continue
			}
			fmt.Println("Static file", file.Name(), "loaded")
			fileData := string(data)
			results += fileData + ","
		}
	}
	results = strings.TrimSuffix(results, ",")
	return results
}

func removeDashboard() {
	awsSession, err := btrzaws.GetAWSSession()
	if err != nil {
		fmt.Print(err, "can't get a session")
		os.Exit(1)
	}
	cloudwatchService := cloudwatch.New(awsSession)
	output, err := cloudwatchService.DeleteDashboards(&cloudwatch.DeleteDashboardsInput{
		DashboardNames: []*string{aws.String("btrz-api-cpu"), aws.String("btrz-dash")},
	})
	if err != nil {
		fmt.Println(err, " error reported deleting ")
	} else {
		fmt.Printf("deleted with %v\n", output)
	}

}

func getActiveServices(awsSession *session.Session) ([]*btrzaws.BetterezInstance, error) {
	tags := make([]*btrzaws.AwsTag, 2)
	tags[0] = btrzaws.NewWithValues("tag:Environment", "production")
	tags[1] = btrzaws.NewWithValues("instance-state-name", "running")
	instances, err := btrzaws.GetInstancesWithTags(awsSession, tags)
	if err != nil {
		return nil, err
	}
	btrzInstances := make([]*btrzaws.BetterezInstance, 0)
	for _, reservation := range instances {
		for _, instance := range reservation.Instances {
			currentInstance := btrzaws.LoadFromAWSInstance(instance)
			if currentInstance.ServiceType == "mongo" || currentInstance.ServiceType == "http" || currentInstance.ServiceType == "worker" {
				btrzInstances = append(btrzInstances, currentInstance)
			}
		}
	}
	return btrzInstances, nil
}

func generateDashboard(awsSession *session.Session) {
	activeInstances, err := getActiveServices(awsSession)
	if err == nil {
		fmt.Println("got the list, ", len(activeInstances), " servers")
	} else {
		return
	}
	widgetsData := generateJSONData(activeInstances, awsSession)
	// file, err := os.OpenFile("dumps/temp.json", os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0777)
	// if err != nil {
	// 	fmt.Println("err", err, "creating file")
	// }
	// file.WriteString(widgetsData)
	// defer file.Close()
	cloudwatchService := cloudwatch.New(awsSession)
	if cloudwatchService == nil {

	}
	//removeDashboard()
	response, err := cloudwatchService.PutDashboard(&cloudwatch.PutDashboardInput{
		DashboardBody: aws.String(widgetsData),
		DashboardName: aws.String("btrz-dynamic"),
	})
	if err != nil {
		fmt.Printf("%v -erorr creating dashboard", err)
	} else {
		fmt.Println("Done", response)
	}
}
func main() {
	awsSession, err := btrzaws.GetAWSSession()
	if err != nil {
		fmt.Print(err, "can't get a session")
		os.Exit(1)
	}
	generateDashboard(awsSession)
}
