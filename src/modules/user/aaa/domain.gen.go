package aaa

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateDomain try to save a new Domain in database
func (m *Manager) CreateDomain(d *Domain) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(d)

	return m.GetDbMap().Insert(d)
}

// UpdateDomain try to update Domain in database
func (m *Manager) UpdateDomain(d *Domain) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(d)

	_, err := m.GetDbMap().Update(d)
	return err
}

// ListDomainsWithFilter try to list all Domains without pagination
func (m *Manager) ListDomainsWithFilter(filter string, params ...interface{}) []Domain {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []Domain
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", DomainTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListDomains try to list all Domains without pagination
func (m *Manager) ListDomains() []Domain {
	return m.ListDomainsWithFilter("")
}

// CountDomainsWithFilter count entity in Domains table with valid where filter
func (m *Manager) CountDomainsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", DomainTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountDomains count entity in Domains table
func (m *Manager) CountDomains() int64 {
	return m.CountDomainsWithFilter("")
}

// ListDomainsWithPaginationFilter try to list all Domains with pagination and filter
func (m *Manager) ListDomainsWithPaginationFilter(
	offset, perPage int, filter string, params ...interface{}) []Domain {
	var res []Domain
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", DomainTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListDomainsWithPagination try to list all Domains with pagination
func (m *Manager) ListDomainsWithPagination(offset, perPage int) []Domain {
	return m.ListDomainsWithPaginationFilter(offset, perPage, "")
}

// FindDomainByID return the Domain base on its id
func (m *Manager) FindDomainByID(id int64) (*Domain, error) {
	var res Domain
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=?", DomainTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindDomainByCName return the Domain base on its cname
func (m *Manager) FindDomainByCName(c string) (*Domain, error) {
	var res Domain
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE cname=?", DomainTableFull),
		c,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PreInsert is gorp hook to prevent Insert without transaction
func (d *Domain) PreInsert(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Insert Domain must be in transaction")
	}
	return nil
}
