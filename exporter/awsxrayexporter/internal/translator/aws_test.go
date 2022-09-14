// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package translator

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"

	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
)

func TestAwsFromEc2Resource(t *testing.T) {
	instanceID := "i-00f7c0bcb26da2a99"
	hostType := "m5.xlarge"
	imageID := "ami-0123456789"
	resource := pcommon.NewResource()
	attrs := pcommon.NewMap()
	attrs.UpsertString(conventions.AttributeCloudProvider, conventions.AttributeCloudProviderAWS)
	attrs.UpsertString(conventions.AttributeCloudPlatform, conventions.AttributeCloudPlatformAWSEC2)
	attrs.UpsertString(conventions.AttributeCloudAccountID, "123456789")
	attrs.UpsertString(conventions.AttributeCloudAvailabilityZone, "us-east-1c")
	attrs.UpsertString(conventions.AttributeHostID, instanceID)
	attrs.UpsertString(conventions.AttributeHostType, hostType)
	attrs.UpsertString(conventions.AttributeHostImageID, imageID)
	attrs.CopyTo(resource.Attributes())

	attributes := make(map[string]pcommon.Value)

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.NotNil(t, awsData.EC2)
	assert.Nil(t, awsData.ECS)
	assert.Nil(t, awsData.Beanstalk)
	assert.Nil(t, awsData.EKS)
	assert.Equal(t, "123456789", *awsData.AccountID)
	assert.Equal(t, &awsxray.EC2Metadata{
		InstanceID:       aws.String(instanceID),
		AvailabilityZone: aws.String("us-east-1c"),
		InstanceSize:     aws.String(hostType),
		AmiID:            aws.String(imageID),
	}, awsData.EC2)
}

func TestAwsFromEcsResource(t *testing.T) {
	instanceID := "i-00f7c0bcb26da2a99"
	containerName := "signup_aggregator-x82ufje83"
	containerID := "0123456789A"
	az := "us-east-1c"
	launchType := "fargate"
	family := "family"
	taskArn := "arn:aws:ecs:us-west-2:123456789123:task/123"
	clusterArn := "arn:aws:ecs:us-west-2:123456789123:cluster/my-cluster"
	containerArn := "arn:aws:ecs:us-west-2:123456789123:container-instance/123"
	resource := pcommon.NewResource()
	attrs := pcommon.NewMap()
	attrs.UpsertString(conventions.AttributeCloudProvider, conventions.AttributeCloudProviderAWS)
	attrs.UpsertString(conventions.AttributeCloudPlatform, conventions.AttributeCloudPlatformAWSECS)
	attrs.UpsertString(conventions.AttributeCloudAccountID, "123456789")
	attrs.UpsertString(conventions.AttributeCloudAvailabilityZone, az)
	attrs.UpsertString(conventions.AttributeContainerImageName, "otel/signupaggregator")
	attrs.UpsertString(conventions.AttributeContainerImageTag, "v1")
	attrs.UpsertString(conventions.AttributeContainerName, containerName)
	attrs.UpsertString(conventions.AttributeContainerID, containerID)
	attrs.UpsertString(conventions.AttributeHostID, instanceID)
	attrs.UpsertString(conventions.AttributeAWSECSClusterARN, clusterArn)
	attrs.UpsertString(conventions.AttributeAWSECSContainerARN, containerArn)
	attrs.UpsertString(conventions.AttributeAWSECSTaskARN, taskArn)
	attrs.UpsertString(conventions.AttributeAWSECSTaskFamily, family)
	attrs.UpsertString(conventions.AttributeAWSECSLaunchtype, launchType)
	attrs.UpsertString(conventions.AttributeHostType, "m5.xlarge")

	attrs.CopyTo(resource.Attributes())

	attributes := make(map[string]pcommon.Value)

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.NotNil(t, awsData.ECS)
	assert.NotNil(t, awsData.EC2)
	assert.Nil(t, awsData.Beanstalk)
	assert.Nil(t, awsData.EKS)
	assert.Equal(t, &awsxray.ECSMetadata{
		ContainerName:    aws.String(containerName),
		ContainerID:      aws.String(containerID),
		AvailabilityZone: aws.String(az),
		ClusterArn:       aws.String(clusterArn),
		ContainerArn:     aws.String(containerArn),
		TaskArn:          aws.String(taskArn),
		TaskFamily:       aws.String(family),
		LaunchType:       aws.String(launchType),
	}, awsData.ECS)
}

