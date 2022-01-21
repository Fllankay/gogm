// Copyright (c) 2021 MindStand Technologies, Inc
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package gogm

import (
	"context"
	"log"
	"reflect"

	"github.com/stretchr/testify/require"
)

func testIndexManagement(req *require.Assertions) {
	gogm, err := getTestGogmWithDefaultStructs()
	req.Nil(err)
	req.NotNil(gogm)
	//delete everything
	req.Nil(dropAllIndexesAndConstraints(context.Background(), gogm))

	//setup structure
	_map := toHashmapStructdecconf(map[string]structDecoratorConfig{
		"TEST1": {
			Label:    "Test1",
			IsVertex: true,
			Fields: map[string]decoratorConfig{
				"UUID": {
					Name:       "uuid",
					PrimaryKey: UUIDPrimaryKeyStrategy.StrategyName,
					Type:       reflect.TypeOf(""),
				},
				"IndexField": {
					Name:  "index_field",
					Index: true,
					Type:  reflect.TypeOf(1),
				},
				"UniqueField": {
					Name:   "unique_field",
					Unique: true,
					Type:   reflect.TypeOf(""),
				},
			},
		},
		"TEST2": {
			Label:    "Test2",
			IsVertex: true,
			Fields: map[string]decoratorConfig{
				"UUID": {
					Name:       "uuid",
					PrimaryKey: UUIDPrimaryKeyStrategy.StrategyName,
					Type:       reflect.TypeOf(""),
				},
				"IndexField1": {
					Name:  "index_field1",
					Index: true,
					Type:  reflect.TypeOf(1),
				},
				"UniqueField1": {
					Name:   "unique_field1",
					Unique: true,
					Type:   reflect.TypeOf(""),
				},
			},
		},
	})

	//create stuff
	req.Nil(createAllIndexesAndConstraints(context.Background(), gogm, _map))

	log.Println("created indices and constraints")

	//validate
	req.Nil(verifyAllIndexesAndConstraints(context.Background(), gogm, _map))

	//clean up
	req.Nil(dropAllIndexesAndConstraints(context.Background(), gogm))
}
