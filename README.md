# ynab-exporter

A Prometheus exporter for the You Need a Budget (YNAB) budgeting software, written in Go.

## Project Status

This project is very much so still in active development, and should be expected to change.

## Exported Timeseries

| Timeseries Name                  | Description                                   | Labels                                             |
|----------------------------------|-----------------------------------------------|----------------------------------------------------|
| `ynab_account_cleared_balance`   | Cleared balance of account                    | `budget_id`, `budget_name`, `account_name`, `type` |
| `ynab_account_uncleared_balance` | Uncleared balance of account                  | `budget_id`, `budget_name`, `account_name`, `type` |
| `ynab_category_budgeted`         | Amount budgeted to category                   | `category_group_name`, `category_name`             |
| `ynab_category_activity`         | Amount of activity in category                | `category_group_name`, `category_name`             |
| `ynab_category_balance`          | Category balance                              | `category_group_name`, `category_name`             |
| `ynab_exporter_api_calls_count`  | Count of calls to the YNAB API                |                                                    |
| `ynab_exporter_build_info`       | Build info for this instance of ynab-exporter | `build_version`, `build_time`                      |

## Running the Exporter

### Configuration

The following environment variables will be read by the exporter:
| Variable Name  | Required? | Default Value | Possible Values                           | Description                                  |
|----------------|-----------|---------------|-------------------------------------------|----------------------------------------------|
| `YNAB_API_TOKEN` | Yes       | None          | Any valid API token                       | Sets the YNAB API token used by the exporter |
| `BIND_ADDR`      | No        | `0.0.0.0`     | Any valid IP address or localhost         | Sets the bind address of the exporter        |
| `PORT`           | No        | `9090`        | Any valid port                            | Sets the port that the exporter binds to     |
| `LOG_LEVEL`      | No        | `INFO`        | `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL` | Sets the log level of the exporter           |

### Running Locally

You can run the exporter in the terminal as follows:
```bash
# Without an explicit build
$ export YNAB_API_TOKEN=${token}
$ go run ./cmd/ynab-exporter

# With an explicit build
$ export YNAB_API_TOKEN=${token}
$ go build -o /tmp/ynab-exporter -ldflags "-X 'github.com/mcbobke/ynab-exporter/cmd/ynab-exporter/version.BuildTime=$(date +%s)' -X 'github.com/mcbobke/ynab-exporter/cmd/ynab-exporter/version.BuildVersion=local'" ./cmd/ynab-exporter
$ /tmp/ynab-exporter
```

### Running Locally in Docker

You can run the exporter in a Docker container as follows:
```bash
$ export YNAB_API_TOKEN=${token}
$ docker build -t localhost/ynab-exporter:latest --no-cache --build-arg BUILD_TIME=$(date +%s) --build-arg BUILD_VERSION=local .
$ docker run -ite YNAB_API_TOKEN=${YNAB_API_TOKEN} --name ynab-exporter --rm --publish 9090:9090/tcp localhost/ynab-exporter:latest
```

## Deploying the Exporter

The exporter can be deployed in a few different ways:
* Deployed to a baremetal or virtual machine and set up as a systemd service (or whatever other init system you use)
* Deployed via docker-compose to a standalone Docker machine
* Deployed to a Kubernetes cluster

See the [examples](./examples) directory.

## TODO

* Proper handling of API token secret (offer alternatives to envvar)
* Testing
