package aws

import (
	"fmt"
	"log"
	//"strings"
	//"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	//"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/cloudwatch"
)

func resourceAwsCloudWatchMetricAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCloudWatchMetricAlarmCreate,
		Read:   resourceAwsCloudWatchMetricAlarmRead,
		Update: resourceAwsCloudWatchMetricAlarmUpdate,
		Delete: resourceAwsCloudWatchMetricAlarmDelete,

		Schema: map[string]*schema.Schema{
			"alarm_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comparison_operator": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"evaluation_periods": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"period": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"statistic": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"threshold": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
			"actions_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"alarm_actions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"alarm_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dimension": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"insufficient_data_actions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"ok_actions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"unit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAwsCloudWatchMetricAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cloudwatchconn

	params := getAwsCloudWatchPutMetricAlarmInput(d)

	log.Printf("[DEBUG] CloudWatch PutMetricAlarm: %#v", params)
	_, err := conn.PutMetricAlarm(&params)
	if err != nil {
		return fmt.Errorf("Error putting metric alarm: %s", err)
	}
	d.SetId(d.Get("alarm_name").(string))
	log.Println("[INFO] CloudWatch Metric Alarm Created")

	return resourceAwsCloudWatchMetricAlarmRead(d, meta)
}

func resourceAwsCloudWatchMetricAlarmRead(d *schema.ResourceData, meta interface{}) error {
	a, err := getAwsCloudWatchMetricAlarm(d, meta)
	if err != nil {
		return err
	}
	if a == nil {
		return nil
	}

	log.Printf("[DEBUG] Read Metric Alarm: %s", d.Get("alarm_name"))

	d.Set("actions_enabled", a.ActionsEnabled)
	d.Set("alarm_actions", a.AlarmActions)
	d.Set("alarm_description", a.AlarmDescription)
	d.Set("alarm_name", a.AlarmName)
	d.Set("comparison_operator", a.ComparisonOperator)
	d.Set("dimensions", a.Dimensions)
	d.Set("evaluation_periods", a.EvaluationPeriods)
	d.Set("insufficient_data_actions", a.InsufficientDataActions)
	d.Set("metric_name", a.MetricName)
	d.Set("namespace", a.Namespace)
	d.Set("ok_actions", a.OKActions)
	d.Set("period", a.Period)
	d.Set("statistic", a.Statistic)
	d.Set("threshold", a.Threshold)
	d.Set("unit", a.Unit)

	return nil
}

func resourceAwsCloudWatchMetricAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cloudwatchconn
	params := getAwsCloudWatchPutMetricAlarmInput(d)

	log.Printf("[DEBUG] CloudWatch PutMetricAlarm: %#v", params)
	_, err := conn.PutMetricAlarm(&params)
	if err != nil {
		return fmt.Errorf("Error putting metric alarm: %s", err)
	}
	log.Println("[INFO] CloudWatch Metric Alarm Updated")

	return resourceAwsCloudWatchMetricAlarmRead(d, meta)
}

func resourceAwsCloudWatchMetricAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	p, err := getAwsCloudWatchMetricAlarm(d, meta)
	if err != nil {
		return err
	}
	if p == nil {
		return nil
	}

	conn := meta.(*AWSClient).cloudwatchconn
	params := cloudwatch.DeleteAlarmsInput{
		AlarmNames: []*string{aws.String(d.Id())},
	}

	if _, err := conn.DeleteAlarms(&params); err != nil {
		cloudwatcherr, ok := err.(aws.APIError)
		if ok && cloudwatcherr.Code != "" {
			fmt.Println("Error: ", cloudwatcherr.Code, cloudwatcherr.Message)
			return nil
		}
		return err
	}

	d.SetId("")
	return nil
}

func getAwsCloudWatchPutMetricAlarmInput(d *schema.ResourceData) cloudwatch.PutMetricAlarmInput {
	params := cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(d.Get("alarm_name").(string)),
		ComparisonOperator: aws.String(d.Get("comparison_operator").(string)),
		EvaluationPeriods:  aws.Long(int64(d.Get("evaluation_periods").(int))),
		MetricName:         aws.String(d.Get("metric_name").(string)),
		Namespace:          aws.String(d.Get("namespace").(string)),
		Period:             aws.Long(int64(d.Get("period").(int))),
		Statistic:          aws.String(d.Get("statistic").(string)),
		Threshold:          aws.Double(d.Get("threshold").(float64)),
	}

	if v := d.Get("actions_enabled"); v != nil {
		params.ActionsEnabled = aws.Boolean(v.(bool))
	}

	if v, ok := d.GetOk("alarm_description"); ok {
		params.AlarmDescription = aws.String(v.(string))
	}

	if v, ok := d.GetOk("unit"); ok {
		params.Unit = aws.String(v.(string))
	}

	var alarmActions []*string
	if v := d.Get("alarm_actions"); v != nil {
		for _, v := range v.(*schema.Set).List() {
			str := v.(string)
			alarmActions = append(alarmActions, aws.String(str))
		}
		params.AlarmActions = alarmActions
	}

	var insufficientDataActions []*string
	if v := d.Get("insufficient_data_actions"); v != nil {
		for _, v := range v.(*schema.Set).List() {
			str := v.(string)
			insufficientDataActions = append(insufficientDataActions, aws.String(str))
		}
		params.InsufficientDataActions = insufficientDataActions
	}

	var okActions []*string
	if v := d.Get("ok_actions"); v != nil {
		for _, v := range v.(*schema.Set).List() {
			str := v.(string)
			okActions = append(okActions, aws.String(str))
		}
		params.OKActions = okActions
	}

	a := d.Get("dimension").(map[string]interface{})
	dimensions := make([]*cloudwatch.Dimension, 0, len(a))
	for k, v := range a {
		dimensions = append(dimensions, &cloudwatch.Dimension{
			Name:  aws.String(k),
			Value: aws.String(v.(string)),
		})
	}
	params.Dimensions = dimensions

	return params
}

func getAwsCloudWatchMetricAlarm(d *schema.ResourceData, meta interface{}) (*cloudwatch.MetricAlarm, error) {
	conn := meta.(*AWSClient).cloudwatchconn

	params := cloudwatch.DescribeAlarmsInput{
		AlarmNames: []*string{aws.String(d.Id())},
	}

	resp, err := conn.DescribeAlarms(&params)
	if err != nil {
		fmt.Errorf("Error retrieve metric alarm: %#v", params)
		_, ok := err.(aws.APIError)
		if ok {
			d.SetId("")
			return nil, nil
		}
		return nil, nil
	}
	for idx, ma := range resp.MetricAlarms {
		if *ma.AlarmName == d.Id() {
			return resp.MetricAlarms[idx], nil
		}
	}
	d.SetId("")
	return nil, nil
}
