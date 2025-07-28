# how to server

when user logs into cal, app sends fetch query to the server
then comparing local and server db, updating the older one
need to keep track of last change date??
think what to send: all db, or difference
need to check rights to send and get things from/to server
maybe transfer data with json
send db or lone events?
make ping
On server: db, posmotrim logicu
proga for client and server to send simple queries

----------------------------------------------------

for grep, TODO: server todos here:

1. [x] basic client server ping pong

2. [ ] verify access and use secure protocols

we need some kind of authentication
so basically we need mutual auth
so theres tokens vs mtls, and i dont think tokens are sutable, yea well go with mtls

2.1. [ ] use https on server
generate self signed ssl certificate?
can i even use those here?
can just use http.
2.2. [ ] mutual tls?
2.3. [ ] success

3. [ ] somehow jsonify, send json, parse it

4. [ ] figure out what to even send

db sync with offline first

----------------------------------------------------

also check out these articles:
https://dev.to/ash_grover/how-i-designed-an-offline-first-app-an-outline-45c
https://felipeemidio.medium.com/offline-first-app-how-bad-is-it-to-build-one-ece1ffff4777
https://medium.com/offline-camp/security-in-offline-first-apps-59bf4800e82a
https://www.reddit.com/r/selfhosted/comments/pufhs0/beginner_guide_how_to_secure_your_selfhosted/
