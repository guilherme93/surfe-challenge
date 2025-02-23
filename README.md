# Surfe Challenge

## Config
 There is a small json config file to set the port where you would like the api to run.
    
Please not if you change the port and you're running in docker you also need to change the port in the docker compose file 
#### Prerequisites

- go 1.23.6

Before running test or lint please run
```shell
make deps
```

#### Run

The command below will compile and run the program

```shell
make run
```

in case you prefer docker you can run
- docker compose up