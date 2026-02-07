package user_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource_basic(t *testing.T) {
	userID := os.Getenv("DISCORD_USER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			if os.Getenv("DISCORD_USER_ID") == "" {
				t.Fatal("DISCORD_USER_ID must be set for acceptance tests")
			}
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.discord_user.test", "id", userID),
					resource.TestCheckResourceAttrSet("data.discord_user.test", "username"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(userID string) string {
	return fmt.Sprintf(`
data "discord_user" "test" {
  id = %[1]q
}
`, userID)
}
