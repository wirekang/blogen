# Install
```
go install github.com/wirekang/blogen@latest
```

# Usage

## Init
```
blogen -i
```
.  
├── blogen.cfg  
├── gen  
├── md  
├── out  
└── template  
  
      
### blogen.cfg
Set base information of blog.
```
title = Title of blog
addr = //blog.example.com
```

### md/
Add posts in directory **md/**. All posts must be written in following format:
```
title= Title of post
date = 2020-1-2
tags = tag1, tag 2, tag3
show = true
##blogen##

...
markdown post
...
```
Without **show** option, the post is ignored.
### template/
Customize blog design.

## Generate
```
blogen
```
Now directory **out/** is static site.

# Example
[example site](https://wirekang.github.io/blogen/example/out)  
[my blog](https://wirekang.github.io)