func TestAwsFromBeanstalkResource(t *testing.T) {
	deployID := "232"
	versionLabel := "4"
	resource := pcommon.NewResource()
	attrs := pcommon.NewMap()
	attrs.UpsertString(conventions.AttributeCloudProvider, conventions.AttributeCloudProviderAWS)
	attrs.UpsertString(conventions.AttributeCloudPlatform, conventions.AttributeCloudPlatformAWSElasticBeanstalk)
	attrs.UpsertString(conventions.AttributeCloudAccountID, "123456789")
	attrs.UpsertString(conventions.AttributeCloudAvailabilityZone, "us-east-1c")
	attrs.UpsertString(conventions.AttributeServiceNamespace, "production")
	attrs.UpsertString(conventions.AttributeServiceInstanceID, deployID)
	attrs.UpsertString(conventions.AttributeServiceVersion, versionLabel)
	attrs.CopyTo(resource.Attributes())

	attributes := make(map[string]pcommon.Value)

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Nil(t, awsData.EC2)
	assert.Nil(t, awsData.ECS)
	assert.NotNil(t, awsData.Beanstalk)
	assert.Nil(t, awsData.EKS)
	assert.Equal(t, &awsxray.BeanstalkMetadata{
		Environment:  aws.String("production"),
		VersionLabel: aws.String(versionLabel),
		DeploymentID: aws.Int64(232),
	}, awsData.Beanstalk)
}

func TestAwsFromEksResource(t *testing.T) {
	instanceID := "i-00f7c0bcb26da2a99"
	containerName := "signup_aggregator-x82ufje83"
	containerID := "0123456789A"
	resource := pcommon.NewResource()
	attrs := pcommon.NewMap()
	attrs.UpsertString(conventions.AttributeCloudProvider, conventions.AttributeCloudProviderAWS)
	attrs.UpsertString(conventions.AttributeCloudPlatform, conventions.AttributeCloudPlatformAWSEKS)
	attrs.UpsertString(conventions.AttributeCloudAccountID, "123456789")
	attrs.UpsertString(conventions.AttributeCloudAvailabilityZone, "us-east-1c")
	attrs.UpsertString(conventions.AttributeContainerImageName, "otel/signupaggregator")
	attrs.UpsertString(conventions.AttributeContainerImageTag, "v1")
	attrs.UpsertString(conventions.AttributeK8SClusterName, "production")
	attrs.UpsertString(conventions.AttributeK8SNamespaceName, "default")
	attrs.UpsertString(conventions.AttributeK8SDeploymentName, "signup_aggregator")
	attrs.UpsertString(conventions.AttributeK8SPodName, "my-deployment-65dcf7d447-ddjnl")
	attrs.UpsertString(conventions.AttributeContainerName, containerName)
	attrs.UpsertString(conventions.AttributeContainerID, containerID)
	attrs.UpsertString(conventions.AttributeHostID, instanceID)
	attrs.UpsertString(conventions.AttributeHostType, "m5.xlarge")
	attrs.CopyTo(resource.Attributes())

	attributes := make(map[string]pcommon.Value)

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.NotNil(t, awsData.EKS)
	assert.NotNil(t, awsData.EC2)
	assert.Nil(t, awsData.ECS)
	assert.Nil(t, awsData.Beanstalk)
	assert.Equal(t, &awsxray.EKSMetadata{
		ClusterName: aws.String("production"),
		Pod:         aws.String("my-deployment-65dcf7d447-ddjnl"),
		ContainerID: aws.String(containerID),
	}, awsData.EKS)
}

