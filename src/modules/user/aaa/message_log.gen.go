package aaa

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateMessageLog try to save a new MessageLog in database
func (m *Manager) CreateMessageLog(ml *MessageLog) error {
	now := time.Now()
	ml.CreatedAt = now
	ml.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ml)

	return m.GetDbMap().Insert(ml)
}

// UpdateMessageLog try to update MessageLog in database
func (m *Manager) UpdateMessageLog(ml *MessageLog) error {
	ml.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ml)

	_, err := m.GetDbMap().Update(ml)
	return err
}

// ListMessageLogs try to list all MessageLogs without pagination
func (m *Manager) ListMessageLogsWithFilter(filter string, params ...interface{}) []MessageLog {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []MessageLog
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", MessageLogTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListMessageLogs try to list all MessageLogs without pagination
func (m *Manager) ListMessageLogs() []MessageLog {
	return m.ListMessageLogsWithFilter("")
}

// CountMessageLogs count entity in MessageLogs table with valid where filter
func (m *Manager) CountMessageLogsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", MessageLogTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountMessageLogs count entity in MessageLogs table
func (m *Manager) CountMessageLogs() int64 {
	return m.CountMessageLogsWithFilter("")
}

// ListMessageLogsWithPaginationFilter try to list all MessageLogs with pagination and filter
func (m *Manager) ListMessageLogsWithPaginationFilter(offset, perPage int, filter string, params ...interface{}) []MessageLog {
	var res []MessageLog
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(params)+1, len(params)+2)
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", MessageLogTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListMessageLogsWithPagination try to list all MessageLogs with pagination
func (m *Manager) ListMessageLogsWithPagination(offset, perPage int) []MessageLog {
	return m.ListMessageLogsWithPaginationFilter(offset, perPage, "")
}

// FindMessageLogByID return the MessageLog base on its id
func (m *Manager) FindMessageLogByID(id int64) (*MessageLog, error) {
	var res MessageLog
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", MessageLogTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetUserMessageLogs return all MessageLogs belong to User
func (m *Manager) GetUserMessageLogs(u *User) []MessageLog {
	var res []MessageLog
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s WHERE user_id=$1",
			MessageLogTableFull,
		),
		u.ID,
	)

	assert.Nil(err)
	return res
}

// CountUserMessageLogs return count MessageLogs belong to User
func (m *Manager) CountUserMessageLogs(u *User) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE user_id=$1",
			MessageLogTableFull,
		),
		u.ID,
	)

	assert.Nil(err)
	return res
}
