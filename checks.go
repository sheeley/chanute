package chanute

type Check string

const (
	CheckAmazonAuroraDBInstanceAccessibility                       = "Amazon Aurora DB Instance Accessibility"
	CheckAmazonEBSProvisionedIOPSSSDVolumeAttachmentConfiguration  = "Amazon EBS Provisioned IOPS (SSD) Volume Attachment Configuration"
	CheckAmazonEBSPublicSnapshots                                  = "Amazon EBS Public Snapshots"
	CheckAmazonEBSSnapshots                                        = "Amazon EBS Snapshots"
	CheckAmazonEC2AvailabilityZoneBalance                          = "Amazon EC2 Availability Zone Balance"
	CheckAmazonEC2ReservedInstanceLeaseExpiration                  = "Amazon EC2 Reserved Instance Lease Expiration"
	CheckAmazonEC2ReservedInstancesOptimization                    = "Amazon EC2 Reserved Instances Optimization"
	CheckAmazonEC2toEBSThroughputOptimization                      = "Amazon EC2 to EBS Throughput Optimization"
	CheckAmazonRDSBackups                                          = "Amazon RDS Backups"
	CheckAmazonRDSIdleDBInstances                                  = "Amazon RDS Idle DB Instances"
	CheckAmazonRDSMultiAZ                                          = "Amazon RDS Multi-AZ"
	CheckAmazonRDSPublicSnapshots                                  = "Amazon RDS Public Snapshots"
	CheckAmazonRDSSecurityGroupAccessRisk                          = "Amazon RDS Security Group Access Risk"
	CheckAmazonRoute53AliasResourceRecordSets                      = "Amazon Route 53 Alias Resource Record Sets"
	CheckAmazonRoute53DeletedHealthChecks                          = "Amazon Route 53 Deleted Health Checks"
	CheckAmazonRoute53FailoverResourceRecordSets                   = "Amazon Route 53 Failover Resource Record Sets"
	CheckAmazonRoute53HighTTLResourceRecordSets                    = "Amazon Route 53 High TTL Resource Record Sets"
	CheckAmazonRoute53LatencyResourceRecordSets                    = "Amazon Route 53 Latency Resource Record Sets"
	CheckAmazonRoute53MXResourceRecordSetsandSenderPolicyFramework = "Amazon Route 53 MX Resource Record Sets and Sender Policy Framework"
	CheckAmazonRoute53NameServerDelegations                        = "Amazon Route 53 Name Server Delegations"
	CheckAmazonS3BucketLogging                                     = "Amazon S3 Bucket Logging"
	CheckAmazonS3BucketPermissions                                 = "Amazon S3 Bucket Permissions"
	CheckAmazonS3BucketVersioning                                  = "Amazon S3 Bucket Versioning"
	CheckAutoScalingGroupHealthCheck                               = "Auto Scaling Group Health Check"
	CheckAutoScalingGroupResources                                 = "Auto Scaling Group Resources"
	CheckAutoScalingGroups                                         = "Auto Scaling Groups"
	CheckAutoScalingLaunchConfigurations                           = "Auto Scaling Launch Configurations"
	CheckAWSCloudTrailLogging                                      = "AWS CloudTrail Logging"
	CheckAWSDirectConnectConnectionRedundancy                      = "AWS Direct Connect Connection Redundancy"
	CheckAWSDirectConnectLocationRedundancy                        = "AWS Direct Connect Location Redundancy"
	CheckAWSDirectConnectVirtualInterfaceRedundancy                = "AWS Direct Connect Virtual Interface Redundancy"
	CheckCloudFormationStacks                                      = "CloudFormation Stacks"
	CheckCloudFrontAlternateDomainNames                            = "CloudFront Alternate Domain Names"
	CheckCloudFrontContentDeliveryOptimization                     = "CloudFront Content Delivery Optimization"
	CheckCloudFrontCustomSSLCertificatesintheIAMCertificateStore   = "CloudFront Custom SSL Certificates in the IAM Certificate Store"
	CheckCloudFrontHeaderForwardingandCacheHitRatio                = "CloudFront Header Forwarding and Cache Hit Ratio"
	CheckCloudFrontSSLCertificateontheOriginServer                 = "CloudFront SSL Certificate on the Origin Server"
	CheckDynamoDBReadCapacity                                      = "DynamoDB Read Capacity"
	CheckDynamoDBWriteCapacity                                     = "DynamoDB Write Capacity"
	CheckEBSActiveSnapshots                                        = "EBS Active Snapshots"
	CheckEBSActiveVolumes                                          = "EBS Active Volumes"
	CheckEBSColdHDDSC1VolumeStorage                                = "EBS Cold HDD (sc1) Volume Storage"
	CheckEBSGeneralPurposeSSDGP2VolumeStorage                      = "EBS General Purpose SSD (gp2) Volume Storage"
	CheckEBSMagneticStandardVolumeStorage                          = "EBS Magnetic (standard) Volume Storage"
	CheckEBSProvisionedIOPSSSSDVolumeAggregateIOPS                 = "EBS Provisioned IOPS (SSD) Volume Aggregate IOPS"
	CheckEBSProvisionedIOPSSSDIO1VolumeStorage                     = "EBS Provisioned IOPS SSD (io1) Volume Storage"
	CheckEBSThroughputOptimizedHDDST1VolumeStorage                 = "EBS Throughput Optimized HDD (st1) Volume Storage"
	CheckEC2ElasticIPAddresses                                     = "EC2 Elastic IP Addresses"
	CheckEC2OnDemandInstances                                      = "EC2 On-Demand Instances"
	CheckEC2ReservedInstanceLeases                                 = "EC2 Reserved Instance Leases"
	CheckEC2ConfigServiceforEC2WindowsInstances                    = "EC2Config Service for EC2 Windows Instances"
	CheckELBApplicationLoadBalancers                               = "ELB Application Load Balancers"
	CheckELBClassicLoadBalancers                                   = "ELB Classic Load Balancers"
	CheckELBConnectionDraining                                     = "ELB Connection Draining"
	CheckELBCrossZoneLoadBalancing                                 = "ELB Cross-Zone Load Balancing"
	CheckELBListenerSecurity                                       = "ELB Listener Security"
	CheckELBNetworkLoadBalancers                                   = "ELB Network Load Balancers"
	CheckELBSecurityGroups                                         = "ELB Security Groups"
	CheckENADriverVersionforEC2WindowsInstances                    = "ENA Driver Version for EC2 Windows Instances"
	CheckExposedAccessKeys                                         = "Exposed Access Keys"
	CheckHighUtilizationAmazonEC2Instances                         = "High Utilization Amazon EC2 Instances"
	CheckIAMAccessKeyRotation                                      = "IAM Access Key Rotation"
	CheckIAMGroup                                                  = "IAM Group"
	CheckIAMInstanceProfiles                                       = "IAM Instance Profiles"
	CheckIAMPasswordPolicy                                         = "IAM Password Policy"
	CheckIAMPolicies                                               = "IAM Policies"
	CheckIAMRoles                                                  = "IAM Roles"
	CheckIAMServerCertificates                                     = "IAM Server Certificates"
	CheckIAMUse                                                    = "IAM Use"
	CheckIAMUsers                                                  = "IAM Users"
	CheckIdleLoadBalancers                                         = "Idle Load Balancers"
	CheckKinesisShardsperRegion                                    = "Kinesis Shards per Region"
	CheckLargeNumberofEC2SecurityGroupRulesAppliedtoanInstance     = "Large Number of EC2 Security Group Rules Applied to an Instance"
	CheckLargeNumberofRulesinanEC2SecurityGroup                    = "Large Number of Rules in an EC2 Security Group"
	CheckLoadBalancerOptimization                                  = "Load Balancer Optimization"
	CheckLowUtilizationAmazonEC2Instances                          = "Low Utilization Amazon EC2 Instances"
	CheckMFAonRootAccount                                          = "MFA on Root Account"
	CheckNVMeDriverVersionforEC2WindowsInstances                   = "NVMe Driver Version for EC2 Windows Instances"
	CheckOverutilizedAmazonEBSMagneticVolumes                      = "Overutilized Amazon EBS Magnetic Volumes"
	CheckPVDriverVersionforEC2WindowsInstances                     = "PV Driver Version for EC2 Windows Instances"
	CheckRDSClusterParameterGroups                                 = "RDS Cluster Parameter Groups"
	CheckRDSClusterRoles                                           = "RDS Cluster Roles"
	CheckRDSClusters                                               = "RDS Clusters"
	CheckRDSDBInstances                                            = "RDS DB Instances"
	CheckRDSDBManualSnapshots                                      = "RDS DB Manual Snapshots"
	CheckRDSDBParameterGroups                                      = "RDS DB Parameter Groups"
	CheckRDSDBSecurityGroups                                       = "RDS DB Security Groups"
	CheckRDSEventSubscriptions                                     = "RDS Event Subscriptions"
	CheckRDSMaxAuthsperSecurityGroup                               = "RDS Max Auths per Security Group"
	CheckRDSOptionGroups                                           = "RDS Option Groups"
	CheckRDSReadReplicasperMaster                                  = "RDS Read Replicas per Master"
	CheckRDSReservedInstances                                      = "RDS Reserved Instances"
	CheckRDSSubnetGroups                                           = "RDS Subnet Groups"
	CheckRDSSubnetsperSubnetGroup                                  = "RDS Subnets per Subnet Group"
	CheckRDSTotalStorageQuota                                      = "RDS Total Storage Quota"
	CheckRoute53HostedZones                                        = "Route 53 Hosted Zones"
	CheckRoute53MaxHealthChecks                                    = "Route 53 Max Health Checks"
	CheckRoute53ReusableDelegationSets                             = "Route 53 Reusable Delegation Sets"
	CheckRoute53TrafficPolicies                                    = "Route 53 Traffic Policies"
	CheckRoute53TrafficPolicyInstances                             = "Route 53 Traffic Policy Instances"
	CheckSecurityGroupsSpecificPortsUnrestricted                   = "Security Groups - Specific Ports Unrestricted"
	CheckSecurityGroupsUnrestrictedAccess                          = "Security Groups - Unrestricted Access"
	CheckSESDailySendingQuota                                      = "SES Daily Sending Quota"
	CheckUnassociatedElasticIPAddresses                            = "Unassociated Elastic IP Addresses"
	CheckUnderutilizedAmazonEBSVolumes                             = "Underutilized Amazon EBS Volumes"
	CheckUnderutilizedAmazonRedshiftClusters                       = "Underutilized Amazon Redshift Clusters"
	CheckVPC                                                       = "VPC"
	CheckVPCElasticIPAddress                                       = "VPC Elastic IP Address"
	CheckVPCInternetGateways                                       = "VPC Internet Gateways"
	CheckVPNTunnelRedundancy                                       = "VPN Tunnel Redundancy"
)

