#!/bin/bash
go-bindata -pkg assets -prefix data data/... && echo REFRESHED || echo FAILED