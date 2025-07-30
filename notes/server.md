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

1. [x] basic client server ping pong

2. [ ] verify access and use secure protocols

we need some kind of authentication
so basically we need mutual auth
so theres tokens vs mtls, and i dont think tokens are sutable, yea well go with mtls

2.1. [ ] use https on server
generate self signed ssl certificate?
can i even use those here?

https://devcenter.heroku.com/articles/ssl-certificate-self

2.2. [ ] mutual tls?
2.3. [ ] success?

4. [x] figure out what to even send and how to sync dbs

### can we just send the sqlite files over https lol?
our db woulndt be very big, like a few mb max and thats with a looot of events
so apparently yes we can and we should try and benchmark

actually were gonna need to send the file anyway, at least initally
but not complicating logging changes and just sending sqlite files over is an explorable idea

### or, send the changes
on the client: after getting the inital file
, we store the changes (patches) as a log and send just those

- how do we know what the change was?
store last patch id we sent that client on the server
the server looks at last patch it sent us and just sends all the patches since then

- we apply changes to our local db untill theres internet duh
then its these 3 scenarios:

#### A. client has changes, server is the same

just send the patches to the server, it applies them and were guchi

#### B. client is the same, server has changes

just get the patches from the server and apply them on the client bam all good

#### C. client and server are both different, skull emoji

the servers db is the authority, thats the only way to do it yea

first get the patches from the server
, roll back the db to last true state
, and apply those

then we try to apply ours somehow
- if theres no conflict then were good
- if server added line and we add the same line, then just dont do shit, everyones happy

- if server changed line and we didnt, thats a buisness logic decision
, who wins can be decided on who was the last to change it or something
- if server deleted line and we changed it
, then we are fucked.
can we somehow ressurect it? so like, have a column in db that says deleted but dont delete the line
wait we cant change events in our app anyways, we delete them and create the same event but different
well that solves the ressurection problem

, and then its just scenario A, send our patches to the server and were done

after applying, erase the logs because we dont need those anymore

### how do we send the initial db?
i think we can just send the binary file lol, its looks the same on all platforms
and binary is the most efficient
i dont think we need to compress it, its a really small db, but maybe in the future

### what would a patch even look like?
can we just save all sql.execs we do along with timestamps
wait this is exactly what we need:

https://github.com/simukti/sqldb-logger

and so yea we can just send the json {sql string, timestamp}
is using sql strings like that even safe? dependency injection hello

### server side, what do we do?
execute patches and log them with {id, sql string, timestamp} :)

----------------------------------------------------

also check out these articles:
https://dev.to/ash_grover/how-i-designed-an-offline-first-app-an-outline-45c
https://felipeemidio.medium.com/offline-first-app-how-bad-is-it-to-build-one-ece1ffff4777
https://medium.com/offline-camp/security-in-offline-first-apps-59bf4800e82a
https://www.reddit.com/r/selfhosted/comments/pufhs0/beginner_guide_how_to_secure_your_selfhosted/
https://ryanisaacg.com/posts/db-internals-basics.html

server related todos here
- TODO: generate self signed ssl sertificates and use https and then mtls
- TODO: figure out how to work with json in go
- TODO: try to execute patches and see what theyd really need to look like
- TODO: log all sql execs in client into patches with https://github.com/simukti/sqldb-logger
- TODO: send the initial sqlite file in binary

i think thats enough for now yea
