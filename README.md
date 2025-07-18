# Omise Go challenge submission

You can access the original README file from [README.statement.md](README.statement.md)

## Project structure
- **cmd**: contains entrypoint of executable
- **data**: contains encrypted CSV of donation records
- **internal**:
  - **cipher**: original Rot128 encryption and decryption code
  - **donation**: contains Omise client and donation process (token creation, charging)
  - **fileio**: handling file and CSV
  - **reporting**: produces donation summary
- **.env.example**: sample env file contains Omise key configuration

## Configure Omise keys

First we need to set Omise public and secret key in env file. The config example has been provided in `.env.example`. You can set the keys by copying the example to the actual env file
```
cp .env.example .env
```
Then replace each keys in `.env`
```
OMISE_PUBLIC_KEY=pkey_xxxxxxxxxxxxxxxxxxxxxxxx
OMISE_SECRET_KEY=skey_xxxxxxxxxxxxxxxxxxxxxxxx
```

## Build and run

Build the executable using command
```
go build -o go-tamboon cmd/go-tamboon/main.go
```

Then you can run the executable. For example, you want to process the donation from file `data/fng.1000.csv.rot128`
```
./go-tamboon data/fng.1000.csv.rot128
```

## Technical notes

This code submission has been designed to allocate memory safely by cleaning all financial information from memory after use to make sure to not leave any sensitive data.

I also try to achieve smallest memory allocation. But in this implementation you can see that the program did store donation amount for all donors in memory, not only the top donors shown in summary after processing.

The reason behind this is because it's possible that some donors can donate multiple times so we need to accumulate all donations for each donors for the accurate top donors records. If we can make sure that all donations has been made by unique donors, we can use more optimal data structure like heap from `container/heap` to build priority queue to keep only top donors, which will allocate memory in much more efficient way.

There are also the effort of throttling API calls to Omise and utilizing CPU cores to work efficient concurrently in donation process.