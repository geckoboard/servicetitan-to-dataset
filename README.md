# ServiceTitan to Geckoboard dataset

Push ServiceTitan reports into Geckoboard datasets

## Quickstart

### 1. Download the app

* macOS [x64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.0.1/servicetitan-to-dataset-darwin-amd64) / [arm64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.0.1/servicetitan-to-dataset-darwin-arm64)
* Linux [x86](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.0.1/servicetitan-to-dataset-linux-386) / [x64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.0.1/servicetitan-to-dataset-linux-amd64)
* Windows [x86](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.0.1/servicetitan-to-dataset-windows-386.exe) / [x64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.0.1/servicetitan-to-dataset-windows-amd64.exe)

#### Make it executable (macOS / Linux)

On macOS and Linux you'll need to open a terminal and run `chmod u+x path/to/file` (replacing `path/to/file` with the actual path to your downloaded app) in order to make the app executable.

### 2. Create an oauth app

 - Create a service titan app in the [developer console here](https://developer.servicetitan.io/custom/my-apps/)
 - Then add the new app to your account under [Settings > integrations here](https://go.servicetitan.com/#/Settings/Api-Apps) and grant it access to reports.
 - Copy the values for later

### 3. Generate a config and update the values

Open up a terminal (on linux/max) or a command prompt (on windows), and run your script.

```
./servicetitan-to-dataset config --generate
```

This will generate an example config and create a file called config.yml by default.
Open this file and replace the servicetitan values with your specific values.


### 4. List categories and possible reports

Now we need to list all the reports and categories in your account.
This is required so that we can store what category and report we want to push to Geckoboard.


Run the following command

```
./servicetitan-to-dataset reports list
```

If the config servicetitan values was setup correctly you should see something
similar to the below, depending on the number categories you have you might see more output.

```sh
2022/10/14 16:10:27 Fetching categories...
2022/10/14 16:10:41 Fetching reports for category Accounting...
2022/10/14 16:10:41 Fetching reports for category Other...
+-----------+----------------------+----------------------+-------------------------+
| REPORT ID |     CATEGORY ID      |    CATEGORY NAME     |       REPORT NAME       |
+-----------+----------------------+----------------------+-------------------------+
|  1234     | accounting           | Accounting           | Revenue by employee     |
|           |                      |                      |                         |
+-----------+----------------------+----------------------+-------------------------+
|  2345     | performance-reports  | Performance Reports  | Sales by technician     |
|           |                      |                      |                         |
+-----------+----------------------+----------------------+-------------------------+
```

### 5. Query the report fields and parameters

Now you have a category ID and report ID (lets take the second example above).
Now we need to query the report fields and parameters

```
./servicetitan-to-dataset reports parameters --report 2345 --category performance-reports
```

That should output something like the following;

```
Report id:  2345
Report name: Sales by technician

Report fields:
+--------------------------+------------------------+--------+
|        FIELD NAME        |         LABEL          |  TYPE  |
+--------------------------+------------------------+--------+
| Technician Name          | TechnicianName         | String |
+--------------------------+------------------------+--------+
| Completed Jobs           | CompletedJobs          | Number |
+--------------------------+------------------------+--------+
| Total revenue            | TotalRevenue           | Number |
+--------------------------+------------------------+--------+

Report parameters:
+-----------------+------------------------------+-----------+--------+-----------+
|  PARAMTER NAME  |            LABEL             | DATA TYPE | ARRAY? | REQUIRED? |
+-----------------+------------------------------+-----------+--------+-----------+
| From            | From                         | Date      | FALSE  | TRUE      |
+-----------------+------------------------------+-----------+--------+-----------+
| To              | To                           | Date      | FALSE  | TRUE      |
+-----------------+------------------------------+-----------+--------+-----------+
| JobTypes        | Job Type                     | String    | FALSE  | FALSE     |
+-----------------+------------------------------+-----------+--------+-----------+
| IncludeInactive | Include Inactive Technicians | Boolean   | FALSE  | FALSE     |
+-----------------+------------------------------+-----------+--------+-----------+
```

### 6. Add a new entry to the config

With all the information now about the report and the fields and parameters we can
add an entry to the config.

When pushing data to Geckoboard at least one required field is expected. This could be
a ID field or a date creation field, or a name field that is always present.

In our example the "Technician Name" will always be present so we add that to the config.

In our report example as well we have 2 required parameters (you might have none in which case you can omit the parameters section)
However if you have required parameters then they must be supplied in the config.

Here we use NOW and NOW-1 a special keyword for date parameters - that are dynamic based on the current datetime.

```yml

entries:
  - report:
      id: 2345
      category_id: performance-reports
      parameters:
        - name: From
          value: "NOW-1"
        - name: To
          value: "NOW"
    dataset:
      required_fields:
        - Technician Name
```

### Custom dataset name

By default the default name is automatically implied from the report name and converted to ensure it is safe.

You may provide a dataset name in the config as such if you want to use the same report with different params for instance

```yml

entries:
  - report:
      ...
    dataset:
      name: "My custom dataset name"
      required_fields:
        ...
```

Please note that in some cases the dataset name you input will be converted based on the rules of what a valid dataset name is
but all letters and numbers are valid values

### refresh_time

Once started, it can query ServiceTitan periodically and push the results to Geckoboard. Use this field to specify the time, in seconds, between refreshes.

Unfortunately due to some limitations with the new reports (beta) endpoint, ServiceTitan only allow query 1 report every 5 minutes.
This means that if you have;
 - 2 entries - then they will update every 10 minutes + the refresh time.
 - 10 entries - then  they will update every 50 minutes + the refresh time.

If you do not wish for it to run on a schedule, omit this option from your config and it will run only once after it has completed all entries.


#### Environment variables

If you wish, you can provide any of the options under servicetitan and geckoboard as environment variables - to prevent storing secrets in the config.
This is possible using the following syntax `"{{YOUR_CUSTOM_ENV}}"`. Make sure to keep the quotes in there! For example:

```yaml
geckoboard:
  api_key: "{{GB_APIKEY}}"
```

### Geckoboard API

Hopefully this is obvious, but this is where your Geckoboard API key goes. You can find yours [here](https://app.geckoboard.com/account/details).
