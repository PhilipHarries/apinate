package main

import (
	"fmt"
	cfg "github.com/PhilipHarries/apinate/config"
	"github.com/PhilipHarries/apinate/exec"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	configfile := ""
	configLocations := []string{
		filepath.Join(os.Getenv("HOME"), ".apinate.toml"),
		filepath.Join(os.Getenv("HOME"), ".apinate.json"),
		filepath.Join(os.Getenv("HOME"), ".apinate.yaml"),
		"/etc/apinate.toml",
		"/etc/apinate.json",
		"/etc/apinate.yaml",
	}
	fmt.Println(configLocations)
	for _, configLocation := range configLocations {
		if _, err := os.Stat(configLocation); err == nil {
			configfile = configLocation
			break
		}
	}
	if configfile == "" {
		fmt.Println("No configfile found")
		os.Exit(1)
	}

	var config cfg.Config
	var err error

	fmt.Println(configfile)
	fmt.Println(filepath.Ext(configfile))
	switch {
	case filepath.Ext(configfile) == ".toml":
		config, err = cfg.LoadConfigTOML(configfile)
	case filepath.Ext(configfile) == ".json":
		config, err = cfg.LoadConfigJSON(configfile)
	case filepath.Ext(configfile) == ".yaml":
		config, err = cfg.LoadConfigYAML(configfile)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router := gin.Default()
	if config.ContentType == "html" {
		router.LoadHTMLGlob("templates/plain.tmpl")
	}

	for _, mapping := range config.Mappings {
		fmt.Println(fmt.Sprintf("mapping %s to command %s", mapping.Resource, mapping.Command))
		res := mapping.Resource
		cmd := mapping.Command
		params := mapping.Params
		var command string
		if params {
			joinedParams := []string{res, "/:params"}
			fmt.Println(joinedParams)
			res = strings.Join(joinedParams, "")
		}
		router.GET(res, func(c *gin.Context) {
			if params {
				fmt.Println(fmt.Sprintf("cmd: %s", cmd))
				fmt.Println(fmt.Sprintf("params: %s", c.Param("params")))
				joinedCmd := []string{cmd, c.Param("params")}
				fmt.Println(fmt.Sprintf("joinedCmd: %s", joinedCmd))
				command = strings.Join(joinedCmd, " ")
			} else {
				command = cmd
			}
			msg, err := exec.Exec(command)
			if err != nil {
				fmt.Println(err)
			}
			if config.ContentType == "json" {
				c.JSON(200, gin.H{"message": msg})
			} else if config.ContentType == "xml" {
				c.XML(200, gin.H{"message": msg})
			} else if config.ContentType == "yaml" {
				c.YAML(200, gin.H{"message": msg})
			} else if config.ContentType == "html" {
				c.HTML(http.StatusOK, "plain.tmpl",
					gin.H{"content": msg})
			} else if config.ContentType == "txt" {
				c.String(http.StatusOK, "%s", strings.Join(msg, "\n"))
			}
		})
	}

	router.Run(":8080")

}
