package main

import (
	"github.com/TheTeaParty/notnotes-platform/internal"
	"github.com/TheTeaParty/notnotes-platform/internal/pkg/logger"
	"go.uber.org/zap"
)

func main() {

	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	a, err := internal.InitializeApplication()
	if err != nil {
		l.With(zap.Error(err)).Fatal("Error init application")
	}

	if err := a.RunHTTP(); err != nil {
		l.With(zap.Error(err)).Fatal("Error running http server")
	}
}
