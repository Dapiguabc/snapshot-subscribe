# snapshot-subscribe
 A backend serverice for the plugin named Subscribe
 
 Front Plugin Repository: 
  * https://github.com/Dapiguabc/snapshot.js 
  * https://github.com/Dapiguabc/snapshot
  
 Backend Server Repository:
  * https://github.com/Dapiguabc/snapshot-subscribe
# Usage
### Download code.
 ```git clone https://github.com/Dapiguabc/snapshot-subscribe.git```
### Enter into the directory.
 ```cd snapshot-subscribe```
### Modify the config.yaml file
 ```
 mailConn:
  name: "Snapshot" // Email alias
  user: "example@gmail.com" // Email address
  pass: "xxxxxxxxxxxxx" // Email password(smtp)
  host: "smtp.gmail.com" // SMTP address
  port: "888"
 graphql:
  url: "https://testnet.snapshot.org/graphql"  // snapshot graphql api
```
### Build server
```go build -i -o subscribe.exe```
# Demo
## Subscribe a proposal
![image](https://github.com/Dapiguabc/snapshot-subscribe/blob/main/demo/font.gif)
## You will reveive an eamil like this when the state of the proposal you subscribed changed 
![image](https://github.com/Dapiguabc/snapshot-subscribe/blob/main/demo/mail.png)
