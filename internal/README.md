# internal
Vital part of the web server, following the hexagonal architecture
You must test everything inside `core`, except the vaults. 100% code coverage should be achieved.
Testing other folders will take more efforts to cover the errors, test as much as you can.

1. **core** - composed of `controllers`, `services` and `vaults` (repositories) which should be mock-able for unit testing dependency injection
   1. controllers will take care of service calls, make sure we can call those services and handle errors
   2. services will call vaults
   3. vaults will connect to any database or well-tested imported dependencies (**You do not have to write tests for those packages
      it would be too difficult.**)
2. **handlers** - handle the requests and extract necessary values from the payload. Also, handle responses from the controllers.
3. **injection** - using `wire`, it proved that we can inject dependencies to the server. We can replace them during tests as stubs.
4. **mocks** - just stubs for tests using `mockgen`
   ```yaml
   //go:generate mockgen -source=$GOFILE -destination=../mocks/base_service.mock.go -package=mocks
   ```
   
ps. renaming `repositories` to `vaults` allow me to alphabetically sort them