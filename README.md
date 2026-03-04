# FME - Flag Matching Engine


FME is a golang written program build atop lvlath library that answer a specific coding challenge.

Let's say you have a set of flags, it can be anything like command parameters or simply characters.
You want to combine them with each other and create a combination of flags.
Now, you want to apply rules to these combinations so not all of them are allowed.
For instance flag A requires flag B for a valid combination.

You need those rules because you have to forbid users to combine any kind of flags so they don't make your system crash. Best example, if you let a user combine command parameters as he wishes, he will inevitably fall into a wrong command.

FME is answering that kind of challenge by leveraging the power of graphs in order to efficiently prevent such issues before they happen.

To learn more about FME, I'm inviting you to read the documentation [docs/](docs/). 

# Usage

Here is a simple use case of FME.

I have an executable called `myexecutable.sh` and several parameters : -i, -t, -filepath, -folderpath and -dest.

I want to tell my engine the following :
- parameter `-i` requires parameter `-filepath` 
- parameter `-t` requires parameter `-folderpath` 
- parameter `-i` and `-t` requires parameter `-dest` 
- parameter `-i` interfers with parameter `-t` 

Here is a set of valid combination (or command in our case) : 
- myexecutable.sh -i -filepath /my/path/file.txt -dest /dest/path
- myexecutable.sh -t -folderpath /my/path/ -dest /dest/path

And here is a set of invalid ones :
- myexecutable.sh -i -filepath /my/path/file.txt
- myexecutable.sh -t -folderpath /my/path/ -i -filepath /my/path/file.txt

```go
func main() {
    script := "myexecutable.sh"
    parameters := []string{"-i", "-t", "-filepath", "-folderpath", "-dest"}
	
    schema := InitSchema()

	schema.AddConstraint("req", "-i", "-filepath")
	schema.AddConstraint("req", "-t", "-folderpath")
    schema.AddConstraint("req", "-i", "-dest")
    schema.AddConstraint("req", "-t", "-dest")
	schema.AddConstraint("interfer", "-i", "-t")

	_, err = NewCombination([]string{"a", "b", "c"}, schema)
	if err == nil {
		t.Fatalf("invalid combination accepted : %v", err)
	}
}
```


# Roadmap

- 04/03/26 : push of the rework + new version of lvlath
- 02/02/26 : rework of the engine with constraint interfaces + fix of the schema validation 
- 15/11/25 : improvements with interfaces and errors management
- 13/11/25 : the first version of FME is written
