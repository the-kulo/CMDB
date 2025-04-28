package azure

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysqlflexibleservers"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// Azure资源的结构
type Resource struct {
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Location string            `json:"location"`
	Owner    string            `json:"owner"`
	Type     string            `json:"type"`
	Tags     map[string]string `json:"tags"`
}

// Azure虚拟机
type VMResource struct {
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Location string            `json:"location"`
	Owner    string            `json:"owner"`
	Size     string            `json:"size,omitempty"`
	Status   string            `json:"status,omitempty"`
	Tags     map[string]string `json:"tags"`
}

// Azure数据库
type DBResource struct {
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Location string            `json:"location"`
	Owner    string            `json:"owner"`
	Server   string            `json:"server,omitempty"`
	DBType   string            `json:"dbType,omitempty"`
	Version  string            `json:"version,omitempty"`
	Status   string            `json:"status,omitempty"`
	Tags     map[string]string `json:"tags"`
}

// AzureHelper 封装Azure认证和资源获取功能
type AzureHelper struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	graphClient            *msgraphsdk.GraphServiceClient
	subscriptionID         string
}

// NewAzureHelper 创建新的AzureHelper实例
func NewAzureHelper() *AzureHelper {
	return &AzureHelper{}
}

// Initialize 初始化Azure认证
func (a *AzureHelper) Initialize() error {
	clientID := os.Getenv("CLIENT_ID")
	tenantID := os.Getenv("TENANT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	a.subscriptionID = os.Getenv("SUBSCRIPTION_ID")

	if clientID == "" || tenantID == "" || clientSecret == "" || a.subscriptionID == "" {
		return fmt.Errorf("环境变量未设置: 请确保设置了CLIENT_ID, TENANT_ID, CLIENT_SECRET和SUBSCRIPTION_ID")
	}

	credOpts := azidentity.ClientSecretCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureChina,
		},
	}

	credential, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, &credOpts)
	if err != nil {
		return fmt.Errorf("创建Azure凭证失败: %v", err)
	}

	a.clientSecretCredential = credential

	// 创建认证提供者
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(a.clientSecretCredential, []string{
		"https://management.chinacloudapi.cn/.default",
	})
	if err != nil {
		return fmt.Errorf("创建认证提供者失败: %v", err)
	}

	// 创建请求适配器
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return fmt.Errorf("创建请求适配器失败: %v", err)
	}

	// 设置中国云端点
	adapter.SetBaseUrl("https://microsoftgraph.chinacloudapi.cn/v1.0")

	// 创建Graph客户端
	a.graphClient = msgraphsdk.NewGraphServiceClient(adapter)

	return nil
}

// GetToken 获取Azure访问令牌（用于调试）
func (a *AzureHelper) GetToken() (string, error) {
	if a.clientSecretCredential == nil {
		if err := a.Initialize(); err != nil {
			return "", err
		}
	}

	token, err := a.clientSecretCredential.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://management.chinacloudapi.cn/.default"},
	})
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %v", err)
	}

	return token.Token, nil
}

// GetResources 获取Azure资源列表
func (a *AzureHelper) GetResources() ([]Resource, error) {
	if a.clientSecretCredential == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	// 创建资源客户端工厂
	clientFactory, err := armresources.NewClientFactory(a.subscriptionID, a.clientSecretCredential, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureChina,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("创建资源客户端工厂失败: %v", err)
	}

	// 获取资源客户端
	resourcesClient := clientFactory.NewClient()
	pager := resourcesClient.NewListPager(nil)

	var resources []Resource

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("获取资源列表失败: %v", err)
		}

		for _, item := range page.Value {
			owner := ""
			if item.Tags != nil && item.Tags["owner"] != nil {
				owner = *item.Tags["owner"]
			}

			resource := Resource{
				Name:     *item.Name,
				ID:       strings.Split(*item.ID, "/")[2],
				Location: *item.Location,
				Owner:    owner,
				Type:     strings.Split(*item.Type, "/")[len(strings.Split(*item.Type, "/"))-1],
				Tags:     make(map[string]string),
			}

			if item.Tags != nil {
				for k, v := range item.Tags {
					if v != nil {
						resource.Tags[k] = *v
					}
				}
			}

			resources = append(resources, resource)
		}
	}

	return resources, nil
}

