package btrzdash

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"strings"
)

// GetLoadBalancersForEnvironment - returns all elb v1 for a tagged environment
func GetLoadBalancersForEnvironment(awsSession *session.Session, envieonment string) ([]*elb.LoadBalancerDescription, error) {
	prodElbs := []*elb.LoadBalancerDescription{}
	elbClient := elb.New(awsSession)
	response, err := elbClient.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}
	for _, currentElb := range response.LoadBalancerDescriptions {
		tags, err := elbClient.DescribeTags(&elb.DescribeTagsInput{
			LoadBalancerNames: []*string{currentElb.LoadBalancerName},
		})
		if err != nil {
			return nil, err
		}
	tagsLoop:
		for _, tags := range tags.TagDescriptions {
			for _, tag := range tags.Tags {
				if *tag.Key == "Environment" && *tag.Value == envieonment {
					prodElbs = append(prodElbs, currentElb)
					break tagsLoop
				}
			}
		}
	}
	return prodElbs, nil
}

// GetApplicationLoadBalancersForEnvironment - return all loadbalancers version 2 that are tagged for an environment
func GetApplicationLoadBalancersForEnvironment(awsSession *session.Session, envieonment string) ([]*elbv2.LoadBalancer, error) {
	prodElbs := []*elbv2.LoadBalancer{}
	elbClient := elbv2.New(awsSession)
	response, err := elbClient.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}
	for _, elb := range response.LoadBalancers {
		tags, err := elbClient.DescribeTags(&elbv2.DescribeTagsInput{
			ResourceArns: []*string{elb.LoadBalancerArn},
		})
		if err != nil {
			return nil, err
		}
	tagsLoop:
		for _, tags := range tags.TagDescriptions {
			for _, tag := range tags.Tags {
				if *tag.Key == "Environment" && *tag.Value == envieonment {
					prodElbs = append(prodElbs, elb)
					break tagsLoop
				}
			}
		}
	}
	return prodElbs, nil
}

// AbbreviateLBArn - remove control part of arn
func AbbreviateLBArn(arn string) string {
	location := strings.Index(arn, "/")
	return arn[location+1:]
}

// BuildDashboardJSONDataFromElb - return dashboard json string for elbv2 slice
func BuildDashboardJSONDataFromElb(loadBalancersV2 []*elbv2.LoadBalancer, loadBalancersV1 []*elb.LoadBalancerDescription) (string, error) {
	elb500Http := ""
	if loadBalancersV2 != nil {
		for _, elb := range loadBalancersV2 {
			elb500Http += fmt.Sprintf(ALB500Patteren, AbbreviateLBArn(*elb.LoadBalancerArn)) + ","
		}
	}
	if loadBalancersV1 != nil {
		for _, elb := range loadBalancersV1 {
			elb500Http += fmt.Sprintf(ELB500Patteren, *elb.LoadBalancerName) + ","
		}
	}
	elb500Http = strings.TrimSuffix(elb500Http, ",")
	widget500 := fmt.Sprintf(WidgetPrefix, 0) + elb500Http + fmt.Sprintf(WidgetPostfix, "Sum", "500 metrics")

	elb400Http := ""
	if loadBalancersV2 != nil {
		for _, elb := range loadBalancersV2 {
			elb400Http += fmt.Sprintf(ALB400Patteren, AbbreviateLBArn(*elb.LoadBalancerArn)) + ","
		}
	}
	if loadBalancersV1 != nil {
		for _, elb := range loadBalancersV1 {
			elb400Http += fmt.Sprintf(ELB400Patteren, *elb.LoadBalancerName) + ","
		}
	}
	elb400Http = strings.TrimSuffix(elb400Http, ",")
	widget400 := fmt.Sprintf(WidgetPrefix, 0) + elb400Http + fmt.Sprintf(WidgetPostfix, "Sum", "400 metrics")

	elb200Http := ""
	if loadBalancersV2 != nil {
		for _, elb := range loadBalancersV2 {
			elb200Http += fmt.Sprintf(ALB200Patteren, AbbreviateLBArn(*elb.LoadBalancerArn)) + ","
		}
	}
	if loadBalancersV1 != nil {
		for _, elb := range loadBalancersV1 {
			elb200Http += fmt.Sprintf(ELB200Patteren, *elb.LoadBalancerName) + ","
		}
	}
	elb200Http = strings.TrimSuffix(elb200Http, ",")
	widget200 := fmt.Sprintf(WidgetPrefix, 0) + elb200Http + fmt.Sprintf(WidgetPostfix, "Sum", "200 metrics")

	return widget500 + "," + widget400 + "," + widget200, nil
}

// Generate500ElbJSON - generate json for 500 error
func Generate500ElbJSON(awsSession *session.Session, envieonment string) (string, error) {
	loadBalancersV2, err := GetApplicationLoadBalancersForEnvironment(awsSession, envieonment)
	if err != nil {
		return "", err
	}
	loadBalancersV1, err := GetLoadBalancersForEnvironment(awsSession, envieonment)
	if err != nil {
		return "", err
	}
	return BuildDashboardJSONDataFromElb(loadBalancersV2, loadBalancersV1)
}
