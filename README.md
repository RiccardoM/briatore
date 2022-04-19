# Briatore - Italy tax report helper
Briatore takes its name from Flavio Briatore, one of the Italy's most famous tax evader of the recent history. 

It represents a tool that can help you fill out your tax report if you live in Italy. To do so, given the details of a Cosmos-based chain, it gets the following data from the chain itself: 
- your balance as per 1st January
- your balance as per 31st December

Balances are returned both in amount of coins as well the corresponding amount of Euros they were worth at those dates.

## Usage
1. Copy the below config somewhere, and edit the provided addresses for each chain you want (delete the accounts for the chains you don't want to parse)
2. Install the binary running `make install`
3. Run the script with the following command: 
    ```
   briatore report 2021-12-31T23:59:59Z cosmos1...,juno1... --home /path/to/dir/where/config/file/is
   ```
   
## Example config file
```yaml
report:
  currency: "eur"

chains:
   - name: "Osmosis"
     rpcAddress: "https://rpc.osmosis.forbole.com:443"
     grpcAddress: "https://grpc.osmosis.forbole.com:443"
     bech32Prefix: "osmo"

  - name: "Cosmos"
    rpcAddress: "https://rpc.cosmoshub.forbole.com:443"
    grpcAddress: "https://grpc.cosmoshub.forbole.com:443"
    bech32Prefix: "cosmos"
```