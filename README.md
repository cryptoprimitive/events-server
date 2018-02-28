# events-server

The events-server provides an api to access the events log. It can also be
run with a '--serverMode testing' flag which opens up address balance
lookup and transaction lookup api's.

##Events API

    http://hostname/events/address
Returns events associated to the contract address

##Address API

    http://hostname/addr/address
Returns balance of address (must have testing enabled)

##Transaction API

    http://hostname/tx/txHash 
Returns the transaction summary json (must have testing enabled)