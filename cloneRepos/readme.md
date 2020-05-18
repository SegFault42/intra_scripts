### This script clone massive repository

1) Add you ssh public key here : https://vogsphere.msk.21-school.ru/user/settings/keys
2) Get the list of repo that you want to clone from metabase (download as json) : https://metabase.21-school.ru/question#eyJkYXRhc2V0X3F1ZXJ5Ijp7ImRhdGFiYXNlIjo1LCJ0eXBlIjoicXVlcnkiLCJxdWVyeSI6eyJzb3VyY2UtdGFibGUiOjI0NH19LCJkaXNwbGF5IjoidGFibGUiLCJ2aXN1YWxpemF0aW9uX3NldHRpbmdzIjp7fX0=
3) Exec the program 

Usage example :

```
	$ go run main.go -j=file.json
```

All repos will be cloned in the working directory.

It will be stored in login diretory
login/{repository-content}
