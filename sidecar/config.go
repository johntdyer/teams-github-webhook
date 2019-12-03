package main

import (
	"flag"
	"github.com/go-ini/ini"
	"log"
)

type relayConfig struct {
	awsRegion   string
	awsAccessId string
	awsSecret   string
	awsUrl      string
	srcProto    string
	dstUrl      string
	insideAuth  string
	dstVerb     []string
}

func parseCmdLineArgs() relayConfig {
	config := flag.String("config", "", "specify file which contains config for GLM")
	flag.Parse()

	if *config == "" {
		log.Fatal("Please specify a config file! Exiting..")
	}

	iniCfg, err := ini.Load(*config)
	if iniCfg == nil || err != nil {
		log.Fatal("Cannot load Config file:!!", err)
	}
	globalSection := iniCfg.Section("globals")
	reg := globalSection.Key("aws_region").String()
	id := globalSection.Key("aws_access_id").String()
	sec := globalSection.Key("aws_secret").String()

	hookRelaySection := iniCfg.Section("hook_relay")
	sp := hookRelaySection.Key("source").String()
	awsurl := hookRelaySection.Key("aws_url").String()
	durl := hookRelaySection.Key("dst_url").String()
	vrb := hookRelaySection.Key("allowed_verbs").Strings(",")
	auth := hookRelaySection.Key("inside_auth").String()

	return relayConfig{
		awsRegion:   reg,
		awsAccessId: id,
		awsSecret:   sec,
		awsUrl:      awsurl,
		srcProto:    sp,
		dstUrl:      durl,
		dstVerb:     vrb,
		insideAuth:  auth,
	}
}
