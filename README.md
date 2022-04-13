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
   briatore report 2021-12-31T23:59:59Z --home /path/to/dir/where/config/file/is
   ```
   
## Example config file
```yaml
report:
  currency: "eur"

chains:
  - name: "Cosmos"
    rpc_address: "https://65.21.93.108:10457"
    grpc_address: "https://65.21.93.108:10490"
  - name: "Akash Network"
    rpc_address: "https://rpc.akash.forbole.com:443"
    grpc_address: "https://grpc.akash.forbole.com:443"
  - name: "Cronos"
    rpc_address: "https://rpc.crypto-org.forbole.com:443"
    grpc_address: "https://grpc.crypto-org.forbole.com:443"
  - name: "Regen"
    rpc_address: "https://rpc.regen.forbole.com:443"
    grpc_address: "https://grpc.regen.forbole.com:443"
  - name: "e-Money"
    rpc_address: "https://rpc.emoney.forbole.com:443"
    grpc_address: "https://grpc.emoney.forbole.com:443"
  - name: "Likecoin"
    rpc_address: "https://rpc.likecoin.forbole.com:443"
    grpc_address: "https://grpc.likecoin.forbole.com:443"
  - name: "Terra"
    rpc_address: "https://rpc.terra.forbole.com:443"
    grpc_address: "https://grpc.terra.forbole.com:443"
  - name: "Bitsong"
    rpc_address: "https://rpc.bitsong.forbole.com:443"
    grpc_address: "https://grpc.bitsong.forbole.com:443"
  - name: "Desmos"
    rpc_address: "https://rpc.desmos.forbole.com:443"
    grpc_address: "https://grpc.desmos.forbole.com:443"
  - name: "Band"
    rpc_address: "https://rpc.band.forbole.com:443"
    grpc_address: "https://grpc.band.forbole.com:443"
  - name: "Sifchain"
    rpc_address: "https://rpc.sifchain.forbole.com:443"
    grpc_address: "https://grpc.sifchain.forbole.com:443"

accounts:
  - chain: "Cosmos"
    addresses:
      - "<Your Cosmos Hub address>"
  
  - chain: "Akash Network"
    addresses:
      - "<Your Akash Network address>"

  - chain: "Cronos"
    addresses:
      - "<Your Chronos address>"
```