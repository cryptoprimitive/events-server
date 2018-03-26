# events-server

The events-server provides an api to access the events log.

## Sync API

    http://hostname/sync
Returns sync status of the node

## Events API

    http://hostname/events/address
Returns events associated to the contract address

## Block Events API

    http://hostname/blockevents/blocknumber
Returns all events from blocknumber

## Address API

    http://hostname/addr/address
Returns balance of address (must have testing enabled)

## Transaction API

    http://hostname/tx/txHash 
Returns the transaction summary json (must have testing enabled)

## Flags

### serverMode flag

Run with '--serverMode testing' to enable additional features
and testing output
 
### fromBlock flag

Set the --fromBlock flag to set what block to start the events search 
from.

### server flag

Set the rpc server to connect to. Default is "http://localhost:8545"