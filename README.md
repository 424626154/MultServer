Hello everyone!

It is a simple network frame of server with golang

This structure of frame
--------------------------------------------------------------
-- Structure
--------------------------------------------------------------
.
├── lanstonetech.com
│   ├── common    					=>> (common packages)
│   │   ├── constant.go
│   │   └── rwnet.go
│   ├── network 					=>>	(network package)
│   │   ├── dispatcher.go 				(dispatcher process of packages)
│   │   ├── message.go 					(protocol message)
│   │   └── network.go 					(socket api)
│   ├── packet						=>>	(protocols)	
│   │   ├── C2M_Req_ShakeHand.go
│   │   ├── ID
│   │   │   └── ID.go
│   │   └── M2C_Resp_ShakeHand.go
│   └── server 						=>>	(servers)
│       └── LoginServer
│           ├── PackageHandler.go
│           └── Server.go
└── server 							=>> (Run)
    └── server.go

--------------------------------------------------------------
--------------------------------------------------------------

Thanks!

