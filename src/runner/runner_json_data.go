package main

import (
	"btrzaws"
	"btrzdash"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"strings"
)

func generateJSONData(instances []*btrzaws.BetterezInstance, awsSession *session.Session) string {
	dashboardJSON := btrzdash.WidgetsPrefix
	cpuUsageApp := ""
	cpuUsageAPI := ""
	cpuUsageMongo := ""
	memUsageMongo := ""
	memUsageApp := ""
	memUsageAPI := ""
	for _, instance := range instances {
		if instance.IsAPIInstance() {
			cpuUsageAPI += fmt.Sprintf(btrzdash.CPUPattern, instance.InstanceID)
			memUsageAPI += fmt.Sprintf(btrzdash.MemoryPattern, instance.InstanceName)
		} else if instance.IsMongoInstance() {
			cpuUsageMongo += fmt.Sprintf(btrzdash.CPUPattern, instance.InstanceID)
			memUsageMongo += fmt.Sprintf(btrzdash.MemoryPattern, instance.InstanceName)
		} else if instance.IsAppInstance() {
			memUsageApp += fmt.Sprintf(btrzdash.MemoryPattern, instance.InstanceName)
			cpuUsageApp += fmt.Sprintf(btrzdash.CPUPattern, instance.InstanceID)
		}
	}
	log.Println("done stringing data")
	cpuUsageAPI = strings.TrimSuffix(cpuUsageAPI, ",")
	memUsageAPI = strings.TrimSuffix(memUsageAPI, ",")
	cpuUsageApp = strings.TrimSuffix(cpuUsageApp, ",")
	memUsageApp = strings.TrimSuffix(memUsageApp, ",")
	cpuUsageMongo = strings.TrimSuffix(cpuUsageMongo, ",")
	memUsageMongo = strings.TrimSuffix(memUsageMongo, ",")
	cpuUsageAPIMetric := fmt.Sprintf(btrzdash.WidgetPrefix, 0) + cpuUsageAPI + fmt.Sprintf(btrzdash.WidgetPostfix, "Average", "API CPU metrics")
	memUsageAPIMetric := fmt.Sprintf(btrzdash.WidgetPrefix, 0) + memUsageAPI + fmt.Sprintf(btrzdash.WidgetPostfix, "Average", "API mem metrics")
	cpuUsageMongoMetric := fmt.Sprintf(btrzdash.WidgetPrefix, 0) + cpuUsageMongo + fmt.Sprintf(btrzdash.WidgetPostfix, "Average", "Mongo CPU metrics")
	memUsageMongoMetric := fmt.Sprintf(btrzdash.WidgetPrefix, 0) + memUsageMongo + fmt.Sprintf(btrzdash.WidgetPostfix, "Average", "Mongo mem metrics")
	cpuUsageAppMetric := fmt.Sprintf(btrzdash.WidgetPrefix, 0) + cpuUsageApp + fmt.Sprintf(btrzdash.WidgetPostfix, "Average", "App CPU metrics")
	memUsageAppMetric := fmt.Sprintf(btrzdash.WidgetPrefix, 0) + memUsageApp + fmt.Sprintf(btrzdash.WidgetPostfix, "Average", "App mem metrics")
	elbString, _ := btrzdash.Generate500ElbJSON(awsSession, "production")
	dashboardJSON += cpuUsageAPIMetric + ","
	dashboardJSON += memUsageAPIMetric + ","
	dashboardJSON += cpuUsageMongoMetric + ","
	dashboardJSON += memUsageMongoMetric + ","
	dashboardJSON += cpuUsageAppMetric + ","
	dashboardJSON += elbString + ","
	dashboardJSON += memUsageAppMetric
	if staticData := loadStaticReports(); staticData != "" {
		dashboardJSON += "," + staticData
	}

	dashboardJSON += btrzdash.WidgetsPostfix
	return "{" + dashboardJSON + "}"
	//return dashboardJSON
}
