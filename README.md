# hd
hd cli utility

## Generate project config
### 1. Go to project workspace dir
```
cd /path/to/workspace
```
### 2. Run `gen config` command
```
hd gen config
```
> - project take current directory name by default
> - multiple app separates by comma `,`, if app is external accessible, append a colon `:` after the app name, followed by the port number, e,g: `gateway:1000`
