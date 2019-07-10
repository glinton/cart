## cart

Cart is a helper tool to fetch **C**ircleCI build **art**ifacts and print them in a consumable manner. This makes it very friendly to post artifacts to users as the output can be pasted directly into a github comment.


#### Installation (requires golang environment)

```sh
go get github.com/glinton/cart
# to run
export CIRCLE_TOKEN="abc123"
$GOPATH/bin/cart -b 1234 -r yourname/repo
```

Add `export CIRCLE_TOKEN="abc123"` to your .bashrc to perist it as an environment variable.


#### Usage

```
cart -b buildID [-e] [-r vcs username/reponame] [-t vcs type]

-b
  Fetch artifacts for this build ID.

-e
  Enables output expansion (takes much more room in the comment as every os/arch is exapanded).

-r
  Set the vcs username/repository name to fetch build for.

-t
  Set the vcs type.
```

