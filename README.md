# ynab-exporter

A Prometheus exporter for the You Need a Budget (YNAB) budgeting software, written in Go.

## Project Status

This project is very much so still in active development, and should be expected to change.

## Exported Timeseries

| Timeseries Name                  | Description                    | Labels                                             |
|----------------------------------|--------------------------------|----------------------------------------------------|
| `ynab_account_cleared_balance`   | Cleared balance of account     | `budget_id`, `budget_name`, `account_name`, `type` |
| `ynab_account_uncleared_balance` | Uncleared balance of account   | `budget_id`, `budget_name`, `account_name`, `type` |
| `ynab_category_budgeted`         | Amount budgeted to category    | `category_group_name`, `category_name`             |
| `ynab_category_activity`         | Amount of activity in category | `category_group_name`, `category_name`             |
| `ynab_category_balance`          | Category balance               | `category_group_name`, `category_name`             |
| `ynab_exporter_api_calls_count`  | Count of calls to the YNAB API |                                                    |

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
$ export YNAB_API_TOKEN=${token}
$ go run ./cmd/ynab-exporter
```

### Running in Docker

You can run the exporter in a Docker container as follows:
```bash
$ export YNAB_API_TOKEN=${token}
$ docker build -t localhost/ynab-exporter:latest --no-cache .
$ docker run -ite YNAB_API_TOKEN=${YNAB_API_TOKEN} --name ynab-exporter --rm --publish 9090:9090/tcp localhost/ynab-exporter:latest
```

## TODO

* CI/CD configured with Github Actions (build, test, release, push to Docker Hub)
* Proper handling of API token secret (offer alternatives to envvar)
* Testing
* Deployment examples (both baremetal and containerized)
