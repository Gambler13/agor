# Start
Agar.io server implementation written in Go. This project is still work in progress
The JavaScript client can be found under: 
## Configuration 
config.yaml
```
---
TickRate: 50        #Tick rate for game loop (1/tickRate)
Port: 8008          #Server port
World:              #World settings
  Width:     1200   #World width
  Height:    800    #World height
  MaxPlayer: 100    #Max players
  Food:      500    #Max food (non regenerative)
```
## Env variables
`AGOR_ORIGIN` Set host for Access-Controll-Allow-Origin (for dev e.g. http://localhost:8084) This is needed during local development
`AGOR_CONFIG` Agor config file path. Default path is `/etc/agor/config/default.yaml`
## Live demo
https://agor.hexhibit.xyz 