// GetVirtualMachines 获取Azure虚拟机资源列表
func (a *AzureHelper) GetVirtualMachines() ([]VMResource, error) {
	if a.clientSecretCredential == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	// 创建计算客户端工厂
	clientFactory, err := armcompute.NewClientFactory(a.subscriptionID, a.clientSecretCredential, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzureChina,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("创建计算客户端工厂失败: %v", err)
	}

	// 获取虚拟机客户端
	vmClient := clientFactory.NewVirtualMachinesClient()
	pager := vmClient.NewListAllPager(nil)

	var vms []VMResource

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("获取虚拟机列表失败: %v", err)
		}

		for _, vm := range page.Value {
			owner := ""
			if vm.Tags != nil && vm.Tags["owner"] != nil {
				owner = *vm.Tags["owner"]
			}

			var osType string
			if vm.Properties != nil && vm.Properties.StorageProfile != nil &&
				vm.Properties.StorageProfile.OSDisk != nil && vm.Properties.StorageProfile.OSDisk.OSType != nil {
				osType = string(*vm.Properties.StorageProfile.OSDisk.OSType)
			}

			status := ""
			if vm.Properties != nil && vm.Properties.ProvisioningState != nil {
				status = *vm.Properties.ProvisioningState
			}

			vmResource := VMResource{
				Name:     *vm.Name,
				ID:       strings.Split(*vm.ID, "/")[2],
				Location: *vm.Location,
				Owner:    owner,
				Size:     osType,
				Status:   status,
				Tags:     make(map[string]string),
			}

			if vm.Tags != nil {
				for k, v := range vm.Tags {
					if v != nil {
						vmResource.Tags[k] = *v
					}
				}
			}

			vms = append(vms, vmResource)
		}
	}

	return vms, nil
}

// GetSQLDatabases 获取Azure SQL数据库资源列表
func (a *AzureHelper) GetSQLDatabases() ([]DBResource, error) {
	if a.clientSecretCredential == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	// 创建SQL客户端工厂
	clientFactory, err := armsql.NewClientFactory(
		a.subscriptionID,
		a.clientSecretCredential,
		&arm.ClientOptions{ClientOptions: azcore.ClientOptions{Cloud: cloud.AzureChina}},
	)
	if err != nil {
		return nil, fmt.Errorf("创建SQL客户端工厂失败: %v", err)
	}

	// 获取SQL服务器列表
	serversClient := clientFactory.NewServersClient()
	srvPager := serversClient.NewListPager(nil)

	var sqlDatabases []DBResource
	for srvPager.More() {
		srvPage, err := srvPager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("列举SQL服务器失败: %v", err)
		}
		for _, srv := range srvPage.Value {
			// 从服务器ID中解析出资源组名称
			parts := strings.Split(*srv.ID, "/")
			if len(parts) < 9 {
				continue
			}
			rgName := parts[4]
			serverName := *srv.Name

			// 获取该服务器下的所有数据库
			dbClient := clientFactory.NewDatabasesClient()
			dbPager := dbClient.NewListByServerPager(rgName, serverName, nil)

			for dbPager.More() {
				dbPage, err := dbPager.NextPage(context.Background())
				if err != nil {
					return nil, fmt.Errorf("列举数据库失败 (%s/%s): %v", rgName, serverName, err)
				}
				for _, db := range dbPage.Value {
					owner := ""
					if db.Tags != nil && db.Tags["owner"] != nil {
						owner = *db.Tags["owner"]
					}

					status := ""
					if db.Properties != nil && db.Properties.Status != nil {
						status = string(*db.Properties.Status)
					}

					sqlDatabases = append(sqlDatabases, DBResource{
						Name:     *db.Name,
						ID:       strings.Split(*db.ID, "/")[2],
						Location: *db.Location,
						Owner:    owner,
						Server:   serverName,
						DBType:   "SQL Database",
						Status:   status,
						Tags:     convertTags(db.Tags),
					})
				}
			}
		}
	}

	return sqlDatabases, nil
}

