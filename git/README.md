# Git package: `git`

The `git` package provides a relatively Git-specific but Github-independent model for objects like Branch, User, Repository and Pull Request. However, it does provide methods to create instances of [Google Go-Github package](https://gihubt.com/google/go-github) objects.  

This package's primary reason to exist is to wrap the Google Go-Github package, to add properties we need, to simplify use of the API by providing a `git.Context` object that can contain the most common parameters normally passed to the functions and methods of the Google Go-Github package.