the questions we have are
what is the db design?
what does the app do?
what are the apis we need?
why redis and how do we create db? do we use an orm or wha?
what kind of latency can we expect?
what is the read/write ratio?

db design- 
no login, only links and their shortened url to be stored
for key generation service, might need a table 

links table: Table Links
1. id
2. actual url
3. shortened url 
4. created at
5. expired at
6. number of clicks

for the database schema, please create a go struct as well. 
   
so for the combinations- 1000 adjectives, 1000 nouns, 4 char hash value. 

now for the kgs, essentially, the application shud not be creating a combination and then checking it in the db, rather it should pop it off a pre existing list of combinations. when a combination is popped off, then a new combination will be generated and stored in the list. 
so the steps are: 1. pre generate pairs of words in the background 
1. store them in a ready to use table or a redis set 
2. when a user requests a short url, the go server simply pops a key from the kgs

note on the click counter? it will prove to be a bottleneck for the db if multiple peopleare clicking on multiple links since updating the counter is a write operation which is time consuming. solution is 
1. handle the redirect immediately in go 
2. send a click event to an asynchronous buffer like a go channel or kafka queue. 
3. a background worker batches these clicks and updates the db once every 10 to 60 seconds. 

SELECT ... LIMIT 1000 FOR UPDATE SKIP LOCKED