var costChecks = []Check{
	CheckLowUtilizationAmazonEC2Instances,
	CheckUnderutilizedAmazonEBSVolumes,
	CheckIdleLoadBalancers,
	CheckAmazonRDSIdleDBInstances,
	CheckUnderutilizedAmazonRedshiftClusters,
	CheckAmazonEC2ReservedInstanceLeaseExpiration,
	CheckAmazonRoute53LatencyResourceRecordSets,
	CheckUnassociatedElasticIPAddresses,
	CheckAmazonEC2ReservedInstancesOptimization,
}

var faultToleranceChecks = []Check{
	CheckAmazonEBSSnapshots,
	CheckAmazonRDSBackups,
	CheckAutoScalingGroupResources,
	CheckAmazonEC2AvailabilityZoneBalance,
	CheckAmazonRDSMultiAZ,
	CheckAmazonS3BucketLogging,
	CheckAmazonS3BucketVersioning,
	CheckAutoScalingGroupHealthCheck,
	CheckEC2ConfigServiceforEC2WindowsInstances,
	CheckELBConnectionDraining,
	CheckELBCrossZoneLoadBalancing,
	CheckLoadBalancerOptimization,
	CheckPVDriverVersionforEC2WindowsInstances,
	CheckAmazonAuroraDBInstanceAccessibility,
	CheckAmazonRoute53DeletedHealthChecks,
	CheckAmazonRoute53FailoverResourceRecordSets,
	CheckAmazonRoute53HighTTLResourceRecordSets,
	CheckAmazonRoute53NameServerDelegations,
	CheckAWSDirectConnectVirtualInterfaceRedundancy,
	CheckVPNTunnelRedundancy,
	CheckAWSDirectConnectConnectionRedundancy,
	CheckAWSDirectConnectLocationRedundancy,
	CheckENADriverVersionforEC2WindowsInstances,
	CheckNVMeDriverVersionforEC2WindowsInstances,
}
var performanceChecks = []Check{
	CheckHighUtilizationAmazonEC2Instances,
	CheckLargeNumberofEC2SecurityGroupRulesAppliedtoanInstance,
	CheckLargeNumberofRulesinanEC2SecurityGroup,
	CheckAmazonEBSProvisionedIOPSSSDVolumeAttachmentConfiguration,
	CheckAmazonEC2toEBSThroughputOptimization,
	CheckAmazonRoute53AliasResourceRecordSets,
	CheckCloudFrontAlternateDomainNames,
	CheckCloudFrontContentDeliveryOptimization,
	CheckCloudFrontHeaderForwardingandCacheHitRatio,
	CheckOverutilizedAmazonEBSMagneticVolumes,
}

