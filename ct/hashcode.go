package ct

import (
	"hash/crc32"
	"strconv"
)

// Migrated from the V1 SDK github.com/hashicorp/terraform-plugin-sdk/helper/hashcode
// https://www.terraform.io/plugin/sdkv2/guides/v2-upgrade-guide#removal-of-helper-hashcode-package

// Hashcode hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func hashcode(s string) string {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return strconv.Itoa(v)
	}
	if -v >= 0 {
		return strconv.Itoa(-v)
	}
	// v == MinInt
	return "0"
}