func TestAwsWithAwsSqsResources(t *testing.T) {
	instanceID := "i-00f7c0bcb26da2a99"
	containerName := "signup_aggregator-x82ufje83"
	containerID := "0123456789A"
	resource := pcommon.NewResource()
	attrs := pcommon.NewMap()
	attrs.UpsertString(conventions.AttributeCloudProvider, conventions.AttributeCloudProviderAWS)
	attrs.UpsertString(conventions.AttributeCloudAccountID, "123456789")
	attrs.UpsertString(conventions.AttributeCloudAvailabilityZone, "us-east-1c")
	attrs.UpsertString(conventions.AttributeContainerName, containerName)
	attrs.UpsertString(conventions.AttributeContainerImageName, "otel/signupaggregator")
	attrs.UpsertString(conventions.AttributeContainerImageTag, "v1")
	attrs.UpsertString(conventions.AttributeK8SClusterName, "production")
	attrs.UpsertString(conventions.AttributeK8SNamespaceName, "default")
	attrs.UpsertString(conventions.AttributeK8SDeploymentName, "signup_aggregator")
	attrs.UpsertString(conventions.AttributeK8SPodName, "my-deployment-65dcf7d447-ddjnl")
	attrs.UpsertString(conventions.AttributeContainerName, containerName)
	attrs.UpsertString(conventions.AttributeContainerID, containerID)
	attrs.UpsertString(conventions.AttributeHostID, instanceID)
	attrs.UpsertString(conventions.AttributeHostType, "m5.xlarge")

	queueURL := "https://sqs.use1.amazonaws.com/Meltdown-Alerts"
	attributes := make(map[string]pcommon.Value)
	attributes[awsxray.AWSOperationAttribute] = pcommon.NewValueString("SendMessage")
	attributes[awsxray.AWSAccountAttribute] = pcommon.NewValueString("987654321")
	attributes[awsxray.AWSRegionAttribute] = pcommon.NewValueString("us-east-2")
	attributes[awsxray.AWSQueueURLAttribute] = pcommon.NewValueString(queueURL)
	attributes["employee.id"] = pcommon.NewValueString("XB477")

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, queueURL, *awsData.QueueURL)
	assert.Equal(t, "us-east-2", *awsData.RemoteRegion)
}

func TestAwsWithRpcAttributes(t *testing.T) {
	resource := pcommon.NewResource()
	attributes := make(map[string]pcommon.Value)
	attributes[conventions.AttributeRPCMethod] = pcommon.NewValueString("ListBuckets")

	_, awsData := makeAws(attributes, resource)

	assert.NotNil(t, awsData)
	assert.Equal(t, "ListBuckets", *awsData.Operation)
}

func TestAwsWithSqsAlternateAttribute(t *testing.T) {
	queueURL := "https://sqs.use1.amazonaws.com/Meltdown-Alerts"
	attributes := make(map[string]pcommon.Value)
	attributes[awsxray.AWSQueueURLAttribute2] = pcommon.NewValueString(queueURL)

	filtered, awsData := makeAws(attributes, pcommon.NewResource())

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, queueURL, *awsData.QueueURL)
}

func TestAwsWithAwsDynamoDbResources(t *testing.T) {
	instanceID := "i-00f7c0bcb26da2a99"
	containerName := "signup_aggregator-x82ufje83"
	containerID := "0123456789A"
	resource := pcommon.NewResource()
	attrs := pcommon.NewMap()
	attrs.UpsertString(conventions.AttributeCloudProvider, conventions.AttributeCloudProviderAWS)
	attrs.UpsertString(conventions.AttributeCloudAccountID, "123456789")
	attrs.UpsertString(conventions.AttributeCloudAvailabilityZone, "us-east-1c")
	attrs.UpsertString(conventions.AttributeContainerName, "signup_aggregator")
	attrs.UpsertString(conventions.AttributeContainerImageName, "otel/signupaggregator")
	attrs.UpsertString(conventions.AttributeContainerImageTag, "v1")
	attrs.UpsertString(conventions.AttributeK8SClusterName, "production")
	attrs.UpsertString(conventions.AttributeK8SNamespaceName, "default")
	attrs.UpsertString(conventions.AttributeK8SDeploymentName, "signup_aggregator")
	attrs.UpsertString(conventions.AttributeK8SPodName, "my-deployment-65dcf7d447-ddjnl")
	attrs.UpsertString(conventions.AttributeContainerName, containerName)
	attrs.UpsertString(conventions.AttributeContainerID, containerID)
	attrs.UpsertString(conventions.AttributeHostID, instanceID)
	attrs.UpsertString(conventions.AttributeHostType, "m5.xlarge")

	tableName := "WIDGET_TYPES"
	attributes := make(map[string]pcommon.Value)
	attributes[conventions.AttributeRPCMethod] = pcommon.NewValueString("IncorrectAWSSDKOperation")
	attributes[awsxray.AWSOperationAttribute] = pcommon.NewValueString("PutItem")
	attributes[awsxray.AWSRequestIDAttribute] = pcommon.NewValueString("75107C82-EC8A-4F75-883F-4440B491B0AB")
	attributes[awsxray.AWSTableNameAttribute] = pcommon.NewValueString(tableName)

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, "PutItem", *awsData.Operation)
	assert.Equal(t, "75107C82-EC8A-4F75-883F-4440B491B0AB", *awsData.RequestID)
	assert.Equal(t, tableName, *awsData.TableName)
}