var securityChecks = []Check{
	CheckIAMAccessKeyRotation,
	CheckSecurityGroupsSpecificPortsUnrestricted,
	CheckSecurityGroupsUnrestrictedAccess,
	CheckAmazonS3BucketPermissions,
	CheckELBListenerSecurity,
	CheckELBSecurityGroups,
	CheckIAMPasswordPolicy,
	CheckAmazonEBSPublicSnapshots,
	CheckAmazonRDSPublicSnapshots,
	CheckAmazonRDSSecurityGroupAccessRisk,
	CheckAmazonRoute53MXResourceRecordSetsandSenderPolicyFramework,
	CheckAWSCloudTrailLogging,
	CheckCloudFrontCustomSSLCertificatesintheIAMCertificateStore,
	CheckCloudFrontSSLCertificateontheOriginServer,
	CheckExposedAccessKeys,
	CheckIAMUse,
	CheckMFAonRootAccount,
}

var serviceLimitChecks = []Check{
	CheckCloudFormationStacks,
	CheckAutoScalingGroups,
	CheckAutoScalingLaunchConfigurations,
	CheckDynamoDBReadCapacity,
	CheckDynamoDBWriteCapacity,
	CheckEBSActiveSnapshots,
	CheckEBSActiveVolumes,
	CheckEBSColdHDDSC1VolumeStorage,
	CheckEBSGeneralPurposeSSDGP2VolumeStorage,
	CheckEBSMagneticStandardVolumeStorage,
	CheckEBSProvisionedIOPSSSSDVolumeAggregateIOPS,
	CheckEBSProvisionedIOPSSSDIO1VolumeStorage,
	CheckEBSThroughputOptimizedHDDST1VolumeStorage,
	CheckEC2ElasticIPAddresses,
	CheckEC2OnDemandInstances,
	CheckEC2ReservedInstanceLeases,
	CheckELBApplicationLoadBalancers,
	CheckELBClassicLoadBalancers,
	CheckELBNetworkLoadBalancers,
	CheckIAMGroup,
	CheckIAMInstanceProfiles,
	CheckIAMPolicies,
	CheckIAMRoles,
	CheckIAMServerCertificates,
	CheckIAMUsers,
	CheckKinesisShardsperRegion,
	CheckRDSClusterParameterGroups,
	CheckRDSClusterRoles,
	CheckRDSClusters,
	CheckRDSDBInstances,
	CheckRDSDBManualSnapshots,
	CheckRDSDBParameterGroups,
	CheckRDSDBSecurityGroups,
	CheckRDSEventSubscriptions,
	CheckRDSMaxAuthsperSecurityGroup,
	CheckRDSOptionGroups,
	CheckRDSReadReplicasperMaster,
	CheckRDSReservedInstances,
	CheckRDSSubnetGroups,
	CheckRDSSubnetsperSubnetGroup,
	CheckRDSTotalStorageQuota,
	CheckRoute53HostedZones,
	CheckRoute53MaxHealthChecks,
	CheckRoute53ReusableDelegationSets,
	CheckRoute53TrafficPolicies,
	CheckRoute53TrafficPolicyInstances,
	CheckSESDailySendingQuota,
	CheckVPC,
	CheckVPCElasticIPAddress,
	CheckVPCInternetGateways,
}

type CheckType string

const (
	CheckTypeCost           CheckType = "CheckTypeCost"
	CheckTypeFaultTolerance CheckType = "CheckTypeFaultTolerance"
	CheckTypePerformance    CheckType = "CheckTypePerformance"
	CheckTypeSecurity       CheckType = "CheckTypeSecurity"
	CheckTypeServiceLimit   CheckType = "CheckTypeServiceLimit"
)

var checkTypeLookup map[Check]CheckType
var typeMap = map[CheckType][]Check{
	CheckTypeCost:           costChecks,
	CheckTypeFaultTolerance: faultToleranceChecks,
	CheckTypePerformance:    performanceChecks,
	CheckTypeSecurity:       securityChecks,
	CheckTypeServiceLimit:   serviceLimitChecks,
}

func init() {
	checkTypeLookup = make(map[Check]CheckType)
	for t, checks := range typeMap {
		for _, c := range checks {
			checkTypeLookup[c] = t
		}
	}
}