// GetMySQLFlexibleServers 获取Azure MySQL灵活服务器资源列表
func (a *AzureHelper) GetMySQLFlexibleServers() ([]DBResource, error) {
	if a.clientSecretCredential == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	// 创建MySQL灵活服务器客户端工厂
	clientFactory, err := armmysqlflexibleservers.NewClientFactory(
		a.subscriptionID,
		a.clientSecretCredential,
		&arm.ClientOptions{ClientOptions: azcore.ClientOptions{Cloud: cloud.AzureChina}},
	)
	if err != nil {
		return nil, fmt.Errorf("创建MySQL灵活服务器客户端工厂失败: %v", err)
	}

	// 获取MySQL灵活服务器列表
	serversClient := clientFactory.NewServersClient()
	srvPager := serversClient.NewListPager(nil)

	var mysqlServers []DBResource
	for srvPager.More() {
		srvPage, err := srvPager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("列举MySQL灵活服务器失败: %v", err)
		}
		for _, srv := range srvPage.Value {
			owner := ""
			if srv.Tags != nil && srv.Tags["owner"] != nil {
				owner = *srv.Tags["owner"]
			}

			version := ""
			if srv.Properties != nil && srv.Properties.Version != nil {
				version = string(*srv.Properties.Version)
			}

			status := ""
			if srv.Properties != nil && srv.Properties.State != nil {
				status = string(*srv.Properties.State)
			}

			mysqlServers = append(mysqlServers, DBResource{
				Name:     *srv.Name,
				ID:       strings.Split(*srv.ID, "/")[2],
				Location: *srv.Location,
				Owner:    owner,
				DBType:   "MySQL Flexible Server",
				Version:  version,
				Status:   status,
				Tags:     convertTags(srv.Tags),
			})
		}
	}

	return mysqlServers, nil
}

// GetSQLServers 获取Azure SQL服务器资源列表
func (a *AzureHelper) GetSQLServers() ([]DBResource, error) {
	if a.clientSecretCredential == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
	}

	// 创建SQL客户端工厂
	clientFactory, err := armsql.NewClientFactory(
		a.subscriptionID,
		a.clientSecretCredential,
		&arm.ClientOptions{ClientOptions: azcore.ClientOptions{Cloud: cloud.AzureChina}},
	)
	if err != nil {
		return nil, fmt.Errorf("创建SQL客户端工厂失败: %v", err)
	}

	// 获取SQL服务器列表
	serversClient := clientFactory.NewServersClient()
	srvPager := serversClient.NewListPager(nil)

	var sqlServers []DBResource
	for srvPager.More() {
		srvPage, err := srvPager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("列举SQL服务器失败: %v", err)
		}
		for _, srv := range srvPage.Value {
			owner := ""
			if srv.Tags != nil && srv.Tags["owner"] != nil {
				owner = *srv.Tags["owner"]
			}

			version := ""
			if srv.Properties != nil && srv.Properties.Version != nil {
				version = *srv.Properties.Version
			}

			sqlServers = append(sqlServers, DBResource{
				Name:     *srv.Name,
				ID:       strings.Split(*srv.ID, "/")[2],
				Location: *srv.Location,
				Owner:    owner,
				DBType:   "SQL Server",
				Version:  version,
				Status:   "Running",
				Tags:     convertTags(srv.Tags),
			})
		}
	}

	return sqlServers, nil
}

// 辅助函数：转换标签
func convertTags(tags map[string]*string) map[string]string {
	result := make(map[string]string)
	for k, v := range tags {
		if v != nil {
			result[k] = *v
		}
	}
	return result
}

// 为了向后兼容，添加全局函数
func GetAzureSQLDatabases() ([]DBResource, error) {
	azHelper := NewAzureHelper()
	return azHelper.GetSQLDatabases()
}

func GetAzureMySQLFlexibleServers() ([]DBResource, error) {
	azHelper := NewAzureHelper()
	return azHelper.GetMySQLFlexibleServers()
}

func GetAzureSQLServers() ([]DBResource, error) {
	azHelper := NewAzureHelper()
	return azHelper.GetSQLServers()
}

// GetDatabases 获取Azure数据库资源列表（向后兼容方法）
func (a *AzureHelper) GetDatabases() ([]DBResource, error) {
	sqlDatabases, err := a.GetSQLDatabases()
	if err != nil {
		return nil, err
	}

	mysqlServers, err := a.GetMySQLFlexibleServers()
	if err != nil {
		return nil, err
	}

	// 合并结果
	return append(sqlDatabases, mysqlServers...), nil
}

// GetAzureResources 获取Azure资源列表（向后兼容函数）
func GetAzureResources() ([]Resource, error) {
	azHelper := NewAzureHelper()
	return azHelper.GetResources()
}

// GetToken 获取Azure访问令牌（向后兼容函数）
func GetToken() (string, error) {
	azHelper := NewAzureHelper()
	return azHelper.GetToken()
}

// GetAzureVirtualMachines 获取Azure虚拟机资源列表（向后兼容函数）
func GetAzureVirtualMachines() ([]VMResource, error) {
	azHelper := NewAzureHelper()
	return azHelper.GetVirtualMachines()
}

func GetAzureDatabases() ([]DBResource, error) {
	azHelper := NewAzureHelper()
	return azHelper.GetDatabases()
}
