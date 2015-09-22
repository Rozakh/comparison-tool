Tool, that compares responses from different hosts and generate html report. Uses Splunk to get urls list.

* Creates config file at first start.
* Url pattern is used to find urls in Splunk.
* Pattern can include named groups to extract certain values. Group name must be placed in double curly braces.

Example:

```
{
  "URLPatternA": "/{{param1}}/{{param2}}/"
  "URLPatternB": "?param1={{param1}}&param2={{param2}}"
}
```

Values of {{param1}} and {{param2}} will be extracted from the Splunk results. Placeholders from the URLPatternB will be filled with
proper values.

* You can specify a regular expression for certain group name in the "Groups" part of the config.

Example:

```
{
  "Groups": {
    "year": "\\d{4}"
  }
}
```

After that, {{year}} will be replace by \d{4} in the Splunk search query.
