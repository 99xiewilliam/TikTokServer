// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
	"log"
)

var DB *gorm.DB

// Init init DB
func Init() {
	var err error

	DB, err = gorm.Open(mysql.Open("engineer6:123456@tcp(localhost:9910)/TikTokDB?charset=utf8&parseTime=True&loc=Local"),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}

	log.Printf("DB initalizing...")
	m := DB.Migrator()
	if m.HasTable(&User{}) {
		log.Printf("User table already exists")
		return
	}
	if err = m.CreateTable(&User{}); err != nil {
		panic(err)
	}
	log.Printf("Created user table already exists")
}