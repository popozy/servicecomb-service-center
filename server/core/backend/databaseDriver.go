/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package backend

import (
	"context"
	"errors"
	"github.com/apache/servicecomb-service-center/pkg/backoff"
	"github.com/apache/servicecomb-service-center/pkg/gopool"
	"github.com/apache/servicecomb-service-center/pkg/log"
	"github.com/apache/servicecomb-service-center/server/plugin"
	"github.com/apache/servicecomb-service-center/server/plugin/databaseDriver"
	"sync"
	"time"
)

var (
	driverInstance *DatabaseDriver
	driverSingletonLock  sync.Mutex
)
func NewDataBaseDriver() (*DatabaseDriver, error) {
	instance := plugin.Plugins().DatabaseDriver()
	if instance == nil {
		return nil, errors.New("register center client plugin does not exist")
	}
	select {
	case err := <-instance.Err():
		plugin.Plugins().Reload(plugin.REGISTRY)
		return nil, err
	case <-instance.Ready():
	}
	return &DatabaseDriver{
		DatabaseDriver:  instance,
		goroutine: gopool.New(context.Background()),
	}, nil
}

func Database() databaseDriver.DatabaseDriver {
	return GetDataBaseDriver()
}

type DatabaseDriver struct {
	databaseDriver.DatabaseDriver
	goroutine *gopool.Pool
}

func GetDataBaseDriver() *DatabaseDriver {
	if engineInstance == nil {
		driverSingletonLock.Lock()
		for i := 0; driverInstance == nil; i++ {
			inst, err := NewDataBaseDriver()
			if err != nil {
				log.Errorf(err, "get register center client failed")
			}
			driverInstance = inst

			if engineInstance != nil {
				driverSingletonLock.Unlock()
				return driverInstance
			}

			t := backoff.GetBackoff().Delay(i)
			log.Errorf(nil, "initialize service center failed, retry after %s", t)
			<-time.After(t)
		}
		driverSingletonLock.Unlock()
	}
	return driverInstance
}