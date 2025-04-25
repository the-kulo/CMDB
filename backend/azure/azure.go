package azure

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// Resource 表示Azure资源的结构
type Resource struct {
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Location string            `json:"location"`
	Owner    string            `json:"owner"`
	Type     string            `json:"type"`
	Tags     map[string]string `json:"tags"`
}

// VMResource 表示Azure虚拟机资源
type VMResource struct {
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Location string            `json:"location"`
	Owner    string            `json:"owner"`
	Size     string            `json:"size,omitempty"`
	Status   string            `json:"status,omitempty"`
	Tags     map[string]string `json:"tags"`
}

// DBResource 表示Azure数据库资源
type DBResource struct {
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Location string            `json:"location"`
	Owner    string            `json:"owner"`
	Server   string            `json:"server,omitempty"`
	Tags     map[string]string `json:"tags"`
}

// AzureHelper 封装Azure认证和资源获取功能
type AzureHelper struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	graphClient           *msgraphsdk.GraphServiceClient
	subscriptionID        string
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

			size := ""
			if vm.Properties != nil && vm.Properties.HardwareProfile != nil && vm.Properties.HardwareProfile.VMSize != nil {
				size = string(*vm.Properties.HardwareProfile.VMSize)
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
				Size:     size,
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

// 为了向后兼容，保留原有的函数

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
