package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/txn2/txeh"
)

//----------
// Handlers
//----------

func createHost(c echo.Context, dbAPI *HostAPI) error {
	h := new(host)
	if err := c.Bind(h); err != nil {
		return err
	}
	id, err := dbAPI.CreateHost(h.IP, h.FQDN)
	if err != nil {
		return err
	}
	h.ID = int(id)
	return c.JSON(http.StatusCreated, h)
}

func updateHost(c echo.Context, dbAPI *HostAPI) error {
	h := new(host)
	if err := c.Bind(h); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	h.ID = id
	_, err := dbAPI.UpdateHost(int64(id), h.IP, h.FQDN)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, h)
}

func deleteHost(c echo.Context, dbAPI *HostAPI) error {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := dbAPI.DeleteHost(int64(id))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func getAllHosts(c echo.Context, dbAPI *HostAPI) error {
	hosts, err := dbAPI.QueryHosts()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, hosts)
}

func generateHosts(c echo.Context, dbAPI *HostAPI, hostFilePath string) error {
	hosts, err := dbAPI.QueryHosts()
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile("", "dnsmasq.hosts")
	if err != nil {
		log.Fatal(err)
	}

	hostsFile, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath: file.Name(),
	})

	if err != nil {
		return err
	}

	for _, host := range hosts {
		hostsFile.AddHost(host.IP, host.FQDN)
	}

	if err := hostsFile.Save(); err != nil {
		return err
	}

	os.Rename(file.Name(), hostFilePath)

	return c.NoContent(http.StatusOK)
}
