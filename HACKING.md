# Hacking on WConn

## Dependencies

Dependencies are handled by go modules. 

### Update and upgrade

You can update dependencies by
```sh
go get -u -t ./...
```

### Vendor

After updating dependencies, you can prepare the packaging of
the project by including all dependencies in a vendor subfolder with:
```sh
go mod vendor
```

## Testing

### GoMock

This framework is used in order to mock some external systems
like dbus. You have to generate mocks by using *mockgen*. First
you should install it:

```sh
sudo apt install -y mockgen
```
