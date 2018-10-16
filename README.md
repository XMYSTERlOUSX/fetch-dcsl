# fetch-dcsl

## Install

`$ go get -v -u github.com/bertrandmt-nflx/fetch-dcsl`

## Run

```
Usage of fetch-dcsl:
  -f	fetch fresh DCSL
  -m	list all manufacturers as CSV
  -s int
    	report on specific system ID
```

Most of the magic has to do with how the tool spits out bona-fide JSON.

Some tricks:

```
$ fetch-dcsl -f |
  jq '.certificateStatusList.certificateStatus |
      map(select(.deviceInfo.manufacturer != null) | .deviceInfo.manufacturer |= ascii_downcase)' |
  jq -s 'map({manufacturer: .[].deviceInfo.manufacturer}) |
         group_by(.manufacturer) |
         map ({manufacturer: .[0].manufacturer, count: length}) |
         sort_by(.count)'
```

```
$ fetch-dcsl -f |
  jq '.certificateStatusList.certificateStatus |
      map(select(.status == "STATUS_REVOKED"))'
```

```
$ fetch-dcsl -f |
  jq '.certificateStatusList.certificateStatus |
      map (select(.deviceInfo.manufacturer != null) |
           .deviceInfo.manufacturer |= ascii_downcase |
           select(.deviceInfo.manufacturer == "netflix")) |
      sort_by(.deviceInfo.systemId)'
```
