package main

import (
	cfg "github.com/PhilipHarries/apinate/config"
	"github.com/PhilipHarries/apinate/exec"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	// log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stderr) // or a file
	log.SetLevel(log.DebugLevel)

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
		log.WithFields(log.Fields{
			"critical": true,
		}).Fatal("No config file could be found")
	}

	var config cfg.Config
	var err error
	var logerr error
	var lf *os.File

	config, err = cfg.LoadConfig(configfile)
	if err != nil {
		log.WithFields(log.Fields{
			"critical":   true,
			"configfile": configfile,
		}).Fatal("Configfile could not be loaded")
	}

	if config.Logfile != "stdout" {
		if _, logerr = os.Stat(config.Logfile); logerr != nil {
			lf, logerr = os.Create(config.Logfile)
		} else {
			lf, logerr = os.OpenFile(config.Logfile, os.O_APPEND|os.O_WRONLY, 0644)
		}
		if logerr != nil {
			log.WithFields(log.Fields{
				"critical": true,
				"logfile":  config.Logfile,
				"error":    logerr,
			}).Fatal("Logfile could not be opened")
		}
		log.SetOutput(lf)
	}
	gin.DefaultWriter = lf
	defer lf.Close()

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
		log.WithFields(log.Fields{
			"resource": mapping.Resource,
			"command":  mapping.Command,
			"params":   mapping.Params,
		}).Info("Mapping defined")
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
				joinedCmd := []string{cmd, c.Param("params")}
				command = strings.Join(joinedCmd, " ")
			} else {
				command = cmd
			}
			var msg []string
			var err error
			if msg, err = exec.Exec(command); err != nil {
				log.WithFields(log.Fields{
					"command": command,
					"message": msg,
					"error":   err,
				}).Info("Error running command.")
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
			} else if config.ContentType == "raw" {
				c.String(http.StatusOK, "%s", strings.Join(msg, "\n"))
			}
		})
	}

	listenString := config.Address + ":" + strconv.Itoa(config.Port)
	router.Run(listenString)

}
