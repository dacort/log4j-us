# DynaLog

Basically, I want a server that can dynamically generate log4j configs.

Simple API that takes a base config file (or empty?) and then adds the different lines.

People should be able to provide a base log4j template (ala codepen - https://codepen.io/kasei-dis/pen/JjYjwza) that they can append trace levels to.

## Usage

```shell
curl https://log4j.us/templates/metabase?trace=metabase.query-processor,metabase.driver
```

Spits out the [default Metabase](https://raw.githubusercontent.com/metabase/metabase/891e128b1f3dfad7e73250e54108148cba491678/resources/log4j.properties) log4j + `metabase.query-processor` set to trace.

