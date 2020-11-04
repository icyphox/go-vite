package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func viteNew(path string) {
	_, file := filepath.Split(path)
	url := strings.TrimSuffix(file, filepath.Ext(file))

	content := fmt.Sprintf(`---
template:
url: %s
title:
subtitle:
date: %s
---`, url, time.Now().Format("2006-01-02"))

	_, err := os.Create(path)
	if err != nil {
		printErr(err)
		return
	}
	ioutil.WriteFile(path, []byte(content), 0644)
	printMsg("created:", path)
}
