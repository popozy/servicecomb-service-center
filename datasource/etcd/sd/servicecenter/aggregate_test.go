// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package servicecenter_test

import (
	"github.com/apache/servicecomb-service-center/datasource/etcd"
	_ "github.com/apache/servicecomb-service-center/test"

	"testing"

	. "github.com/apache/servicecomb-service-center/datasource/etcd/sd/servicecenter"
)

func TestNewSCClientAggregate(t *testing.T) {
	etcd.Configuration().ClusterAddresses = "sc-1=127.0.0.1:2379,127.0.0.2:2379"
	etcd.Configuration().InitClusterInfo()
	c := GetOrCreateSCClient()
	if len(*c) == 0 {
		t.Fatalf("TestNewSCClientAggregate failed")
	}
}
