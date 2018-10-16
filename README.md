# fetch-dcsl

## Install

`$ go get -v -u github.com/bertrandmt-nflx/fetch-dcsl`

## Run

```Usage of fetch-dcsl:
  -f	fetch fresh DCSL
  -m	list all manufacturers as CSV
  -s int
    	report on specific system ID
```

Most of the magic has to do with how the tool spits out bona-fide JSON.

Some tricks:
`$ fetch-dcsl -f | jq '.certificateStatusList.certificateStatus[].deviceInfo.manufacturer' | sort -f | uniq -ci | sort -k1 -n`
`$ fetch-dcsl -f | jq '.certificateStatusList.certificateStatus[] | select(.status == "STATUS_REVOKED")'`