func TestAwsWithDynamoDbAlternateAttribute(t *testing.T) {
	tableName := "MyTable"
	attributes := make(map[string]pcommon.Value)
	attributes[awsxray.AWSTableNameAttribute2] = pcommon.NewValueString(tableName)

	filtered, awsData := makeAws(attributes, pcommon.NewResource())

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, tableName, *awsData.TableName)
}

func TestAwsWithRequestIdAlternateAttribute(t *testing.T) {
	requestid := "12345-request"
	attributes := make(map[string]pcommon.Value)
	attributes[awsxray.AWSRequestIDAttribute2] = pcommon.NewValueString(requestid)

	filtered, awsData := makeAws(attributes, pcommon.NewResource())

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, requestid, *awsData.RequestID)
}

func TestJavaSDK(t *testing.T) {
	attributes := make(map[string]pcommon.Value)
	resource := pcommon.NewResource()
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKName, "opentelemetry")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKLanguage, "java")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKVersion, "1.2.3")

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, "opentelemetry for java", *awsData.XRay.SDK)
	assert.Equal(t, "1.2.3", *awsData.XRay.SDKVersion)
}

func TestJavaAutoInstrumentation(t *testing.T) {
	attributes := make(map[string]pcommon.Value)
	resource := pcommon.NewResource()
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKName, "opentelemetry")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKLanguage, "java")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKVersion, "1.2.3")
	resource.Attributes().UpsertString(conventions.AttributeTelemetryAutoVersion, "3.4.5")

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, "opentelemetry for java", *awsData.XRay.SDK)
	assert.Equal(t, "1.2.3", *awsData.XRay.SDKVersion)
	assert.True(t, *awsData.XRay.AutoInstrumentation)
}

func TestGoSDK(t *testing.T) {
	attributes := make(map[string]pcommon.Value)
	resource := pcommon.NewResource()
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKName, "opentelemetry")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKLanguage, "go")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKVersion, "2.0.3")

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, "opentelemetry for go", *awsData.XRay.SDK)
	assert.Equal(t, "2.0.3", *awsData.XRay.SDKVersion)
}

func TestCustomSDK(t *testing.T) {
	attributes := make(map[string]pcommon.Value)
	resource := pcommon.NewResource()
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKName, "opentracing")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKLanguage, "java")
	resource.Attributes().UpsertString(conventions.AttributeTelemetrySDKVersion, "2.0.3")

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, "opentracing for java", *awsData.XRay.SDK)
	assert.Equal(t, "2.0.3", *awsData.XRay.SDKVersion)
}

func TestLogGroups(t *testing.T) {
	cwl1 := awsxray.LogGroupMetadata{
		LogGroup: awsxray.String("group1"),
	}
	cwl2 := awsxray.LogGroupMetadata{
		LogGroup: awsxray.String("group2"),
	}

	attributes := make(map[string]pcommon.Value)
	resource := pcommon.NewResource()
	ava := resource.Attributes().UpsertEmptySlice(conventions.AttributeAWSLogGroupNames)
	ava.EnsureCapacity(2)
	ava.AppendEmpty().SetStringVal("group1")
	ava.AppendEmpty().SetStringVal("group2")

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, 2, len(awsData.CWLogs))
	assert.Contains(t, awsData.CWLogs, cwl1)
	assert.Contains(t, awsData.CWLogs, cwl2)
}

func TestLogGroupsFromArns(t *testing.T) {
	group1 := "arn:aws:logs:us-east-1:123456789123:log-group:group1"
	cwl1 := awsxray.LogGroupMetadata{
		LogGroup: awsxray.String("group1"),
		Arn:      awsxray.String(group1),
	}
	group2 := "arn:aws:logs:us-east-1:123456789123:log-group:group2:*"
	cwl2 := awsxray.LogGroupMetadata{
		LogGroup: awsxray.String("group2"),
		Arn:      awsxray.String(group2),
	}

	attributes := make(map[string]pcommon.Value)
	resource := pcommon.NewResource()
	ava := resource.Attributes().UpsertEmptySlice(conventions.AttributeAWSLogGroupARNs)
	ava.EnsureCapacity(2)
	ava.AppendEmpty().SetStringVal(group1)
	ava.AppendEmpty().SetStringVal(group2)

	filtered, awsData := makeAws(attributes, resource)

	assert.NotNil(t, filtered)
	assert.NotNil(t, awsData)
	assert.Equal(t, 2, len(awsData.CWLogs))
	assert.Contains(t, awsData.CWLogs, cwl1)
	assert.Contains(t, awsData.CWLogs, cwl2)
}
