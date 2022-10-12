package util

import "github.com/sirupsen/logrus"

var ApiLogger = logrus.WithField("system", "api")
var FsLogger = logrus.WithField("system", "fs")
var DbLogger = logrus.WithField("system", "db")