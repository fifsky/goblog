package pagination

// Copyright 2015 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Paginater(t *testing.T) {
	p := New(0, -1, -1, 0)
	assert.Equal(t, p.PagingNum(), 1)
	assert.True(t, p.IsFirst())
	assert.False(t, p.HasPrevious())
	assert.Equal(t, p.Previous(), 1)
	assert.False(t, p.HasNext())
	assert.Equal(t, p.Next(), 1)
	assert.True(t, p.IsLast())
	assert.Equal(t, p.Total(), 0)

	p = New(1, 10, 2, 0)
	assert.Equal(t, p.PagingNum(), 10)
	assert.True(t, p.IsFirst())
	assert.False(t, p.HasPrevious())
	assert.False(t, p.HasNext())
	assert.True(t, p.IsLast())

	p = New(10, 10, 1, 0)
	assert.Equal(t, p.PagingNum(), 10)
	assert.True(t, p.IsFirst())
	assert.False(t, p.HasPrevious())
	assert.False(t, p.HasNext())
	assert.True(t, p.IsLast())

	p = New(11, 10, 1, 0)
	assert.Equal(t, p.PagingNum(), 10)
	assert.True(t, p.IsFirst())
	assert.False(t, p.HasPrevious())
	assert.True(t, p.HasNext())
	assert.Equal(t, p.Next(), 2)
	assert.False(t, p.IsLast())

	p = New(11, 10, 2, 0)
	assert.Equal(t, p.PagingNum(), 10)
	assert.False(t, p.IsFirst())
	assert.True(t, p.HasPrevious())
	assert.Equal(t, p.Previous(), 1)
	assert.False(t, p.HasNext())
	assert.True(t, p.IsLast())

	p = New(20, 10, 2, 0)
	assert.Equal(t, p.PagingNum(), 10)
	assert.False(t, p.IsFirst())
	assert.True(t, p.HasPrevious())
	assert.False(t, p.HasNext())
	assert.True(t, p.IsLast())

	p = New(25, 10, 2, 0)
	assert.Equal(t, p.PagingNum(), 10)
	assert.False(t, p.IsFirst(), )
	assert.True(t, p.HasPrevious())
	assert.True(t, p.HasNext())
	assert.False(t, p.IsLast())
}