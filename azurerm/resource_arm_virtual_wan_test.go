package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualWan_basic(t *testing.T) {
	resourceName := "azurerm_virtual_wan.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualWan_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_virtual_wan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualWan_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_wan"),
			},
		},
	})
}

func TestAccAzureRMVirtualWan_complete(t *testing.T) {
	resourceName := "azurerm_virtual_wan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_complete(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMVirtualWanDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualWanClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_wan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual WAN still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMVirtualWanExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualWanName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Virtual WAN: %s", virtualWanName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualWanClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualWanName)
		if err != nil {
			return fmt.Errorf("Bad: Get on virtualWanClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual WAN %q (resource group: %q) does not exist", virtualWanName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMVirtualWan_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, rInt, location, rInt)
}

func testAccAzureRMVirtualWan_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualWan_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_wan" "import" {
  name                = azurerm_virtual_wan.test.name
  resource_group_name = azurerm_virtual_wan.test.resource_group_name
  location            = azurerm_virtual_wan.test.location
}
`, template)
}

func testAccAzureRMVirtualWan_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  disable_vpn_encryption            = false
  allow_branch_to_branch_traffic    = true
  allow_vnet_to_vnet_traffic        = true
  office365_local_breakout_category = "All"

  tags = {
    Hello = "There"
    World = "Example"
  }
}
`, rInt, location, rInt)
}
