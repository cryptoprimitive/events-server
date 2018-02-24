# events-server

The events-server provides an api to access the events log. It can also be
run with a '--serverMode testing' flag which opens up address balance
lookup and transaction lookup api's.

events API

/events/address returns events associated to the contract address

Address API

/addr/address returns balance of address (must have testing enabled)

Transaction API
/tx/txHash returns the transaction summary json (must have testing enabled)