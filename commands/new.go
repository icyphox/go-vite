package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func New(path string) error {
	_, file := filepath.Split(path)
	url := strings.TrimSuffix(file, filepath.Ext(file))

	content := fmt.Sprintf(`---
template:
slug: %s
title:
subtitle:
date: %s
---`, url, time.Now().Format("2006-01-02"))

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
		os.WriteFile(path, []byte(content), 0755)
		fmt.Printf("vite: created new post at %s\n", path)
		return nil
	}

	fmt.Printf("error: %s already exists\n", path)
	return nil
}
