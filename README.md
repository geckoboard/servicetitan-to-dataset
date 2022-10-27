# ServiceTitan to Geckoboard dataset

Push ServiceTitan reports into Geckoboard datasets

## Quickstart

### 1. Download the app

* macOS [x64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.2.0/servicetitan-to-dataset-darwin-amd64) / [arm64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.2.0/servicetitan-to-dataset-darwin-arm64)
* Linux [x86](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.2.0/servicetitan-to-dataset-linux-x86) / [x64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.2.0/servicetitan-to-dataset-linux-amd64)
* Windows [x86](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.2.0/servicetitan-to-dataset-windows-x86.exe) / [x64](https://github.com/geckoboard/servicetitan-to-dataset/releases/download/v0.2.0/servicetitan-to-dataset-windows-amd64.exe)

#### Make it executable (macOS / Linux)

On macOS and Linux you'll need to open a terminal and run `chmod u+x path/to/file` (replacing `path/to/file` with the actual path to your downloaded app) in order to make the app executable.

### 2. Create an oauth app

 - Create a service titan app in the [developer console here](https://developer.servicetitan.io/custom/my-apps/)
 - Then add the new app to your account under [Settings > integrations here](https://go.servicetitan.com/#/Settings/Api-Apps) and grant it access to reports.
 - Copy the values for later

From the developer console the **application key** field in the image below (ak1...) maps to the servicetitan `app_id` in the config.
<img width="650" alt="developer_oauth_app" src="https://user-images.githubusercontent.com/4930249/196000283-1630f560-20a9-4ff7-90bf-80fe51cb4ca5.png">

From the integrations page under the settings section those fields map nicely to the config.
<img width="650" alt="integrations_api_app_admin" src="https://user-images.githubusercontent.com/4930249/196000369-2dc791c0-6363-456a-8147-bfcb337b4e11.png">

### 3. Generate a config and update the values

Open up a terminal (on linux/max) or a command prompt (on windows), and run your script.

```
./servicetitan-to-dataset config --generate
```

This will generate an example config and create a file called config.yml by default.
Open this file and replace the servicetitan values with your specific values.

#### Geckoboard API

Hopefully this is obvious, but this is where your Geckoboard API key goes. You can find yours [here](https://app.geckoboard.com/account/details).


### 4. List categories and possible reports

Now we need to list all the reports and categories in your account.
This is required so that we can store what category and report we want to push to Geckoboard.


Run the following command

```
./servicetitan-to-dataset reports list
```

If the servicetitan config values was setup correctly you should see something
similar to the below, depending on the number categories you have you might see more output.

```
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

## Other

#### Custom dataset name

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

#### Dataset type

By default when pushing data to Geckoboard - we query just the first page of data and always replace the dataset contents
with the latest report data returned from ServiceTitan.

In some very rare cases - it maybe preferred to have the dataset append data, however for this to work you must have a
required and unique column or multiple columns..

For instance if you report contains a Date, Name and Number fields, the required fields would be Date and Name as these are the unique constraints.
By doing this we can build up data overtime if the query only returns data for "today"

To support appending data - you need to specify the dataset type as "append"

```yml
dataset:
  name: ...
  type: append
  required_fields:
    - Date
    - Name
```

#### Dynamic date parameters

If you're report requires a date parameter, you can hardcode a specific date such as 2022-10-19 however do so would require
the configuration to be updated every day.

To support a dynamic date value - you can use NOW - which for a date parameter will translate to the current date.
Also for cases that require a range e.g a from and to range. There is also a NOW-2 or NOW+2 which minus or plus n days

```yml
  parameters:
    - name: From
      value: "NOW-1"
    - name: To
      value: "NOW"
```

#### Refresh time

Once started, it can query ServiceTitan periodically and push the results to Geckoboard. Use this field to specify the time, in seconds, between refreshes.

Unfortunately due to some limitations with the new reports (beta) endpoint, ServiceTitan only allow query 1 report every 5 minutes.
This means that if you have;
 - 2 entries - then they will update every 10 minutes + the refresh time.
 - 10 entries - then  they will update every 50 minutes + the refresh time.

If you do not wish for it to run on a schedule, omit this option from your config and it will run only once after it has completed all entries.

```yml
refresh_time: 60
```

#### Environment variables

If you wish, you can provide any of the options under servicetitan and geckoboard as environment variables - to prevent storing secrets in the config.
This is possible using the following syntax `"{{YOUR_CUSTOM_ENV}}"`. Make sure to keep the quotes in there! For example:

```yaml
geckoboard:
  api_key: "{{GB_APIKEY}}"
```

