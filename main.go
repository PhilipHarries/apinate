package main

import (
	"fmt"
	cfg "github.com/PhilipHarries/apinate/config"
	"github.com/PhilipHarries/apinate/exec"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	configfile := ""
	thisConfigFile := ""
	configLocations := []string{
		filepath.Join(os.Getenv("HOME"), ".apinate"),
		"/etc/apinate",
	}
	configFiles := []string{
		"apinate.toml",
		"apinate.json",
		"apinate.yaml",
	}
	for _, configLocation := range configLocations {
		for _, configFile := range configFiles {
			thisConfigFile = filepath.Join(configLocation, configFile)
			if _, err := os.Stat(thisConfigFile); err == nil {
				configfile = thisConfigFile
				break
			}
		}
	}
	if configfile == "" {
		fmt.Println("No configfile found")
		os.Exit(1)
	}

	var config cfg.Config
	var err error

	config, err = cfg.LoadConfig(configfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	if config.ContentType == "html" {
		templateLocations := []string{
			filepath.Join(os.Getenv("HOME"), ".apinate/templates"),
			"/usr/share/apinate/templates",
		}
		for _, templateLocation := range templateLocations {
			if _, err := os.Stat(templateLocation); err == nil {
				router.LoadHTMLGlob(filepath.Join(templateLocation, "*.tmpl"))
				break
			}
		}
	}

	for _, mapping := range config.Mappings {
		fmt.Println(fmt.Sprintf("mapping %s to command %s", mapping.Resource, mapping.Command))
		res := mapping.Resource
		cmd := mapping.Command
		params := mapping.Params
		template := mapping.Template
		var command string
		if params {
			joinedParams := []string{res, "/:params"}
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
			var msg []string
			var err error
			if msg, err = exec.Exec(command); err != nil {
				fmt.Println(err)
			}
			if config.ContentType == "json" {
				c.JSON(200, gin.H{"message": msg})
			} else if config.ContentType == "xml" {
				c.XML(200, gin.H{"message": msg})
			} else if config.ContentType == "yaml" {
				c.YAML(200, gin.H{"message": msg})
			} else if config.ContentType == "html" {
				if template == "" {
					template = "plain.tmpl"
				}
				c.HTML(http.StatusOK, template,
					gin.H{"content": msg})
			} else if config.ContentType == "txt" {
				c.String(http.StatusOK, "%s", strings.Join(msg, "\n"))
			}
		})
	}

	listenString := config.Address + ":" + strconv.Itoa(config.Port)
	router.Run(listenString)

}
