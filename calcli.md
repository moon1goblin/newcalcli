# ok cli for our calendar
call it kal? idk
could remap cal to our app or whatever

## new
-b --begin
-e --end
-n --name
-y --yes

cal new -b 12 07 12:55 -e 13 07 00:00 -n name

cal new -b 12 july 12 55 -n walk dog

maybe 12 07 always have two numbers and then if time follows then great, like then 12 is hours and the next number is minutes and so on up to ms or whatever well support
valid separators : . / or whitespace
can also make am/pm after time

what about stacking? -ben 12 07 15 00 n nah this doesnt work

should print out what it thinks we want to create and ask to confirm (--yes to pypass)


## list
-n --name
-b --begin
-e --end
-d --duration
-p --pretty

would return events that match that
if no arguments then we complain

-b is today by default
-e is 1 day by default

cal ls -b 12 07 -d 7
cal ls
cal ls

-e and -d cant be used simultaniously

june 12 12:00 name
july 12 14:00 shitfuck

for pretty it would group by day maybe?


## delete
-n --name
-b --begin
-e --end
-d --duration
-f --forever idk
-y --yes

cal rm walk dog

takes the output of cal ls and deletes it after conformation
should print what it found and ask for conformation (-y to bypass and confirm automatically)

should have a trash bin so we can do backsies 
delete trash after a week or something
or delete forever -f flag, then the event wouldnt go in the bin just straight up deleted
