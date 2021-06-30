package main

import (
	"database/sql"
)

type HostAPI struct {
	db *sql.DB
}

func (h *HostAPI) CreateHost(ip string, fqdn string, comment string) (int64, error) {
	stmt, err := h.db.Prepare("INSERT INTO hosts VALUES(NULL, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(ip, fqdn, comment)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *HostAPI) UpdateHost(id int64, ip string, fqdn string, comment string) (int64, error) {
	stmt, err := h.db.Prepare("UPDATE hosts SET ip=?, fqdn=?, comment=? WHERE id=?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(ip, fqdn, comment, id)
	if err != nil {
		return 0, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affect, nil
}

func (h *HostAPI) QueryHosts() ([]host, error) {

	rows, err := h.db.Query("SELECT * FROM hosts ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hosts := make([]host, 0)
	var currentHost host

	for rows.Next() {
		err = rows.Scan(&currentHost.ID, &currentHost.IP, &currentHost.FQDN, &currentHost.COMMENT)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, currentHost)
	}
	return hosts, nil
}

func (h *HostAPI) DeleteHost(id int64) (int64, error) {
	stmt, err := h.db.Prepare("DELETE FROM hosts WHERE id=?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affect, nil
}
