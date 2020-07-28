package chttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortRoutes_PartsLen(t *testing.T) {
	input := []Route{
		{Path: "/foo"},
		{Path: "/foo/bar"},
	}
	expected := []Route{
		{Path: "/foo/bar"},
		{Path: "/foo"},
	}

	sortRoutes(input)

	assert.Equal(t, expected, input)
}

func TestSortRoutes_Regex(t *testing.T) {
	input := []Route{
		{Path: "/foo/{uuid}"},
		{Path: "/foo/bar"},
	}
	expected := []Route{
		{Path: "/foo/bar"},
		{Path: "/foo/{uuid}"},
	}

	sortRoutes(input)

	assert.Equal(t, expected, input)
}

func TestSortRoutes_Index(t *testing.T) {
	input := []Route{
		{Path: "/"},
		{Path: "/foo"},
	}
	expected := []Route{
		{Path: "/foo"},
		{Path: "/"},
	}

	sortRoutes(input)

	assert.Equal(t, expected, input)
}
