package provider

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ignoreCase(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func keyFromStrings(s []string) string {
	ss := make([]string, len(s))
	copy(ss, s)
	sort.Strings(ss)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(ss, "|"))))
}

func protectedFromDelete(resourceData *schema.ResourceData) bool {
	return resourceData.Get("deletion_protection").(bool)
}

// checkRequiredTestVars returns true if the given environment variables are not empty
func checkRequiredTestVars(t *testing.T, vars ...string) {
	for _, v := range vars {
		if strings.TrimSpace(os.Getenv(v)) == "" {
			t.Skipf("skipping test due to missing %s environment variable", v)
		}
	}
}
