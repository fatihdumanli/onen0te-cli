# onen0te-CLI

create & view onenote notes on terminal.

![Build Status](https://img.shields.io/github/workflow/status/fatihdumanli/onen0te-cli/Test)
[![Build Status](https://dev.azure.com/fatihdumanli0884/onenote-cli/_apis/build/status/build-and-test?branchName=master)](https://dev.azure.com/fatihdumanli0884/onenote-cli/_build/latest?definitionId=24&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fatihdumanli/onen0te-cli)](https://goreportcard.com/report/github.com/fatihdumanli/onen0te-cli)


![onenote](https://github.com/fatihdumanli/onen0te-cli/assets/7475244/9f05d75a-9f75-4fc0-a913-8e5d71956a44)

# What is it?
This tool (will be referenced as 'nnote') is used to view & create onenote pages on the terminal using Microsoft Graph APIs. 

Please bear in mind that this is an unofficial Onenote client and in order to use this app you need to authorize this app to access and write your OneNote notes. See the **Authentication** section for more info.

# Features
- Take inline notes
- Take notes on your favorite editor
- Import contents of a file as a Onenote Page
- Use aliases to quick access your sections
- Browse in your notes
- Search in notebooks
- Display notes on the web browser
- Display notes on the Onenote client


# Installation 
This tool is written in Go. Run the command below to install the Onenote CLI.

```bash
$ go install github.com/fatihdumanli/onen0te-cli/cmd/nnote@latest
```

## Authentication
Authentication is done during your very first interaction with nnote. To use this application, you must authorize nnote to access/write your Onenote notebooks and sections.

Feel free to change the **ClientId** and **TenantId** variables with yours [here](https://github.com/fatihdumanli/onen0te-cli/blob/master/options.go#L14). You can grab yours on the Azure portal. See the following link for further information.

[https://docs.microsoft.com/en-us/graph/auth-v2-user?context=graph%2Fapi%2F1.0&view=graph-rest-1.0#1-register-your-app](https://docs.microsoft.com/en-us/graph/auth-v2-user?context=graph/api/1.0&view=graph-rest-1.0#1-register-your-app)

# Usage
```

Usage:
  nnote [command]

Available Commands:
  alias       add/list and remove alias
  browse      browse the pages within a onenote section
  help        Help about any command
  new         create a new note
  search      do a search in your notes
```

## Creating a New Note

Use `new` command to take notes on your Onenote account. You can use the following flags when taking notes.

### Usage

```
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
$ nnote new -i "taking inline notes is fun" 
``` 
**Example 2 - Taking an inline note using an alias**

```bash
$ nnote new -i "aliases help you to save time" -a foo
```

**Example 3 - Taking an inline note with a title**

```bash
$ nnote new -i "titles makes it easy to locate your notes!" -t "this is crazy important note!"
``` 

![save-inline](https://github.com/fatihdumanli/onen0te-cli/assets/7475244/e4484166-525b-4dfe-be60-dd366e5b4901)

<br/>

### Take notes using your favorite editor
To take a note using your favorite editor, do not specify any flag or argument along with the `new` command. To launch your default text editor to save a Onenote page run the following command.

```
$ nnote new
```
You'll be prompted to choose the notebook and section to which the note will be uploaded right after quitting the editor.

![new-editor](https://github.com/fatihdumanli/onen0te-cli/assets/7475244/43828c4f-7bba-4e5a-8a97-80711a233de8)

<br/>

### Import the contents of a file as a OneNote page
Use `-f` flag to save the contents of a raw file as a Onenote page.


```bash
$ nnote new -f /path/to/the/file.txt -t "title-of-the-page" -a quicknotes
```

> NOTE: You'll be prompted to choose notebook and section if you don't specify `-a` flag.

## View your notes

Use the `browse` command to browse in your notes. You'll be prompted to select a notebook, section, and page. The page content will be rendered on the terminal. You can also view your note content in web browser or Onenote desktop client.

You can navigate between your OneNote section/pages while displaying note content. You don't need to run `browse` command each time you want to go through your notes.

```
$ nnote browse
```

![browse](https://github.com/fatihdumanli/onen0te-cli/assets/7475244/89afdc76-942d-43ef-a5a2-3166840280b2)


## Search in your notes

Use the `search` command to search in your notes. This command will perform a search in your all notebooks and will prompt you to select one of the results to view the selected note. You can also view your notes in the web browser, or Onenote desktop client.

**Example**
```
$ nnote search "redis"
```

![search-in-notes](https://github.com/fatihdumanli/onen0te-cli/assets/7475244/1b240126-4cbb-43fe-943b-61fd7bd75be6)

## Aliases
It's recommended you to tag your sections with aliases. It'll facilitate the process of taking a new note as you'll be skipping time-consuming HTTP requests to fetch your notebooks and sections.

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

```bash
$ nnote alias new [alias]
```

You'll be prompted to select notebook and section in order. 

**Example**
```bash
$ nnote alias new "quick notes"
```

```bash
$ nnote alias new "elasticsearch"
```

### Listing aliases
Use the `list` command to list your aliases and the corresponding sections.

```bash
$ nnote alias list


Alias         | Section                | Notebook
git           | Git                    | Fatih's Notebook
kafka         | Kafka                  | Fatih's Notebook
microservices | Building Microservices | Fatih's Notebook
postgresql    | PostgreSQL             | Fatih's Notebook
qn            | Quick Notes            | Fatih's Notebook
sql           | SQL Cookbook           | Fatih's Notebook
vocab         | Inner Vocabulary       | Fatih's Notebook
```

### Removing an alias
You may want to change the mapping for an alias. Use the `remove` command to remove an alias.

```bash
$ nnote alias remove <alias>
``` 

**Example**
```bash
$ nnote alias remove "qn"
```


# Thanks 

- [prologic/bitcask](https://git.mills.io/prologic/bitcask)
- [AlecAivazis/survey](https://github.com/AlecAivazis/survey)
- [k3a/html2text](https://github.com/k3a/html2text)
- [pterm](https://github.com/pterm/pterm)
- [spf13/cobra](https://github.com/spf13/cobra)

# Contribution
All PRs are welcome! If you think you've discovered a bug, please open an issue first.

# Licence
Released under the terms of the [MIT Licence](https://opensource.org/licenses/MIT).

