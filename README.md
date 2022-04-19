# Briatore - Italy tax report helper
Briatore takes its name from [Flavio Briatore](https://en.wikipedia.org/wiki/Flavio_Briatore), one of the Italy's most famous tax evader of the recent history. 

It represents a tool that can help you fill out your tax report if you live in Italy. To do so, given the details of a Cosmos-based chain, it gets the following data from the chain itself: 
- your balance as per 1st January
- your balance as per 31st December

Balances are returned both in amount of coins as well the corresponding amount of Euros they were worth at those dates.

## Usage
1. Copy the below config somewhere
2. Install the binary running `make install`
3. Run the script with the following command: 
    ```
   briatore report 2021-12-31T23:59:59Z cosmos1...,juno1... --home /path/to/dir/where/config/file/is
   ```

> NOTE  
> The reported value is currently returned in Euro (EUR).
   
## Example config file
```yaml
report:
  currency: "eur"

chains:
   - name: "Osmosis"
     rpcAddress: "https://rpc....:443"
     grpcAddress: "https://grpc....:443"
     bech32Prefix: "osmo"

  - name: "Cosmos"
    rpcAddress: "https://rpc....:443"
    grpcAddress: "https://grpc....:443"
    bech32Prefix: "cosmos"
```

## APIs
Aside from the `report` command, Briatore also contains the `start` command that allows to start a new REST server exposing the following endpoints.

### Endpoints
#### GET `/report`
Returns the report data for the provided addresses and date, in the given output format. 

|  Parameter  |                             Type                             | Description                                                                    |
|:-----------:|:------------------------------------------------------------:|:-------------------------------------------------------------------------------|
|   `date`    | [RFC339 Date](https://datatracker.ietf.org/doc/html/rfc3339) | Date for which to get the report (ideally end of year - `2021-12-31T23:59:59Z` |
| `addresses` |                String <br/>(comma separated)                 | List of addresses for which to get the report                                  |
|  `outpout`  |                            String                            | Format in which to return the data (supported formats: `csv`, `text`, `json`)  |

### Live instance 
If you don't want to run your own instance by specifying your own nodes, you can use the one running at `http://162.55.171.213:8080/`:

```
http://162.55.171.213:8080/report?output=csv&date=2021-12-31T23:59:59Z&addresses=cosmos1...,juno1...
```