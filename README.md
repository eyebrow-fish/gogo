# gogo

`gogo` is an interpreted language for simple applications and shell-like use.
The syntax is created with pure simplicity in mind.

I was in inspired by [bash](https://www.gnu.org/software/bash/), not because it's good, but because
it's syntax always gets in the way of making simple scripts. Despite that, I use my CLI constantly.
I want to make my experience a little nicer. Maybe it will be nice for you too?

# mini-tutorial

`gogo` has variables:

```gogo
foo := "bar" # assigns "bar" to foo
foo = "bear" # foo is now "bear"
```

`gogo` has functions:

```gogo
print("Hello, World!")
```

`gogo` has semi-colon termination for beautiful multi-lining:

```gogo
prefix := "my name is"; suffix := " Alex"; print(prefix + suffix)
```

`gogo` can use shell programs:

```gogo
!echo 123 # everything after a "!" is executed as bash
```

`gogo` has basic control flow (and operators):

```gogo
a := 1; b := 2; c := a + b
if c == 3 then
    print("wow!")
else
    print("how?")
end
```
