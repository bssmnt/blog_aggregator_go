# blog aggregator project built in golang utilising postgresql

you need to install postgres and go to run this program, can install them either via brew or via the packages you can find online

homebrew instructions:
install homebrew (macos and linux)
````
/bin/bash -c "$(curl -fsSL
https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
````

then install postgres
````
brew install postgresql
````

and you can install go via homebrew as well, but i recommend downloading it directly from https://go.dev, as it's easier when it comes to setting up GOROOT and GOPATH etc;
alternatively if you want to install via brew you can use
````
brew install go
````

you'll also need to install gator, which is the tool i've use for evaluating gatekeeper constraint templates and constraints in a local env
````
brew install gator
````

### some of the commands available in the program:

reset: resets program and runs down migrations followed by up migrations (reset)

register: register a new user (register "username")

login: login to an already registered user (login "username")

users: lists all users currently on the system, tells you who is currently logged in (users)

agg: aggregates data from all currently followed feeds, showing timestamps and post titles (agg "time")

browse: allows you to browse more info about a certain number of posts, default is 2 unless specified (browse "number of posts")

follow: allows logged-in user to follow a feed (follow "url")

unfollow: allows logged-in user to unfollow a feed (unfollow "url")