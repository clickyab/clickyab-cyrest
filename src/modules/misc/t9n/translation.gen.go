package t9n

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateTranslation try to save a new Translation in database
func (m *Manager) CreateTranslation(t *Translation) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(t)

	return m.GetDbMap().Insert(t)
}

// UpdateTranslation try to update Translation in database
func (m *Manager) UpdateTranslation(t *Translation) error {
	t.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(t)

	_, err := m.GetDbMap().Update(t)
	return err
}

// ListTranslations try to list all Translations without pagination
func (m *Manager) ListTranslationsWithFilter(filter string, params ...interface{}) []Translation {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []Translation
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", TranslationTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListTranslations try to list all Translations without pagination
func (m *Manager) ListTranslations() []Translation {
	return m.ListTranslationsWithFilter("")
}

// CountTranslations count entity in Translations table with valid where filter
func (m *Manager) CountTranslationsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", TranslationTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountTranslations count entity in Translations table
func (m *Manager) CountTranslations() int64 {
	return m.CountTranslationsWithFilter("")
}

// ListTranslationsWithPaginationFilter try to list all Translations with pagination and filter
func (m *Manager) ListTranslationsWithPaginationFilter(offset, perPage int, filter string, params ...interface{}) []Translation {
	var res []Translation
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(params)+1, len(params)+2)
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", TranslationTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListTranslationsWithPagination try to list all Translations with pagination
func (m *Manager) ListTranslationsWithPagination(offset, perPage int) []Translation {
	return m.ListTranslationsWithPaginationFilter(offset, perPage, "")
}
