#!/bin/sh

mockgen -destination mocks/busobject.go -imports github.com/godbus/dbus -package mocks github.com/godbus/dbus BusObject

