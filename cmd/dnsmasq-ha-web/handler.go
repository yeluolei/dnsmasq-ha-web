package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

//----------
// Handlers
//----------

func createHost(c echo.Context, dbAPI *HostAPI) error {
	h := new(host)
	if err := c.Bind(h); err != nil {
		return err
	}
	id, err := dbAPI.CreateHost(h.IP, h.FQDN, h.COMMENT)
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
	_, err := dbAPI.UpdateHost(int64(id), h.IP, h.FQDN, h.COMMENT)
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

	hostsMap := map[string][]string{}
	for _, host := range hosts {
		if currentLine, ok := hostsMap[host.IP]; ok {
			hostsMap[host.IP] = append(currentLine, host.FQDN)
		} else {
			hostsMap[host.IP] = []string{host.FQDN}
		}
	}

	hostLines := []string{}

	for ip, fqdns := range hostsMap {
		hostLines = append(hostLines, fmt.Sprintf("%-16s %s\n", ip, strings.Join(fqdns, " ")))
	}

	content := strings.Join(hostLines, "")

	err = ioutil.WriteFile(file.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	os.Rename(file.Name(), hostFilePath)
	os.Chmod(hostFilePath, 0666)

	return c.NoContent(http.StatusOK)
}
