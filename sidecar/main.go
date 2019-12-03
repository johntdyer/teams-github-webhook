package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type outputEvent struct {
	Body        string            `json:"Body"`
	Headers     map[string]string `json:"Headers"`
	Method      string            `json:"Method"`
	QueryParams string            `json:"QueryParams"`
}

func executeWHRelay(cfg relayConfig) {
	if cfg.srcProto != "AWS_SQS" {
		log.Fatal("AWS SQS is only supported...")
	}
	sess := setupReceiveSession(cfg)
	timeout := int64(10) // 10 seconds
	go func() {
		for {
			msg, enc := receiveMsgOnSession(sess, timeout)
			if msg == nil {
                            continue
                        }
                        sDec := []byte("")
                        if (enc == "base64") {
                            sDec, _ = b64.StdEncoding.DecodeString(*msg)
                        }
                        fmt.Println("Input Event:", string(sDec))
                        var evt outputEvent
                        err := json.Unmarshal(sDec, &evt)
                        if err != nil {
                            log.Println("Error unmarshalling...")
                            continue
                        }
                        repoName, ok1 := evt.Headers["repoName"]
                        prNum, ok := evt.Headers["pullRequestId"]
                        fmt.Println("repoName: ", repoName, "prNum:", prNum)
                        if ok && ok1 {
                                putUrl := cfg.dstUrl + repoName + "/pulls/" + prNum + "/merge"
                                dispatch("PUT", putUrl, cfg.insideAuth, "")
                        } else {
                            log.Println("Did NOT find repoName and/or pullRequestId")
                        }
		}
	}()
}

func main() {
	cfg := parseCmdLineArgs()
	fmt.Println("cfg:", cfg)
	executeWHRelay(cfg)
	time.Sleep(100000 * time.Second)
}
