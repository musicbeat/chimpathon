// Copyright 2015 Blackhawk Network, Inc.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chimpmail

import (
	"fmt"
	"testing"
)

func TestGetTemplateInfo(t *testing.T) {
	fmt.Println("Test: GetTemplateInfo")
	err := GetTemplateInfo()
	fmt.Printf("err: %#v\n", err)
}
