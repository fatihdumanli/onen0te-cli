# Onenote-CLI

A compact CLI tool to create/view the notes on your Onenote notebooks.

[![Build Status](https://dev.azure.com/fatihdumanli0884/onenote-cli/_apis/build/status/build-and-test?branchName=master)](https://dev.azure.com/fatihdumanli0884/onenote-cli/_build/latest?definitionId=24&branchName=master)
![Build Status](https://img.shields.io/github/workflow/status/fatihdumanli/onenote/build-and-test)



![preview](./img/nnote.gif)

# What is it?
You can quickly take notes on your terminal and save them as Onenote pages. It's also possible to take a note on your favorite text editor, onenote cli tool will save the note to the specified section upon quitting the editor. 

Please note that this is an unofficial onenote client and in order to use this app you need to authorize this app to access and write your OneNote notes. See the **Authentication** section for more info.

# Features
- Take inline notes
- Take notes on your favorite editor
- Import contents of a file as a Onenote Page
- Use aliases to quick access to your sections
- [ ] Browse your notes
- [ ] Search in notebooks

# Installation 
This tool is written in Go. Run the command below to install the Onenote CLI.

```bash
$ go install github.com/fatihdumanli/onenote@latest
```

## Authentication
Authentication is done during your very first interaction with nnote. To use this application, you must authorize nnote to access/write your Onenote notebooks and sections.

Feel free to change the **ClientId** and **TenantId** variables with yours [here](https://github.com/fatihdumanli/onenote-cli/blob/master/options.go#L12). You can grab yours on Azure portal. See the following link for further information.

[https://docs.microsoft.com/en-us/graph/auth-v2-user?context=graph%2Fapi%2F1.0&view=graph-rest-1.0#1-register-your-app](https://docs.microsoft.com/en-us/graph/auth-v2-user?context=graph/api/1.0&view=graph-rest-1.0#1-register-your-app)

# Usage
```bash
Usage:
  nnote [command]

Available Commands:
  alias       add/list and remove alias
  browse      browse the pages within a onenote section
  new         Create a new note
```

## Creating a New Note

Use `new` command to take notes on your Onenote account. You can use the following flags when taking notes.

### Usage

```bash
Usage:
  nnote new [flags]

Aliases:
  new, add, save

Flags:
- `-a`: Alias for the section
- `-t`: Title for the page (default is empty)
- `-i`: Use this flag to save inline note. Wrap your text with double quotes right after the flag literal.
- `-f`: Read contents of the file and save it as a Onenote page.
```
### Take an inline note
To take an inline note, use the command `new` with flag `-i`


**Example 1 - Taking an inline note**
```bash
$ nnote new -i "taking inline notes are fun" 
``` 

**Example 2 - Taking an inline note using an alias**

```bash
$ nnote new -i "aliases help you to save time" -a fooSection
```

**Example 3 - Taking an inline note with a title**

```bash
$ nnote new -i "titles makes it easy to locate your notes!" -t "this is crazy important note!"
``` 

### Take note using your favorite editor
To take a note using your favorite editor, do not send any flags or argument along with the `new` command. To launch your default text editor to save a Onenote page run the following command.

```bash
$ nnote new
```
Upon quitting the editor, you'll be prompted to choose the notebook and section to which the note'll be uploaded.

### Import the contents of a file as Onenote page
Use `-f` flag to save the contents of a file as a Onenote page.


```bash
$ cnote new -f /path/to/the/file.txt -t "it's recommended to speficy a title although it's not required!" -a barSection
```

> NOTE: You'll be prompted to choose notebook and section if you don't specify `-a` flag.

## Aliases
It's encourged you to tag your sections with aliases. It'll facilitate the process of taking a new note as you'll be skipping time consuming HTTP requests to fetch your notebooks and sections.

### Usage
```bash
Usage:
  nnote alias [command]

Available Commands:
  list        display alias list
  new         create a new alias.
  remove      remove an alias
```

### Creating a new alias
Use the following command to create a new alias.

TODO: Currently It does not work in this way
```bash
$ nnote alias new [alias]
```

You'll be prompted to select notebook and section in order. 

**Example**
```bash
$ nnote alias new "quick notes"
```

```bash
$ nnote alias new "asp.net"
```

### Listing aliases
Use `list` command to list your aliases and the corresponding sections.

```bash
$ nnote alias list
```

### Removing aliases
You may want to change the mapping for an alias. Use `remove` command to remove an alias.

```bash
$ nnote alias remove <alias>
``` 

**Example**
```bash
$ nnote alias remove "asp.net"
```

## View your notes
TODO: Update
## Edit your notes
TODO: Update

# Contribution
# Licence




