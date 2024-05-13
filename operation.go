package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	oas "github.com/getkin/kin-openapi/openapi3"
)

type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

type operation struct {
	path    string
	method  Method
	summary string
	tags    []string
}

func (e operation) Title() string {
	return fmt.Sprintf("%s %s", e.method, e.path)
}

func (e operation) Description() string {
	return fmt.Sprintf("%s %s", e.summary, strings.Join(e.tags, ","))
}
func (e operation) FilterValue() string { return e.path }

func NewOperation(path string, method Method, operationItem *oas.Operation) operation {
	return operation{
		path:    path,
		summary: operationItem.Summary,
		method:  method,
		tags:    operationItem.Tags,
	}
}

func NewOperationsList(doc *oas.T) list.Model {
	operations := make([]list.Item, 0)
	for path, pathItem := range doc.Paths.Map() {
		if operationItem := pathItem.Get; operationItem != nil {
			operations = append(operations, NewOperation(path, Get, operationItem))
		}

		if operationItem := pathItem.Post; operationItem != nil {
			operations = append(operations, NewOperation(path, Post, operationItem))
		}

		if operationItem := pathItem.Put; operationItem != nil {
			operations = append(operations, NewOperation(path, Put, operationItem))
		}

		if operationItem := pathItem.Patch; operationItem != nil {
			operations = append(operations, NewOperation(path, Patch, operationItem))
		}

		if operationItem := pathItem.Delete; operationItem != nil {
			operations = append(operations, NewOperation(path, Delete, operationItem))
		}
	}

	methodOrder := map[Method]int{
		Get:    1,
		Post:   2,
		Put:    3,
		Patch:  4,
		Delete: 5,
	}
	sort.Slice(operations, func(i, j int) bool {
		// sort by operation path first then by method
		if operations[i].(operation).path == operations[j].(operation).path {
			return methodOrder[operations[i].(operation).method] < methodOrder[operations[j].(operation).method]
		}

		return operations[i].(operation).path < operations[j].(operation).path
	})

	operationsList := list.New(operations, list.NewDefaultDelegate(), 0, 0)
	operationsList.Title = "endpoints"
	operationsList.SetShowPagination(false)

	return operationsList
}
