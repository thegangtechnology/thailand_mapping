# GEP Backend Template

*Go + echo + PostgreSQL = web server*

In this template, we will implement the reporting system.
You should also be able to host a server through docker compose.

## What to do?
On your first initial commit, use this command
```
make init
```
It will automatically set the module name to 
`your-repository` name.

## Current Spec
```
Go version 1.20
Echo v.4
PostgreSQL:14
```

### Tools
Must use tools to make our life easier.

1. [GVM - go version manager](https://github.com/moovweb/gvm)
2. [Pre-commit - save the commits and CI/CD runs on remote](https://github.com/dnephin/pre-commit-golang)
3. [Make - use the aliases instead of typing the long commands](https://makefiletutorial.com/)
4. [GoLand - all essential tools](https://www.jetbrains.com/go/download/#section=mac) or just use VS code if you want to custom your editors a lot
5. [Docker + AIR - blazing-ly fast hot reload](https://formulae.brew.sh/formula/docker)

#### Optional
6. [OpenAPI - REST building tools](https://www.openapis.org/), but the generated code could be terrible sometimes.
### Directories
If you wonder what each folder does, just read the `README.md` inside.

## Install packages

```sh
go mod tidy
```

## Install linter

Install [golangci-lint](https://golangci-lint.run/) for linting.

Once installed, you should be able these lint commands.

Lint check and automatically fix.
```sh
make lzl
```

Make your life easier by integrate it to [format your code.](https://golangci-lint.run/usage/integrations/)
```
GOLAND TIP!
argument: run --fix
```

## Setup Git Hooks
```
brew install pre-commit

precommit install
```
[Read more about git hooks](https://pre-commit.com/#install)

### Setup environment variables

Create `.env.development` based on `.env.example` for development. For production, create `.env.production` and set `APP_ENV=production`.

##### @dev - Please use docker compose in thegang


## Usage Examples

## Dev (Docker Compose)
In dev stage, you need creds from your PM or Devs and put them to `secrets` before running the following command.
```
docker compose up
```

[Debugging in GoLand](https://blog.jetbrains.com/go/2020/05/06/debugging-a-go-application-inside-a-docker-container/
)

The exposed port should be `32423` and `2345`

## OpenAPI

I suggest writing the `spec.yaml` on separated platform, use them for communication.

## Run

Build and run the project.

```sh
make run
```

## Build

Build the project.

```sh
make compile
```

## Test
Inject the stubs and test by functions. Also, keep the code coverage up to 80%+, you may leave the failed cases on a rush.

Test the project.

```sh
make test
```

### Mocking

```sh
make lzm
```

Using mockgen, we can just run the above command to lazy mock.

You can generated mock client using the following command as an example. Mock only works on interface
```bash
mockgen --source=services/secret_manager_service.go --package=mocked SecretClient > controllers/mocked/secret_service.mocked.go
```

Or use go:generate magic comment
```bash
mockgen -source=$GOFILE -destination=../mocks/your_file.mock.go -package=mocks
```
```bash
go generate ./repositories ./services
```

## Set up the database && Mock data
Only in docker, you can migrate tables
```bash
make mig-up
```

## Reset database (Do not use it on production)
Only in docker, you can remove tables
```bash
make mig-down
```

## Inject dependencies
Define each service init functions and pack it up with the Initialize function.

```bash
make wire
```

## Code Convention
Follow the Controller-Service-Repositories pattern. We can test only just the controllers and services, leaving out some coverages to mock offline.
1. Any loop must be stored in a function (1 for-loop per function)
2. Nested If is cursed code. Please replace it with nested functions
3. Bundling up erroneous functions into 1 function and up to 15 functions can be stored
4. Define structs for any result from function can help reducing the complexity (Export_full_report.go is a good example)
5. Public functions should be at the top, and leave private functions at the bottom.
6. Structs should better be generated from spec.yml
7. If possible, create a struct for any function that require expansion later on. (Server Parameters for example)
8. Proofread your code again and make sure that it is understandable without comments in the code.

### IDE Tips
This could efficiently save your time.
1. If you name your function uniquely, it is quicker for IDE to search the project for usages (CreateSomething, PatchSomething, UpdateSomething).
2. Use debugger (TBC on docker)
3. You can let the lint format automatically by saving your files

## Deployment Guide
Set up your environment as followed, you can either check your docker compose
1. `LOG_LEVEL` - Set to `error`, `warn`, or `debug`.
2. `APP_ENV`- Set to `docker` or `production`

### Deployment Flow

This project uses [Git Flow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow#:~:text=Gitflow%20is%20a%20legacy%20Git,software%20development%20and%20DevOps%20practices.)

- To deploy dev, just push to develop branch
- To deploy QA, merged from develop to a new branch name release/\*, formatted it as release/<semver>
- To deploy UAT and Production, merged from release branch to main branch and create a version tag (v<semver>)

## Set up from dump file
Assume that you have the container,
```zsh
	docker cp path/to/dump.sql <container-name>:/
	docker exec -it <container-id> /bin/sh
	psql -U <username> -d <database> -f dump.sql
```

## FAQ
**Question 1** : Why my docker container cannot find the database?

Check the database port in `.env.docker`, make sure it is matched to the container exposed port.
We found that it switches between `5432` and `5458`.

**Question 2** : I need assistance ASAP!
Please contact Possawat or Narongdej in MS team.