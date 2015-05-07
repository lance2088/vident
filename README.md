# Vident
Vident is a scripting language written in Go for a tutorial series. It's an interpreter implemented
in the Go programming language.

## Blog
I've been moving my blog around a ton, I've finally settled
with using Ghost, so go check out my blog [here](http://blog.felixangell.com).

[Here's a link to the first article in the series](http://blog.felixangell.com/part-1-lets-build-an-interpreted-language-in-go/).

## Example:

    let x = 5 + 5
    let add(a, b) = {
        a + b
    }
    let z = add(x, 5)
    print(z)
