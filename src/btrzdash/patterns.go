package btrzdash

const (
	versionNumber = 1
	// WidgetsPrefix - WidgetsPrefix
	WidgetsPrefix = `"widgets":[`
	// WidgetsPostfix -WidgetsPostfix
	WidgetsPostfix = "]"
	WidgetPrefix   = `
			{
				 "type":"metric",
				 "x":%d,
				 "y":0,
				 "width":24,
				 "height":8,
				 "properties":{
						"metrics":[`
	WidgetPostfix = `],
		"period":300,
		"stat":"%s",
		"region":"us-east-1",
		"title":"%s"
 }
}`
	MemoryPattern = `[
"%s-memory",
"Total Memory",
"Ram Percentage",
"Used Memory Percentage"
],`
	CPUPattern = `   [
			"AWS/EC2",
			"CPUUtilization",
			"InstanceId",
			"%s"
	 ],`
	ALB500Patteren = `[
		"AWS/ApplicationELB",
		"HTTPCode_ELB_5XX_Count",
		"LoadBalancer",
		"%s"
	]`
	ELB500Patteren = `[
		"AWS/ELB",
		"HTTPCode_Backend_5XX",
		"LoadBalancerName",
		"%s"
	]`
	ALB400Patteren = `[
		"AWS/ApplicationELB",
		"HTTPCode_ELB_4XX_Count",
		"LoadBalancer",
		"%s"
	]`
	ELB400Patteren = `[
		"AWS/ELB",
		"HTTPCode_ELB_4XX_Count",
		"LoadBalancerName",
		"%s"
	]`
	ALB200Patteren = `[
		"AWS/ApplicationELB",
		"HTTPCode_Target_2XX_Count",
		"LoadBalancer",
		"%s"
	]`
	ELB200Patteren = `[
		"AWS/ELB",
		"HTTPCode_Backend_2XX",
		"LoadBalancerName",
		"%s"
	]`
)
